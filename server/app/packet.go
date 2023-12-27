package app

import (
	"github.com/mattermost/mattermost/server/public/model"
)

type PacketActionService interface {

	// MessageHasBeenPosted suggests playbooks to the user if triggered
	MessageHasBeenPosted(post *model.Post)

	GetPacket(customerId string) (CustomerPacketValues, error)
	GetConfig(customerId string) (model.Config, error)
	GetPlugins(customerId string) ([]CustomerPluginValues, error)
}

type PacketActionStore interface {
	GetPacket(customerId string) (CustomerPacketValues, error)
	GetConfig(customerId string) (model.Config, error)
	GetPlugins(customerId string) ([]CustomerPluginValues, error)
}
