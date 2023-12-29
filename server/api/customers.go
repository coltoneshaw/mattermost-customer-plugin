package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"

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
	customersRouter.HandleFunc("", withContext(handler.getCustomers)).Methods(http.MethodGet)

	//
	customerRouter := customersRouter.PathPrefix("/{id:[A-Za-z0-9]+}").Subrouter()
	customerRouter.HandleFunc("", withContext(handler.getCustomer)).Methods(http.MethodGet)

	return handler
}

func (h *CustomerHandler) getCustomers(c *Context, w http.ResponseWriter, r *http.Request) {
	opts, err := parseGetPlaybooksOptions(r.URL)
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

func parseGetPlaybooksOptions(u *url.URL) (app.CustomerFilterOptions, error) {
	params := u.Query()

	var sortField app.SortField
	param := strings.ToLower(params.Get("sort"))
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
	case "siteURL":
		sortField = app.SortBySiteURL
	case "licensedTo":
		sortField = app.SortByLicensedTo
	default:
		return app.CustomerFilterOptions{}, errors.Errorf("bad parameter 'sort' (%s): it should be empty or one of 'name', 'customerSuccessManager', 'accountExecutive', 'technicalAccountManager', 'type', 'siteURL', 'licensedTo'", param)
	}

	var sortDirection app.SortDirection
	param = strings.ToLower(params.Get("direction"))
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

	perPageParam := params.Get("per_page")
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
		Sort:      sortField,
		Direction: sortDirection,
		Page:      page,
		PerPage:   perPage,
	}, nil
}
