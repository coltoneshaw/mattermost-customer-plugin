package sqlstore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/mattermost/mattermost/server/public/model"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"
)

// customerStore holds the information needed to fulfill the methods in the store interface.
type customerStore struct {
	pluginAPI          PluginAPIClient
	store              *SQLStore
	queryBuilder       sq.StatementBuilderType
	customerSelect     sq.SelectBuilder
	packetValuesSelect sq.SelectBuilder
	configValuesSelect sq.SelectBuilder
	pluginValuesSelect sq.SelectBuilder
}

type sqlCustomers struct {
	app.FullCustomerInfo
}

type sqlPacket struct {
	app.CustomerPacketValues
}

type sqlConfig struct {
	Config json.RawMessage `db:"config"`
}

const (
	customerTable = "crm_customers"
	packetTable   = "crm_packetValues"
	configTable   = "crm_configValues"
	pluginTable   = "crm_pluginValues"
	auditTable    = "crm_audit"
)

type UpdateType string

const (
	Packet UpdateType = "packet"
	User   UpdateType = "user"
)

func applyCustomerFilterOptionsSort(builder sq.SelectBuilder, options app.CustomerFilterOptions) (sq.SelectBuilder, error) {
	var searchTerm string
	if options.SearchTerm != "" {
		searchTerm = "%" + options.SearchTerm + "%"
		builder = builder.Where(sq.Or{
			sq.ILike{"ci.name": searchTerm},
			sq.ILike{"ci.licensedTo": searchTerm},
			sq.ILike{"ci.siteURL": searchTerm},
		})
	}

	var sort string
	switch options.Sort {
	case app.SortByName, "":
		sort = "name"
	case app.SortByCSM:
		sort = "customerSuccessManager"
	case app.SortByAE:
		sort = "accountExecutive"
	case app.SortByTAM:
		sort = "technicalAccountManager"
	case app.SortByType:
		sort = "type"
	case app.SortBySiteURL:
		sort = "siteURL"
	case app.SortByLicensedTo:
		sort = "licensedTo"
	case app.SortByLastUpdated:
		sort = "lastUpdated"
	default:
		return sq.SelectBuilder{}, errors.Errorf("unsupported sort parameter '%s'", options.Sort)
	}

	var direction string
	switch options.Direction {
	case app.DirectionAsc:
		direction = "ASC"
	case app.DirectionDesc:
		direction = "DESC"
	case "":
		// Default to an ascending sort if none explicitly provided.
		direction = "ASC"
	default:
		return sq.SelectBuilder{}, errors.Errorf("unsupported direction parameter '%s'", options.Direction)
	}

	builder = builder.OrderByClause(fmt.Sprintf("%s %s", sort, direction))

	page := options.Page
	perPage := options.PerPage
	if page < 0 {
		page = 0
	}
	if perPage < 0 {
		perPage = 0
	}

	builder = builder.
		Offset(uint64(page * perPage)).
		Limit(uint64(perPage))

	return builder, nil
}

// NewCustomerStore creates a new store for customers ServiceImpl.
func NewCustomerStore(pluginAPI PluginAPIClient, sqlStore *SQLStore) app.CustomerStore {
	customerSelect := sqlStore.builder.
		Select(
			"ci.id", "ci.name", "ci.type",
			"ci.customerSuccessManager", "ci.accountExecutive",
			"ci.technicalAccountManager", "ci.salesforceId", "ci.zendeskId",
			"ci.siteUrl", "ci.licensedTo", "ci.gdriveLink", "ci.customerChannel", "ci.lastUpdated").
		From(customerTable + " as ci")

	packetValuesSelect := sqlStore.builder.
		Select(
			"cp.licensedTo", "cp.version", "cp.serverOS", "cp.serverArch",
			"cp.databaseType", "cp.databaseVersion", "cp.databaseSchemaVersion",
			"cp.fileDriver", "cp.activeUsers", "cp.dailyActiveUsers", "cp.monthlyActiveUsers",
			"cp.inactiveUserCount", "cp.licenseSupportedUsers", "cp.totalPosts", "cp.totalChannels", "cp.totalTeams").
		From(packetTable + " as cp")

	configValuesSelect := sqlStore.builder.
		Select("ccv.config").
		From(configTable + " as ccv")

	pluginValuesSelect := sqlStore.builder.
		Select("cpv.pluginId", "cpv.version", "cpv.isActive", "cpv.name").
		From(pluginTable + " as cpv")

	return &customerStore{
		pluginAPI:          pluginAPI,
		store:              sqlStore,
		queryBuilder:       sqlStore.builder,
		customerSelect:     customerSelect,
		packetValuesSelect: packetValuesSelect,
		configValuesSelect: configValuesSelect,
		pluginValuesSelect: pluginValuesSelect,
	}
}

