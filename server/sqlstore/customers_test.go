package sqlstore

import (
	"reflect"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	mock_sqlstore "github.com/coltoneshaw/mattermost-plugin-customers/server/sqlstore/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
)

func setupCustomerStore(t *testing.T, db *sqlx.DB) app.CustomerStore {
	mockCtrl := gomock.NewController(t)

	configAPI := mock_sqlstore.NewMockConfigurationAPI(mockCtrl)
	pluginAPIClient := PluginAPIClient{
		Configuration: configAPI,
	}

	sqlStore := setupSQLStore(t, db)

	return NewCustomerStore(pluginAPIClient, sqlStore)
}

func TestGetCustomerByID(t *testing.T) {
	db := setupTestDB(t)
	customerStore := setupCustomerStore(t, db)

	t.Run("fail to get customer id", func(t *testing.T) {
		_, err := customerStore.GetCustomerByID("1")
		if errors.Cause(err) != app.ErrNotFound {
			t.Fatal(err)
		}
	})

	t.Run("no id provided", func(t *testing.T) {
		_, err := customerStore.GetCustomerByID("")
		if err.Error() != "ID cannot be empty" {
			t.Fatal(err)
		}
	})

	t.Run("get customer id", func(t *testing.T) {
		_, err := db.Exec(`INSERT INTO crm_customers (id, name, type, lastUpdated) VALUES ($1, $2, $3, $4)`, "1", "test", "cloud", 0)
		if err != nil {
			t.Fatal(err)
		}

		customer, err := customerStore.GetCustomerByID("1")
		if err != nil {
			t.Fatal(err)
		}

		if customer.ID != "1" {
			t.Fatal("customer id does not match")
		}
	})
}

func TestGetCustomerID(t *testing.T) {
	db := setupTestDB(t)
	customerStore := setupCustomerStore(t, db)

	t.Run("no siteurl or licensedto provided", func(t *testing.T) {
		_, err := customerStore.GetCustomerID("", "")
		if err.Error() != "must include siteURL or Licensedto" {
			t.Fatal(err)
		}
	})

	t.Run("creates new profile", func(t *testing.T) {
		ID, err := customerStore.GetCustomerID("www.test.com", "test")
		if err != nil {
			t.Fatal(err)
		}

		if len(ID) != 26 {
			t.Fatal("does not appear to be a valid id")
		}

		customer, err := customerStore.GetCustomerByID(ID)

		if err != nil {
			t.Fatal(err)
		}

		if customer.ID != ID {
			t.Fatal("customer id does not match")
		}
	})

	t.Run("returns existing profile", func(t *testing.T) {
		_, err := db.Exec(`INSERT INTO crm_customers (id, name, type, siteurl, licensedto, lastupdated) VALUES ($1, $2, $3, $4, $5, $6)`, "1", "test", "cloud", "www.1.com", "1", 0)
		if err != nil {
			t.Fatal(err)
		}

		ID, err := customerStore.GetCustomerID("www.1.com", "1")
		if err != nil {
			t.Fatal(err)
		}

		if ID != "1" {
			t.Fatal("customer id does not match")
		}
	})

	t.Run("creates new due to too many exact matches", func(t *testing.T) {
		_, err := db.Exec(`INSERT INTO crm_customers (id, name, type, siteurl, licensedto, lastupdated) VALUES ($1, $2, $3, $4, $5, $6)`, "2", "test", "cloud", "www.1.com", "1", 0)
		if err != nil {
			t.Fatal(err)
		}

		ID, err := customerStore.GetCustomerID("www.1.com", "1")
		if err != nil {
			t.Fatal(err)
		}

		if ID == "1" || ID == "2" {
			t.Fatal("customer id does not match")
		}

		if len(ID) != 26 {
			t.Fatal("does not appear to be a valid id")
		}
	})

	t.Run("siteurl was changed licensedto remains", func(t *testing.T) {
		ID, err := customerStore.GetCustomerID("www.2.com", "1")
		if err != nil {
			t.Fatal(err)
		}

		if ID == "1" || ID == "2" {
			t.Fatal("customer id does not match")
		}

		if len(ID) != 26 {
			t.Fatal("does not appear to be a valid id")
		}
	})
}

