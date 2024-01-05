package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/coltoneshaw/mattermost-plugin-customers/server/config"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
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
	customersRouter.HandleFunc("", withContext(handler.getCustomers)).Methods(http.MethodGet)

	//
	customerRouter := customersRouter.PathPrefix("/{id:[A-Za-z0-9]+}").Subrouter()
	customerRouter.HandleFunc("", withContext(handler.getCustomer)).Methods(http.MethodGet)
	customerRouter.HandleFunc("", withContext(handler.updateCustomer)).Methods(http.MethodPut)

	configRouter := customerRouter.PathPrefix("/config").Subrouter()
	configRouter.HandleFunc("", withContext(handler.updateCustomerConfig)).Methods(http.MethodPut)

	packetRouter := customerRouter.PathPrefix("/packet").Subrouter()
	packetRouter.HandleFunc("", withContext(handler.updateCustomerPacket)).Methods(http.MethodPut)

	pluginRouter := customerRouter.PathPrefix("/plugins").Subrouter()
	pluginRouter.HandleFunc("", withContext(handler.updateCustomerPlugins)).Methods(http.MethodPut)

	return handler
}

func (h *CustomerHandler) getCustomers(c *Context, w http.ResponseWriter, r *http.Request) {
	opts, err := parseGetCustomerOptions(r.URL)
	if err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, fmt.Sprintf("failed to get customers: %s", err.Error()), nil)
		return
	}

	customerResults, err := h.customerService.GetCustomers(opts)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}
	ReturnJSON(w, customerResults, http.StatusOK)
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

func (h *CustomerHandler) updateCustomer(c *Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// userID := r.Header.Get("Mattermost-User-ID")
	var customer app.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to decode customer", err)
		return
	}

	customer.ID = vars["id"]
	err := h.customerService.UpdateCustomer(customer)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}

	fullCustomer, err := h.customerService.GetCustomerByID(customer.ID)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}

	ReturnJSON(w, &fullCustomer, http.StatusOK)
}

func (h *CustomerHandler) updateCustomerConfig(c *Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.Header.Get("Mattermost-User-ID")
	var config model.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to decode customer config", err)
		return
	}

	// need to validate config here. When i attempted to do it I keep getting a nil pointer dereference, with valid and non valid config.
	// if err := config.IsValid(); err != nil {
	// 	h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "invalid config", err)
	// 	return
	// }

	customerID := vars["id"]
	err := h.customerService.UpdateCustomerData(customerID, userID, nil, &config, nil)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}

	fullCustomer, err := h.customerService.GetCustomerByID(customerID)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}

	ReturnJSON(w, &fullCustomer, http.StatusOK)
}

func (h *CustomerHandler) updateCustomerPacket(c *Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.Header.Get("Mattermost-User-ID")
	var packet app.CustomerPacketValues
	if err := json.NewDecoder(r.Body).Decode(&packet); err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to decode customer packet info", err)
		return
	}

	customerID := vars["id"]
	err := h.customerService.UpdateCustomerData(customerID, userID, &packet, nil, nil)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}

	fullCustomer, err := h.customerService.GetCustomerByID(customerID)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}

	ReturnJSON(w, &fullCustomer, http.StatusOK)
}

func (h *CustomerHandler) updateCustomerPlugins(c *Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.Header.Get("Mattermost-User-ID")
	var plugins []app.CustomerPluginValues
	if err := json.NewDecoder(r.Body).Decode(&plugins); err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to decode customer plugins", err)
		return
	}

	customerID := vars["id"]
	err := h.customerService.UpdateCustomerData(customerID, userID, nil, nil, plugins)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}

	fullCustomer, err := h.customerService.GetCustomerByID(customerID)
	if err != nil {
		h.HandleError(w, c.logger, err)
		return
	}

	ReturnJSON(w, &fullCustomer, http.StatusOK)
}

func parseGetCustomerOptions(u *url.URL) (app.CustomerFilterOptions, error) {
	params := u.Query()

	var searchTerm string
	param := strings.ToLower(params.Get("searchTerm"))
	if param != "" {
		searchTerm = param
	}

	var sortField app.SortField
	param = strings.ToLower(params.Get("sort"))
	switch param {
	case "name", "":
		sortField = app.SortByName
	case "csm":
		sortField = app.SortByCSM
	case "ae":
		sortField = app.SortByAE
	case "tam":
		sortField = app.SortByTAM
	case "type":
		sortField = app.SortByType
	case "site_url":
		sortField = app.SortBySiteURL
	case "licensed_to":
		sortField = app.SortByLicensedTo
	case "last_updated":
		sortField = app.SortByLicensedTo
	default:
		return app.CustomerFilterOptions{}, errors.Errorf("bad parameter 'sort' (%s): it should be empty or one of 'name', 'customerSuccessManager', 'accountExecutive', 'technicalAccountManager', 'type', 'siteURL', 'licensedTo'", param)
	}

	var sortDirection app.SortDirection
	param = strings.ToLower(params.Get("order"))
	switch param {
	case "asc", "":
		sortDirection = app.DirectionAsc
	case "desc":
		sortDirection = app.DirectionDesc
	default:
		return app.CustomerFilterOptions{}, errors.Errorf("bad parameter 'direction' (%s): it should be empty or one of 'asc' or 'desc'", param)
	}

	pageParam := params.Get("page")
	if pageParam == "" {
		pageParam = "0"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		return app.CustomerFilterOptions{}, errors.Wrapf(err, "bad parameter 'page': it should be a number")
	}
	if page < 0 {
		return app.CustomerFilterOptions{}, errors.Errorf("bad parameter 'page': it should be a positive number")
	}

	perPageParam := params.Get("perPage")
	if perPageParam == "" || perPageParam == "0" {
		perPageParam = "1000"
	}
	perPage, err := strconv.Atoi(perPageParam)
	if err != nil {
		return app.CustomerFilterOptions{}, errors.Wrapf(err, "bad parameter 'per_page': it should be a number")
	}
	if perPage < 0 {
		return app.CustomerFilterOptions{}, errors.Errorf("bad parameter 'per_page': it should be a positive number")
	}

	return app.CustomerFilterOptions{
		Sort:       sortField,
		Direction:  sortDirection,
		SearchTerm: searchTerm,
		Page:       page,
		PerPage:    perPage,
	}, nil
}
