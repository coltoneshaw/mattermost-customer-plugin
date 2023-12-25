package app

import (
	"github.com/coltoneshaw/mattermost-plugin-customers/server/bot"
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
		api:    api}
}

func (s *customerService) Get(id string) (Customer, error) {
	return s.store.Get(id)
}
