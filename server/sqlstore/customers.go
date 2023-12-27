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

// NewCustomerStore creates a new store for customers ServiceImpl.
func NewCustomerStore(pluginAPI PluginAPIClient, sqlStore *SQLStore) app.CustomerStore {
	customerSelect := sqlStore.builder.
		Select(
			"ci.id", "ci.name", "ci.type",
			"ci.customerSuccessManager", "ci.accountExecutive",
			"ci.technicalAccountManager", "ci.salesforceId", "ci.zendeskId",
			"ci.siteUrl", "ci.licensedTo").
		From("crm_customers as ci")

	packetValuesSelect := sqlStore.builder.
		Select(
			"cpv.licensedTo", "cpv.version", "cpv.serverOS", "cpv.serverArch",
			"cpv.databaseType", "cpv.databaseVersion", "cpv.databaseSchemaVersion",
			"cpv.fileDriver", "cpv.activeUsers", "cpv.dailyActiveUsers", "cpv.monthlyActiveUsers",
			"cpv.inactiveUserCount", "cpv.licenseSupportedUsers", "cpv.totalPosts", "cpv.totalChannels", "cpv.totalTeams").
		From("crm_packetValues as cpv")

	configValuesSelect := sqlStore.builder.
		Select("ccv.config").
		From("crm_configValues as ccv")

	pluginValuesSelect := sqlStore.builder.
		Select("cpv.pluginId", "cpv.version", "cpv.isActive", "cpv.name").
		From("crm_pluginValues as cpv")

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

// func (s *customerStore) GetId(siteUrl string, licensedTo string) (id string, err error) {
// 	if siteUrl == "" || licensedTo == "" {
// 		return "", errors.New("must include siteUrl or Licensedto")
// 	}

// 	tx, err := s.store.db.Beginx()
// 	if err != nil {
// 		return "", errors.Wrap(err, "could not begin transaction")
// 	}
// 	defer s.store.finalizeTransaction(tx)

// 	query := s.queryBuilder.
// 		Select("id").
// 		From("crm_customers").
// 		Where(sq.Or{
// 			sq.Eq{"siteUrl": siteUrl},
// 			sq.Eq{"licensedTo": licensedTo},
// 		})

// 	var rawIds []sqlCustomers

// 	err = s.store.selectBuilder(tx, &rawCustomers, s.customerSelect.Where(sq.Eq{"ci.ID": id}))
// 	if err == sql.ErrNoRows {
// 		return app.Customer{}, errors.Wrapf(app.ErrNotFound, "customer does not exist for id '%s'", id)
// 	} else if err != nil {
// 		return app.Customer{}, errors.Wrapf(err, "failed to get customer by id '%s'", id)
// 	}
// }

func (s *customerStore) GetPacket(customerId string) (app.CustomerPacketValues, error) {
	if customerId == "" {
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
			Where(sq.Eq{"cpv.customerId": customerId}).
			Where(sq.Eq{"cpv.current": true}),
	)

	// TODO - this length could possibly be > 1. How to ensure it's always 1.
	if err != nil {
		return app.CustomerPacketValues{}, errors.Wrapf(err, "failed to get packet data for customer id '%s'", customerId)
	}

	if err = tx.Commit(); err != nil {
		return app.CustomerPacketValues{}, errors.Wrap(err, "could not commit transaction")
	}

	return rawPacket.CustomerPacketValues, nil
}

func (s *customerStore) GetConfig(customerId string) (model.Config, error) {
	if customerId == "" {
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
			Where(sq.Eq{"ccv.customerId": customerId}).
			Where(sq.Eq{"ccv.current": true}),
	)

	// TODO - this length could possibly be > 1. How to ensure it's always 1.
	if err != nil {
		return model.Config{}, errors.Wrapf(err, "failed to get config data for customer id '%s'", customerId)
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

func (s *customerStore) GetPlugins(customerId string) ([]app.CustomerPluginValues, error) {
	if customerId == "" {
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
			Where(sq.Eq{"cpv.customerId": customerId}).
			Where(sq.Eq{"cpv.current": true}),
	)
	if err != nil {
		return []app.CustomerPluginValues{}, errors.Wrapf(err, "failed to get plugin data for customer id '%s'", customerId)
	}

	if err = tx.Commit(); err != nil {
		return []app.CustomerPluginValues{}, errors.Wrap(err, "could not commit transaction")
	}

	return rawPlugins, nil
}
