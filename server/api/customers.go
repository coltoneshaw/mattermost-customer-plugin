package api

import (
	"errors"
	"net/http"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/coltoneshaw/mattermost-plugin-customers/server/config"
	"github.com/gorilla/mux"
	pluginapi "github.com/mattermost/mattermost/server/public/pluginapi"
)

// CusotmerHandler is the API handler.
type CustomerHandler struct {
	*ErrorHandler
	customerService app.CustomerService
	pluginAPI       *pluginapi.Client
	config          config.Service
}

// NewCustomerHandler returns a new customer api handler
func NewCustomerHandler(router *mux.Router, customerService app.CustomerService, api *pluginapi.Client, configService config.Service) *CustomerHandler {
	handler := &CustomerHandler{
		ErrorHandler:    &ErrorHandler{},
		customerService: customerService,
		pluginAPI:       api,
		config:          configService,
	}

	customersRouter := router.PathPrefix("/customers").Subrouter()

	// customerRouter.HandleFunc("", withContext(handler.createCustomer)).Methods(http.MethodPost)

	//
	customerRouter := customersRouter.PathPrefix("/{id:[A-Za-z0-9]+}").Subrouter()
	customerRouter.HandleFunc("", withContext(handler.getCustomer)).Methods(http.MethodGet)

	return handler
}

func (h *CustomerHandler) getCustomer(c *Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["id"]

	customer, err := h.customerService.GetCustomerByID(customerID)

	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			h.HandleErrorWithCode(w, c.logger, http.StatusNotFound, "No customer found for this ID", err)
			return
		}
		h.HandleError(w, c.logger, err)
		return
	}

	ReturnJSON(w, &customer, http.StatusOK)
}