func TestGetCustomers(t *testing.T) {
	db := setupTestDB(t)
	customerStore := setupCustomerStore(t, db)
	t.Run("Returns empty array with no customers", func(t *testing.T) {
		customers, err := customerStore.GetCustomers(app.CustomerFilterOptions{})
		if err != nil {
			t.Fatal(err)
		}

		if len(customers.Customers) != 0 {
			t.Fatal("Incorrect amount of customers", len(customers.Customers))
		}
	})

	t.Run("correctly returns all customers", func(t *testing.T) {
		testCases := []struct {
			domain string
			id     string
		}{
			{"www.1.com", "1"},
			{"www.2.com", "2"},
			{"www.3.com", "3"},
		}

		for _, tc := range testCases {
			_, err := customerStore.GetCustomerID(tc.domain, tc.id)
			if err != nil {
				t.Fatalf("Failed to get customer ID for domain %s and id %s: %v", tc.domain, tc.id, err)
			}
		}

		customers, err := customerStore.GetCustomers(app.CustomerFilterOptions{
			PerPage: 10,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(customers.Customers) != 3 {
			t.Fatal("Incorrect amount of customers", customers)
		}
	})
}

func TestStoreCustomerData(t *testing.T) {
	db := setupTestDB(t)
	customerStore := setupCustomerStore(t, db)

	t.Run("store customer data", func(t *testing.T) {
		customer := app.Customer{
			SiteURL:    "www.test.com",
			LicensedTo: "test",
			Name:       "test",
		}

		customerID, err := customerStore.GetCustomerID(customer.SiteURL, customer.LicensedTo)

		if err != nil {
			t.Fatal(err)
		}

		customer.ID = customerID

		packet := &model.SupportPacket{
			LicenseTo:     "test",
			ServerVersion: "5.0.0",
			TotalPosts:    100,
		}
		config := &model.Config{
			ServiceSettings: model.ServiceSettings{
				SiteURL: &customer.SiteURL,
			},
		}
		plugins := &model.PluginsResponse{
			Active: []*model.PluginInfo{
				{
					Manifest: model.Manifest{
						Id:      "1",
						Name:    "1name",
						Version: "1.0.0",
					},
				},
				{
					Manifest: model.Manifest{
						Id:      "2",
						Name:    "2name",
						Version: "1.0.0",
					},
				},
			},
			Inactive: []*model.PluginInfo{
				{
					Manifest: model.Manifest{
						Id:      "3",
						Name:    "3name",
						Version: "1.0.0",
					},
				},
				{
					Manifest: model.Manifest{
						Id:      "4",
						Name:    "4name",
						Version: "1.0.0",
					},
				},
			},
		}

		pluginsResponse := []app.CustomerPluginValues{
			{
				PluginID: plugins.Active[0].Id,
				Name:     plugins.Active[0].Name,
				Version:  plugins.Active[0].Version,
				IsActive: true,
			},
			{
				PluginID: plugins.Active[1].Id,
				Name:     plugins.Active[1].Name,
				Version:  plugins.Active[1].Version,
				IsActive: true,
			},
			{
				PluginID: plugins.Inactive[0].Id,
				Name:     plugins.Inactive[0].Name,
				Version:  plugins.Inactive[0].Version,
				IsActive: false,
			},
			{
				PluginID: plugins.Inactive[1].Id,
				Name:     plugins.Inactive[1].Name,
				Version:  plugins.Inactive[1].Version,
				IsActive: false,
			},
		}

		packetResponse := app.CustomerPacketValues{
			LicensedTo: packet.LicenseTo,
			Version:    packet.ServerVersion,
			TotalPosts: packet.TotalPosts,
		}

		err = customerStore.UpdateCustomerThroughUpload(customerID, packet, config, plugins)
		if err != nil {
			t.Fatal(err)
		}

		customerInfo, err := customerStore.GetCustomerByID(customerID)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, *config.ServiceSettings.SiteURL, *customerInfo.Config.ServiceSettings.SiteURL, "config")

		if customerInfo.LastUpdated == 0 {
			t.Fatal("last updated not set in customer data")
		}
		customerInfo.LastUpdated = 0
		assertEqual(t, customer, customerInfo.Customer, "customer info")
		assertEqual(t, packetResponse, customerInfo.PacketValues, "packet data")
		assertEqual(t, pluginsResponse, customerInfo.Plugins, "plugin data")
	})
}

func TestUpdateCustomer(t *testing.T) {
	db := setupTestDB(t)
	customerStore := setupCustomerStore(t, db)
	t.Run("update customer", func(t *testing.T) {
		customer := app.Customer{
			SiteURL:                 "www.test.com",
			LicensedTo:              "test",
			Name:                    "test",
			AccountExecutive:        "bob",
			CustomerSuccessManager:  "jane",
			TechnicalAccountManager: "joe",
			SalesforceID:            "123",
			ZendeskID:               "456",
			Type:                    "cloud",
		}

		// creates the customer
		customerID, err := customerStore.GetCustomerID(customer.SiteURL, customer.LicensedTo)

		if err != nil {
			t.Fatal(err)
		}

		customer.ID = customerID

		err = customerStore.UpdateCustomer(customer)
		if err != nil {
			t.Fatal(err)
		}

		customerInfo, err := customerStore.GetCustomerByID(customerID)
		if err != nil {
			t.Fatal(err)
		}

		if customerInfo.LastUpdated == 0 {
			t.Fatal("last updated not set in customer data")
		}
		customerInfo.LastUpdated = 0
		assertEqual(t, customer, customerInfo.Customer, "customer info")
	})
}

func assertEqual(t *testing.T, expected, actual interface{}, name string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Incorrect %s. Expected %v got %v", name, expected, actual)
	}
}
