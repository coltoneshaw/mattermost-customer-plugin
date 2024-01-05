package sqlstore

import (
	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"
)

func (s *customerStore) rawPacketToPacket(rawPacket *model.SupportPacket) *app.CustomerPacketValues {
	var packet app.CustomerPacketValues

	packet.LicensedTo = rawPacket.LicenseTo
	packet.Version = rawPacket.ServerVersion
	packet.ServerOS = rawPacket.ServerOS
	packet.ServerArch = rawPacket.ServerArchitecture
	packet.DatabaseType = rawPacket.DatabaseType
	packet.DatabaseVersion = rawPacket.DatabaseVersion
	packet.DatabaseSchemaVersion = rawPacket.DatabaseSchemaVersion
	packet.FileDriver = rawPacket.FileDriver
	packet.ActiveUsers = rawPacket.ActiveUsers
	packet.DailyActiveUsers = rawPacket.DailyActiveUsers
	packet.MonthlyActiveUsers = rawPacket.MonthlyActiveUsers
	packet.InactiveUserCount = rawPacket.InactiveUserCount
	packet.LicenseSupportedUsers = rawPacket.LicenseSupportedUsers
	packet.TotalPosts = rawPacket.TotalPosts
	packet.TotalChannels = rawPacket.TotalChannels
	packet.TotalTeams = rawPacket.TotalTeams

	return &packet
}

func (s *customerStore) rawPluginstoPlugins(rawPlugins *model.PluginsResponse) []app.CustomerPluginValues {
	var parsedPlugins []app.CustomerPluginValues

	for _, activePlugins := range rawPlugins.Active {
		parsedPlugins = append(parsedPlugins, app.CustomerPluginValues{
			PluginID: activePlugins.Id,
			Version:  activePlugins.Version,
			IsActive: true,
			Name:     activePlugins.Name,
		})
	}

	for _, activePlugins := range rawPlugins.Inactive {
		parsedPlugins = append(parsedPlugins, app.CustomerPluginValues{
			PluginID: activePlugins.Id,
			Version:  activePlugins.Version,
			IsActive: false,
			Name:     activePlugins.Name,
		})
	}

	return parsedPlugins
}

func (s *customerStore) UpdateCustomerThroughUpload(customerID string, packet *model.SupportPacket, config *model.Config, plugins *model.PluginsResponse) error {
	if customerID == "" {
		return errors.New("customerID cannot be empty")
	}

	if packet == nil && config == nil && plugins == nil {
		return errors.New("must include at least one of packet, config, or plugins")
	}

	if packet != nil {
		rawPacket := s.rawPacketToPacket(packet)

		err := s.storePacket("", customerID, rawPacket)
		if err != nil {
			return errors.Wrap(err, "failed to store packet")
		}
	}

	if config != nil {
		err := s.storeConfig("", customerID, config)
		if err != nil {
			return errors.Wrap(err, "failed to store config")
		}
	}

	if plugins != nil {
		rawPlugins := s.rawPluginstoPlugins(plugins)
		err := s.storePlugins("", customerID, rawPlugins)
		if err != nil {
			return errors.Wrap(err, "failed to store plugins")
		}
	}

	return nil
}
