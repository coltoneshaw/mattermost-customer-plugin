package app

import (
	"github.com/coltoneshaw/mattermost-plugin-customers/server/bot"
	"github.com/mattermost/mattermost/server/public/model"
	pluginapi "github.com/mattermost/mattermost/server/public/pluginapi"
)

type customerService struct {
	store  CustomerStore
	poster bot.Poster
	api    *pluginapi.Client
}

// NewCustomerService returns a new customer service
func NewCustomerService(store CustomerStore, poster bot.Poster, api *pluginapi.Client) CustomerService {
	return &customerService{
		store:  store,
		poster: poster,
		api:    api,
	}
}

func (s *customerService) GetCustomerByID(id string) (Customer, error) {

	customer, err := s.store.GetCustomerByID(id)
	if err != nil {
		return Customer{}, err
	}

	config, err := s.GetConfig(customer.ID)
	if err != nil {
		return Customer{}, err
	}

	customer.Config = config

	plugins, err := s.GetPlugins(customer.ID)
	if err != nil {
		return Customer{}, err
	}
	customer.Plugins = plugins

	packet, err := s.GetPacket(id)
	if err != nil {
		return Customer{}, err
	}

	customer.PacketValues = packet

	return customer, nil
}

// func (s *customerService) GetId(siteUrl string, licensedTo string) (id string, err error) {

// }

func (p *customerService) GetPacket(customerId string) (CustomerPacketValues, error) {
	return p.store.GetPacket(customerId)
}

// func (p *customerService) StorePacket(updateId string, packet CustomerPacketValues) error {
// 	return p.store.StorePacket(updateId, packet)
// }

func (p *customerService) GetConfig(customerId string) (model.Config, error) {
	return p.store.GetConfig(customerId)
}

func (p *customerService) GetPlugins(customerId string) ([]CustomerPluginValues, error) {
	return p.store.GetPlugins(customerId)
}
