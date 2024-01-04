package sqlstore

import (
	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/mattermost/mattermost/server/public/model"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"
)

func (s *customerStore) uploadStorePacket(updateID string, customerID string, packet *model.SupportPacket) error {
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

	// updating licensedTo in the customer table to always keep it up to date
	_, err = s.store.execBuilder(s.store.db, sq.
		Update(customerTable).
		SetMap(map[string]interface{}{
			"licensedTo": packet.LicenseTo,
		}).
		Where(sq.Eq{"id": customerID}))

	if err != nil {
		return errors.Wrap(err, "failed to update licensedTo from packet update")
	}

	return nil
}

func (s *customerStore) uploadStorePlugins(updateID string, customerID string, plugins *model.PluginsResponse) error {
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

func (s *customerStore) UpdateCustomerThroughUpload(customerID string, packet *model.SupportPacket, config *model.Config, plugins *model.PluginsResponse) error {
	if customerID == "" {
		return errors.New("customerID cannot be empty")
	}

	if packet == nil && config == nil && plugins == nil {
		return errors.New("must include at least one of packet, config, or plugins")
	}

	auditID, err := s.createAuditRow(customerID, "")
	if err != nil {
		return errors.Wrap(err, "failed to create audit row")
	}

	if packet != nil {
		err = s.uploadStorePacket(auditID, customerID, packet)
		if err != nil {
			return errors.Wrap(err, "failed to store packet")
		}
	}

	if config != nil {
		err = s.storeConfig(auditID, customerID, config)
		if err != nil {
			return errors.Wrap(err, "failed to store config")
		}
	}

	if plugins != nil {
		err = s.uploadStorePlugins(auditID, customerID, plugins)
		if err != nil {
			return errors.Wrap(err, "failed to store plugins")
		}
	}

	return nil
}
