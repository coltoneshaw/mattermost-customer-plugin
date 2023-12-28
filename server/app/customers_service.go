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

func (s *customerService) GetCustomerID(siteURL string, licensedTo string) (id string, err error) {
	return s.store.GetCustomerID(siteURL, licensedTo)
}

func (s *customerService) GetPacket(customerID string) (CustomerPacketValues, error) {
	return s.store.GetPacket(customerID)
}

func (s *customerService) GetConfig(customerID string) (model.Config, error) {
	return s.store.GetConfig(customerID)
}

func (s *customerService) GetPlugins(customerID string) ([]CustomerPluginValues, error) {
	return s.store.GetPlugins(customerID)
}
