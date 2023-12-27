package sqlstore

import (
	"encoding/json"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/mattermost/mattermost/server/public/model"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"
)

// packetActionStore holds the information needed to fulfill the methods in the store interface.
type packetActionStore struct {
	pluginAPI          PluginAPIClient
	store              *SQLStore
	queryBuilder       sq.StatementBuilderType
	packetValuesSelect sq.SelectBuilder
	configValuesSelect sq.SelectBuilder
	pluginValuesSelect sq.SelectBuilder
}

type sqlPacket struct {
	app.CustomerPacketValues
}

type sqlConfig struct {
	Config json.RawMessage `db:"config"`
}

func NewPacketActionStore(pluginAPI PluginAPIClient, sqlStore *SQLStore) app.PacketActionStore {
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

	return &packetActionStore{
		pluginAPI:          pluginAPI,
		store:              sqlStore,
		queryBuilder:       sqlStore.builder,
		packetValuesSelect: packetValuesSelect,
		configValuesSelect: configValuesSelect,
		pluginValuesSelect: pluginValuesSelect,
	}
}

func (s *packetActionStore) GetPacket(customerId string) (app.CustomerPacketValues, error) {
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

func (s *packetActionStore) StorePacket(updateId string, packet app.CustomerPacketValues) error {

	return nil
}

func (s *packetActionStore) GetConfig(customerId string) (model.Config, error) {
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

func (s *packetActionStore) GetPlugins(customerId string) ([]app.CustomerPluginValues, error) {
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
