package sqlstore

import (
	"database/sql"
	"encoding/json"

	"github.com/mattermost/mattermost/server/public/model"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"
)

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

func (s *customerStore) storeConfig(userID string, customerID string, config *model.Config) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	existingConfig, err := s.GetConfig(customerID)

	if err != nil {
		return errors.Wrap(err, "failed to get existing config")
	}

	diff, err := diffConfig(&existingConfig, config)

	if err != nil {
		return errors.Wrap(err, "failed to diff config")
	}

	auditID, err := s.createAuditRow(customerID, userID, diff)
	if err != nil {
		return errors.Wrap(err, "failed to create audit row")
	}

	_, err = s.store.execBuilder(s.store.db, sq.
		Update(configTable).
		SetMap(map[string]interface{}{
			"current": false,
		}).
		Where(sq.Eq{"customerId": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to set old config inactive")
	}

	// taking the site url from the config and storing it in the customer table for matching later.
	_, err = s.store.execBuilder(s.store.db, sq.
		Update(customerTable).
		SetMap(map[string]interface{}{
			"siteURL": config.ServiceSettings.SiteURL,
		}).
		Where(sq.Eq{"id": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to update siteURL from config change")
	}

	_, err = s.store.execBuilder(s.store.db, sq.
		Insert(configTable).
		SetMap(map[string]interface{}{
			"ID":         model.NewId(),
			"AuditID":    auditID,
			"Current":    true,
			"CustomerId": customerID,
			"Config":     string(configJSON),
		}))
	if err != nil {
		return errors.Wrap(err, "failed to store config")
	}

	return nil
}