func (s *customerStore) GetCustomers(opts app.CustomerFilterOptions) (app.GetCustomersResult, error) {
	queryForResults, err := applyCustomerFilterOptionsSort(s.customerSelect, opts)

	if err != nil {
		return app.GetCustomersResult{}, errors.Wrap(err, "failed to apply sort options")
	}

	queryForTotal := s.store.builder.
		Select("COUNT(*)").
		From(customerTable)

	if opts.SearchTerm != "" {
		queryForTotal = queryForTotal.Where(sq.Or{
			sq.ILike{"name": "%" + opts.SearchTerm + "%"},
			sq.ILike{"siteURL": "%" + opts.SearchTerm + "%"},
			sq.ILike{"licensedTo": opts.SearchTerm},
		})
	}

	var customers []app.Customer
	err = s.store.selectBuilder(s.store.db, &customers, queryForResults)

	if err == sql.ErrNoRows {
		return app.GetCustomersResult{}, errors.Wrap(app.ErrNotFound, "no customers found")
	} else if err != nil {
		return app.GetCustomersResult{}, errors.Wrap(err, "failed to get customers")
	}

	var total int

	if err = s.store.getBuilder(s.store.db, &total, queryForTotal); err != nil {
		return app.GetCustomersResult{}, errors.Wrap(err, "failed to get total customers")
	}

	pageCount := 0
	if opts.PerPage > 0 {
		pageCount = int(math.Ceil(float64(total) / float64(opts.PerPage)))
	}
	hasMore := opts.Page+1 < pageCount

	return app.GetCustomersResult{
		Customers:  customers,
		TotalCount: total,
		PageCount:  pageCount,
		HasMore:    hasMore,
	}, nil
}

func (s *customerStore) GetCustomerByID(id string) (app.FullCustomerInfo, error) {
	var customer app.FullCustomerInfo
	if id == "" {
		return app.FullCustomerInfo{}, errors.New("ID cannot be empty")
	}

	tx, err := s.store.db.Beginx()
	if err != nil {
		return app.FullCustomerInfo{}, errors.Wrap(err, "could not begin transaction")
	}
	defer s.store.finalizeTransaction(tx)
	var rawCustomers sqlCustomers
	err = s.store.getBuilder(tx, &rawCustomers, s.customerSelect.Where(sq.Eq{"ci.ID": id}))
	if err == sql.ErrNoRows {
		return app.FullCustomerInfo{}, errors.Wrapf(app.ErrNotFound, "customer does not exist for id '%s'", id)
	} else if err != nil {
		return app.FullCustomerInfo{}, errors.Wrapf(err, "failed to get customer by id '%s'", id)
	}

	if err = tx.Commit(); err != nil {
		return app.FullCustomerInfo{}, errors.Wrap(err, "could not commit transaction")
	}

	customer.Customer = rawCustomers.Customer

	config, err := s.GetConfig(customer.ID)
	if err != nil {
		return app.FullCustomerInfo{}, err
	}

	customer.Config = config

	plugins, err := s.GetPlugins(customer.ID)
	if err != nil {
		return app.FullCustomerInfo{}, err
	}
	customer.Plugins = plugins

	packet, err := s.GetPacket(id)
	if err != nil {
		return app.FullCustomerInfo{}, err
	}

	customer.PacketValues = packet

	return customer, nil
}

