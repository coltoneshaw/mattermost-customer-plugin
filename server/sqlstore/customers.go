package sqlstore

import (
	"database/sql"
	"encoding/json"

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
	app.Customer
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
)

// NewCustomerStore creates a new store for customers ServiceImpl.
func NewCustomerStore(pluginAPI PluginAPIClient, sqlStore *SQLStore) app.CustomerStore {
	customerSelect := sqlStore.builder.
		Select(
			"ci.id", "ci.name", "ci.type",
			"ci.customerSuccessManager", "ci.accountExecutive",
			"ci.technicalAccountManager", "ci.salesforceId", "ci.zendeskId",
			"ci.siteUrl", "ci.licensedTo").
		From(customerTable + " as ci")

	packetValuesSelect := sqlStore.builder.
		Select(
			"cpv.licensedTo", "cpv.version", "cpv.serverOS", "cpv.serverArch",
			"cpv.databaseType", "cpv.databaseVersion", "cpv.databaseSchemaVersion",
			"cpv.fileDriver", "cpv.activeUsers", "cpv.dailyActiveUsers", "cpv.monthlyActiveUsers",
			"cpv.inactiveUserCount", "cpv.licenseSupportedUsers", "cpv.totalPosts", "cpv.totalChannels", "cpv.totalTeams").
		From(packetTable + " as cpv")

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

// get customer data
// get the latest from crm_packetValues with current = true
// this should just be one

// get the latest from config where current = true
// this should just be one

// get the latest from pluginValues where current = true
// this can be an array of values

func (s *customerStore) GetCustomerByID(id string) (app.Customer, error) {
	if id == "" {
		return app.Customer{}, errors.New("ID cannot be empty")
	}

	tx, err := s.store.db.Beginx()
	if err != nil {
		return app.Customer{}, errors.Wrap(err, "could not begin transaction")
	}
	defer s.store.finalizeTransaction(tx)
	var rawCustomers sqlCustomers
	err = s.store.getBuilder(tx, &rawCustomers, s.customerSelect.Where(sq.Eq{"ci.ID": id}))
	if err == sql.ErrNoRows {
		return app.Customer{}, errors.Wrapf(app.ErrNotFound, "customer does not exist for id '%s'", id)
	} else if err != nil {
		return app.Customer{}, errors.Wrapf(err, "failed to get customer by id '%s'", id)
	}

	if err = tx.Commit(); err != nil {
		return app.Customer{}, errors.Wrap(err, "could not commit transaction")
	}

	return rawCustomers.Customer, nil
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
			"LicensedTo":              "",
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
			Where(sq.Eq{"cpv.customerId": customerID}).
			Where(sq.Eq{"cpv.current": true}),
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
