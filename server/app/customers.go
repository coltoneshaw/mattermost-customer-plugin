package app

// NOTE: When adding a column to the db, search for "When adding a PlaybookRun column" to see where
// that column needs to be added in the sqlstore code.
type Customer struct {
	// ID is the unique identifier of the customer.
	ID string `json:"id"`

	// Name is the customers name
	Name string `json:"name"`

	// LicenseType is the type of license a customer can have
	// It can be "cloud", "enterprise", "professional", "free"
	Type string `json:"type"`
}

type CustomerService interface {
	// Get retrieves a customer based on id
	Get(id string) (Customer, error)
}

type CustomerStore interface {
	// GetCustomers returns filtered customers and the total count before paging.
	Get(id string) (Customer, error)
}
