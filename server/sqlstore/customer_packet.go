package sqlstore

import (
	"database/sql"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/mattermost/mattermost/server/public/model"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"
)

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

func (s *customerStore) storePacket(userID string, customerID string, packet *app.CustomerPacketValues) error {
	existingPacket, err := s.GetPacket(customerID)

	if err != nil {
		return errors.Wrap(err, "failed to get existing packet")
	}

	_, err = s.store.execBuilder(s.store.db, sq.
		Update(packetTable).
		SetMap(map[string]interface{}{
			"current": false,
		}).
		Where(sq.Eq{"customerId": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to delete old packet data")
	}

	diff, err := diffPacket(&existingPacket, packet)

	if err != nil {
		return errors.Wrap(err, "failed to diff packet")
	}

	auditID, err := s.createAuditRow(customerID, userID, diff)
	if err != nil {
		return errors.Wrap(err, "failed to create audit row")
	}

	newID := model.NewId()
	_, err = s.store.execBuilder(s.store.db, sq.
		Insert(packetTable).
		SetMap(map[string]interface{}{
			"ID":                    newID,
			"customerId":            customerID,
			"updateDataId":          auditID,
			"licensedTo":            packet.LicensedTo,
			"version":               packet.Version,
			"serverOS":              packet.ServerOS,
			"serverArch":            packet.ServerArch,
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

	// updating licensedTo in the customer table to always keep it up to date
	_, err = s.store.execBuilder(s.store.db, sq.
		Update(customerTable).
		SetMap(map[string]interface{}{
			"licensedTo": packet.LicensedTo,
		}).
		Where(sq.Eq{"id": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to update licensedTo from packet update")
	}

	return nil
}
