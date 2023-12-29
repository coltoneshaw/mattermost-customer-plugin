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

func applyCustomerFilterOptionsSort(builder sq.SelectBuilder, options app.CustomerFilterOptions) (sq.SelectBuilder, error) {
	var sort string
	switch options.Sort {
	case app.SortByName:
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
	case "":
		// Default to a stable sort if none explicitly provided.
		sort = "ID"
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
			"ci.siteUrl", "ci.licensedTo", "ci.gdriveLink", "ci.customerChannel").
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

	return rawCustomers.FullCustomerInfo, nil
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

func (s *customerStore) GetPacket(customerID string) (app.CustomerPacketValues, error) {
	if customerID == "" {
		return app.CustomerPacketValues{}, errors.New("ID cannot be empty")
	}

	tx, err := s.store.db.Beginx()
	if err != nil {
		return app.CustomerPacketValues{}, errors.Wrap(err, "could not begin transaction")
	}
	defer s.store.finalizeTransaction(tx)

	var rawPacket sqlPacket
	err = s.store.getBuilder(
		tx,
		&rawPacket,
		s.packetValuesSelect.
			Where(sq.Eq{"cp.customerId": customerID}).
			Where(sq.Eq{"cp.current": true}),
	)

	if err == sql.ErrNoRows {
		return app.CustomerPacketValues{}, nil
	} else if err != nil {
		return app.CustomerPacketValues{}, errors.Wrapf(err, "failed to get packet data for customer id '%s'", customerID)
	}

	if err = tx.Commit(); err != nil {
		return app.CustomerPacketValues{}, errors.Wrap(err, "could not commit transaction")
	}

	return rawPacket.CustomerPacketValues, nil
}

func (s *customerStore) GetConfig(customerID string) (model.Config, error) {
	if customerID == "" {
		return model.Config{}, errors.New("ID cannot be empty")
	}

	tx, err := s.store.db.Beginx()
	if err != nil {
		return model.Config{}, errors.Wrap(err, "could not begin transaction")
	}
	defer s.store.finalizeTransaction(tx)
	var rawConfig sqlConfig
	err = s.store.getBuilder(
		tx,
		&rawConfig,
		s.configValuesSelect.
			Where(sq.Eq{"ccv.customerId": customerID}).
			Where(sq.Eq{"ccv.current": true}),
	)

	if err == sql.ErrNoRows {
		return model.Config{}, nil
	} else if err != nil {
		return model.Config{}, errors.Wrapf(err, "failed to get config data for customer id '%s'", customerID)
	}

	var config model.Config
	err = json.Unmarshal(rawConfig.Config, &config)
	if err != nil {
		return model.Config{}, err
	}

	if err = tx.Commit(); err != nil {
		return model.Config{}, errors.Wrap(err, "could not commit transaction")
	}

	return config, nil
}

func (s *customerStore) GetPlugins(customerID string) ([]app.CustomerPluginValues, error) {
	if customerID == "" {
		return []app.CustomerPluginValues{}, errors.New("ID cannot be empty")
	}

	tx, err := s.store.db.Beginx()
	if err != nil {
		return []app.CustomerPluginValues{}, errors.Wrap(err, "could not begin transaction")
	}
	defer s.store.finalizeTransaction(tx)
	var rawPlugins []app.CustomerPluginValues
	err = s.store.selectBuilder(
		tx,
		&rawPlugins,
		s.pluginValuesSelect.
			Where(sq.Eq{"cpv.customerId": customerID}).
			Where(sq.Eq{"cpv.current": true}),
	)

	if err == sql.ErrNoRows {
		return []app.CustomerPluginValues{}, nil
	} else if err != nil {
		return []app.CustomerPluginValues{}, errors.Wrapf(err, "failed to get plugin data for customer id '%s'", customerID)
	}

	if err = tx.Commit(); err != nil {
		return []app.CustomerPluginValues{}, errors.Wrap(err, "could not commit transaction")
	}

	return rawPlugins, nil
}

func (s *customerStore) storePacket(updateID string, customerID string, packet *model.SupportPacket) error {
	_, err := s.store.execBuilder(s.store.db, sq.
		Update(packetTable).
		SetMap(map[string]interface{}{
			"current": false,
		}).
		Where(sq.Eq{"customerId": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to delete old packet data")
	}
	newID := model.NewId()
	_, err = s.store.execBuilder(s.store.db, sq.
		Insert(packetTable).
		SetMap(map[string]interface{}{
			"ID":                    newID,
			"customerId":            customerID,
			"updateDataId":          updateID,
			"licensedTo":            packet.LicenseTo,
			"version":               packet.ServerVersion,
			"serverOS":              packet.ServerOS,
			"serverArch":            packet.ServerArchitecture,
			"databaseType":          packet.DatabaseType,
			"databaseVersion":       packet.DatabaseVersion,
			"databaseSchemaVersion": packet.DatabaseSchemaVersion,
			"fileDriver":            packet.FileDriver,
			"activeUsers":           packet.ActiveUsers,
			"dailyActiveUsers":      packet.DailyActiveUsers,
			"monthlyActiveUsers":    packet.MonthlyActiveUsers,
			"inactiveUserCount":     packet.InactiveUserCount,
			"licenseSupportedUsers": packet.LicenseSupportedUsers,
			"totalPosts":            packet.TotalPosts,
			"totalChannels":         packet.TotalChannels,
			"totalTeams":            packet.TotalTeams,
			"current":               true,
		}))
	if err != nil {
		return errors.Wrap(err, "failed to store packet")
	}

	return nil
}

func (s *customerStore) storeConfig(updateID string, customerID string, config *model.Config) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	_, err = s.store.execBuilder(s.store.db, sq.
		Update(configTable).
		SetMap(map[string]interface{}{
			"current": false,
		}).
		Where(sq.Eq{"customerId": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to delete old config data")
	}

	_, err = s.store.execBuilder(s.store.db, sq.
		Insert(configTable).
		SetMap(map[string]interface{}{
			"ID":           model.NewId(),
			"customerId":   customerID,
			"updateDataId": updateID,
			"config":       string(configJSON),
			"current":      true,
		}))
	if err != nil {
		return errors.Wrap(err, "failed to store config")
	}

	return nil
}

func (s *customerStore) storePlugins(updateID string, customerID string, plugins *model.PluginsResponse) error {
	_, err := s.store.execBuilder(s.store.db, sq.
		Update(pluginTable).
		SetMap(map[string]interface{}{
			"current": false,
		}).
		Where(sq.Eq{"customerId": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to delete old plugin data")
	}

	var parsedPlugins []app.CustomerPluginValues

	for _, activePlugins := range plugins.Active {
		parsedPlugins = append(parsedPlugins, app.CustomerPluginValues{
			PluginID: activePlugins.Id,
			Version:  activePlugins.Version,
			IsActive: true,
			Name:     activePlugins.Name,
		})
	}

	for _, activePlugins := range plugins.Inactive {
		parsedPlugins = append(parsedPlugins, app.CustomerPluginValues{
			PluginID: activePlugins.Id,
			Version:  activePlugins.Version,
			IsActive: false,
			Name:     activePlugins.Name,
		})
	}

	for _, plugin := range parsedPlugins {
		_, err := s.store.execBuilder(s.store.db, sq.
			Insert(pluginTable).
			SetMap(map[string]interface{}{
				"ID":           model.NewId(),
				"customerId":   customerID,
				"updateDataId": updateID,
				"pluginId":     plugin.PluginID,
				"version":      plugin.Version,
				"isActive":     plugin.IsActive,
				"name":         plugin.Name,
				"current":      true,
			}))
		if err != nil {
			return errors.Wrap(err, "failed to store plugin")
		}
	}

	return nil
}

func (s *customerStore) createAuditRow(customerID string) (id string, err error) {
	id = model.NewId()
	_, err = s.store.execBuilder(s.store.db, sq.
		Insert("crm_audit").
		SetMap(map[string]interface{}{
			"ID":         id,
			"customerId": customerID,
			"updatedBy":  "",
			"updatedAt":  model.GetMillis(),
			"updateType": "packet",
			"path":       "",
		}))
	if err != nil {
		return "", errors.Wrap(err, "failed to store audit row")
	}

	return id, nil
}

func (s *customerStore) UpdateCustomerData(customerID string, packet *model.SupportPacket, config *model.Config, plugins *model.PluginsResponse) error {
	if customerID == "" {
		return errors.New("customerID cannot be empty")
	}

	auditID, err := s.createAuditRow(customerID)
	if err != nil {
		return errors.Wrap(err, "failed to create audit row")
	}

	err = s.storePacket(auditID, customerID, packet)
	if err != nil {
		return errors.Wrap(err, "failed to store packet")
	}

	err = s.storeConfig(auditID, customerID, config)
	if err != nil {
		return errors.Wrap(err, "failed to store config")
	}

	err = s.storePlugins(auditID, customerID, plugins)
	if err != nil {
		return errors.Wrap(err, "failed to store plugins")
	}

	return nil
}
