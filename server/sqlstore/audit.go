package sqlstore

import (
	"encoding/json"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/mattermost/mattermost/server/public/model"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/r3labs/diff"
)

func diffConfig(old, new *model.Config) (diff.Changelog, error) {
	return diff.Diff(old, new)
}

func diffPacket(old, new *app.CustomerPacketValues) (diff.Changelog, error) {
	return diff.Diff(old, new)
}

func diffPlugins(old, new []app.CustomerPluginValues) (diff.Changelog, error) {
	return diff.Diff(old, new)
}

func (s *customerStore) createAuditRow(customerID string, updatedBy string, diff diff.Changelog) (id string, err error) {
	if customerID == "" {
		return "", errors.New("customerID cannot be empty")
	}

	var updateType UpdateType
	if updatedBy != "" {
		updateType = User
	} else {
		updateType = Packet
	}

	changelogJSON, err := json.Marshal(diff)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal changelog")
	}

	id = model.NewId()
	lastUpdated := model.GetMillis()
	_, err = s.store.execBuilder(s.store.db, sq.
		Insert(auditTable).
		SetMap(map[string]interface{}{
			"ID":         id,
			"customerId": customerID,
			"updatedBy":  updatedBy,
			"updatedAt":  lastUpdated,
			"updateType": updateType,
			"path":       "",
			"diff":       string(changelogJSON),
		}))
	if err != nil {
		return "", errors.Wrap(err, "failed to store audit row")
	}

	_, err = s.store.execBuilder(s.store.db, sq.
		Update(customerTable).
		SetMap(map[string]interface{}{
			"lastUpdated": lastUpdated,
		}).
		Where(sq.Eq{"ID": customerID}),
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to store audit row")
	}
	return id, nil
}