func (s *customerStore) createCustomer(siteURL string, licensedTo string) (string, error) {
	newID := model.NewId()

	_, err := s.store.execBuilder(s.store.db, sq.
		Insert(customerTable).
		// TODO - Should this use some kind of app.customer struct?
		// the one that I have now pulls the other attributes that should not be stored here.
		SetMap(map[string]interface{}{
			"ID":                      newID,
			"Name":                    licensedTo,
			"CustomerSuccessManager":  "",
			"AccountExecutive":        "",
			"TechnicalAccountManager": "",
			"SalesforceId":            "",
			"ZendeskId":               "",
			"LicensedTo":              licensedTo,
			"SiteUrl":                 siteURL,
			"Type":                    "",
			"GdriveLink":              "",
			"CustomerChannel":         "",
			"LastUpdated":             model.GetMillis(),
		}))
	if err != nil {
		return "", errors.Wrap(err, "failed to store new customer")
	}

	return newID, nil
}
func (s *customerStore) GetCustomerID(siteURL string, licensedTo string) (id string, err error) {
	if siteURL == "" || licensedTo == "" {
		return "", errors.New("must include siteURL or Licensedto")
	}

	tx, err := s.store.db.Beginx()
	if err != nil {
		return "", errors.Wrap(err, "could not begin transaction")
	}
	defer s.store.finalizeTransaction(tx)

	query := s.queryBuilder.
		Select("*").
		From(customerTable).
		Where(sq.Or{
			sq.Eq{"siteUrl": siteURL},
			sq.Eq{"licensedTo": licensedTo},
		})

	var rawCustomers []sqlCustomers

	err = s.store.selectBuilder(tx, &rawCustomers, query)

	if err != nil {
		return "", errors.Wrapf(err, "failed find customer with siteURL: '%s' and licensedTo: '%s'", siteURL, licensedTo)
	}

	if len(rawCustomers) == 0 {
		return s.createCustomer(siteURL, licensedTo)
	}

	if err = tx.Commit(); err != nil {
		return "", errors.Wrap(err, "could not commit transaction")
	}

	if len(rawCustomers) == 1 {
		return rawCustomers[0].ID, nil
	}

	var matchingLicense []app.Customer
	var matchingBoth []app.Customer

	for _, customer := range rawCustomers {
		if customer.LicensedTo == licensedTo && customer.SiteURL == siteURL {
			matchingBoth = append(matchingBoth, customer.Customer)
		}

		if customer.LicensedTo == licensedTo {
			matchingLicense = append(matchingLicense, customer.Customer)
		}
	}

	if len(matchingBoth) > 1 || len(matchingLicense) > 1 {
		// Need to evaluate this logic.
		return s.createCustomer(siteURL, licensedTo)
	} else if len(matchingBoth) == 1 || len(matchingLicense) == 1 {
		return matchingBoth[0].ID, nil
	}

	return "", errors.Wrapf(err, "No customer found with siteURL: '%s' and licensedTo: '%s'", siteURL, licensedTo)
}

func (s *customerStore) UpdateCustomer(customer app.Customer) error {
	if customer.ID == "" {
		return errors.New("customerID cannot be empty")
	}
	_, err := s.store.execBuilder(s.store.db, sq.
		Update(customerTable).
		SetMap(map[string]interface{}{
			"name":                    customer.Name,
			"customerSuccessManager":  customer.CustomerSuccessManager,
			"accountExecutive":        customer.AccountExecutive,
			"technicalAccountManager": customer.TechnicalAccountManager,
			"salesforceId":            customer.SalesforceID,
			"zendeskId":               customer.ZendeskID,
			"type":                    customer.Type,
			"customerChannel":         customer.CustomerChannel,
			"gdriveLink":              customer.GDriveLink,
			"lastUpdated":             model.GetMillis(),
		}).
		Where(sq.Eq{"id": customer.ID}))

	return err
}

func (s *customerStore) UpdateCustomerData(customerID string, userID string, packet *app.CustomerPacketValues, config *model.Config, plugins []app.CustomerPluginValues) error {
	if customerID == "" {
		return errors.New("customerID cannot be empty")
	}

	if packet == nil && config == nil && plugins == nil {
		return errors.New("must include at least one of packet, config, or plugins")
	}

	if packet != nil {
		err := s.storePacket(userID, customerID, packet)
		if err != nil {
			return errors.Wrap(err, "failed to store packet")
		}
	}

	if config != nil {
		err := s.storeConfig(userID, customerID, config)
		if err != nil {
			return errors.Wrap(err, "failed to store config")
		}
	}

	if plugins != nil {
		err := s.storePlugins(userID, customerID, plugins)
		if err != nil {
			return errors.Wrap(err, "failed to store plugins")
		}
	}

	return nil
}
