package sqlstore

import (
	"database/sql"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/mattermost/mattermost/server/public/model"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"
)

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

func (s *customerStore) storePlugins(userID string, customerID string, plugins []app.CustomerPluginValues) error {
	_, err := s.store.execBuilder(s.store.db, sq.
		Update(pluginTable).
		SetMap(map[string]interface{}{
			"current": false,
		}).
		Where(sq.Eq{"customerId": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to delete old plugin data")
	}

	existingPlugins, err := s.GetPlugins(customerID)
	if err != nil {
		return errors.Wrap(err, "failed to get existing plugins")
	}

	diff, err := diffPlugins(existingPlugins, plugins)

	if err != nil {
		return errors.Wrap(err, "failed to diff plugins")
	}

	auditID, err := s.createAuditRow(customerID, userID, diff)
	if err != nil {
		return errors.Wrap(err, "failed to create audit row")
	}

	for _, plugin := range plugins {
		_, err := s.store.execBuilder(s.store.db, sq.
			Insert(pluginTable).
			SetMap(map[string]interface{}{
				"ID":          model.NewId(),
				"AuditID":     auditID,
				"CustomerID":  customerID,
				"Current":     true,
				"PluginID":    plugin.PluginID,
				"Version":     plugin.Version,
				"IsActive":    plugin.IsActive,
				"Name":        plugin.Name,
				"HomePageURL": plugin.HomePageURL,
			}))
		if err != nil {
			return errors.Wrap(err, "failed to store plugin")
		}
	}

	return nil
}
