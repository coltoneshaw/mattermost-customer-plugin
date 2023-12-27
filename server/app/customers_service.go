package app

import (
	"github.com/coltoneshaw/mattermost-plugin-customers/server/bot"
	pluginapi "github.com/mattermost/mattermost/server/public/pluginapi"
)

type customerService struct {
	store  CustomerStore
	poster bot.Poster
	api    *pluginapi.Client
	packet PacketActionService
}

// NewCustomerService returns a new customer service
func NewCustomerService(store CustomerStore, poster bot.Poster, api *pluginapi.Client, packetGetter PacketActionService) CustomerService {
	return &customerService{
		store:  store,
		poster: poster,
		api:    api,
		packet: packetGetter,
	}
}

func (s *customerService) Get(id string) (Customer, error) {

	customer, err := s.store.Get(id)
	if err != nil {
		return Customer{}, err
	}

	config, err := s.packet.GetConfig(customer.ID)
	if err != nil {
		return Customer{}, err
	}

	customer.Config = config

	plugins, err := s.packet.GetPlugins(customer.ID)
	if err != nil {
		return Customer{}, err
	}
	customer.Plugins = plugins

	packet, err := s.packet.GetPacket(id)
	if err != nil {
		return Customer{}, err
	}

	customer.PacketValues = packet

	return customer, nil
}
