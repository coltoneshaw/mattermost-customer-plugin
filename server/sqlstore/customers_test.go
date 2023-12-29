package sqlstore

import (
	"testing"

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

func TestGetCustomerId(t *testing.T) {
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
		_, err := db.Exec(`INSERT INTO crm_customers (id, name, type) VALUES ($1, $2, $3)`, "1", "test", "cloud")
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
		_, err := db.Exec(`INSERT INTO crm_customers (id, name, type, siteurl, licensedto) VALUES ($1, $2, $3, $4, $5)`, "1", "test", "cloud", "www.1.com", "1")
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
		_, err := db.Exec(`INSERT INTO crm_customers (id, name, type, siteurl, licensedto) VALUES ($1, $2, $3, $4, $5)`, "2", "test", "cloud", "www.1.com", "1")
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
