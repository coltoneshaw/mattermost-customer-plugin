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

func (s *customerService) GetCustomers(opts CustomerFilterOptions) (GetCustomersResult, error) {
	return s.store.GetCustomers(opts)
}

func (s *customerService) UpdateCustomer(customer Customer, userID string) error {
	return s.store.UpdateCustomer(customer, userID)
}

func (s *customerService) UpdateCustomerData(customerID string, userID string, packet *CustomerPacketValues, config *model.Config, plugins []CustomerPluginValues) error {
	return s.store.UpdateCustomerData(customerID, userID, packet, config, plugins)
}

func (s *customerService) UpdateCustomerThroughUpload(customerID string, packet *model.SupportPacket, config *model.Config, plugins *model.PluginsResponse) error {
	return s.store.UpdateCustomerThroughUpload(customerID, packet, config, plugins)
}

func (s *customerService) GetCustomerByID(id string) (FullCustomerInfo, error) {
	return s.store.GetCustomerByID(id)
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
