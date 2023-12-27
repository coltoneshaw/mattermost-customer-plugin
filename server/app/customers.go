package app

import "github.com/mattermost/mattermost/server/public/model"

type Customer struct {
	// ID is the unique identifier of the customer.
	ID string `json:"id"`

	// Name is the customers name and used as a somewhat unique identifier
	Name string `json:"name"`

	// The named CSM of the account.
	CustomerSuccessManager string `json:"customerSuccessManager"`

	// The named AE of the account.
	AccountExecutive string `json:"accountExecutive"`

	// The named TAM of the account.
	TechnicalAccountManager string `json:"technicalAccountManager"`

	// The Salesforce ID of the customer, manually added to the customer
	SalesforceID string `json:"salesforceId"`

	// The Zendesk Org ID of the customer, manually added to the customer
	ZendeskID string `json:"zendeskId"`

	// This field may be removed eventually, but this is just a way to try and do a check on
	// who this customer belongs to when a support packet comes in
	LicensedTo string `json:"licensed_to"`
	SiteURL    string `json:"siteURL"`

	// LicenseType is the type of license a customer can have
	// It can be "cloud", "enterprise", "professional", "free"
	Type string `json:"type"`

	PacketValues CustomerPacketValues `json:"packet"`

	Plugins []CustomerPluginValues `json:"plugins"`

	Config model.Config `json:"config"`
}

type CustomerPacketValues struct {
	LicensedTo            string `json:"licensedTo"`
	Version               string `json:"version"`
	ServerOS              string `json:"serverOS"`
	ServerArch            string `json:"serverArch"`
	DatabaseType          string `json:"databaseType"` // `mysql` or `postgres`
	DatabaseVersion       string `json:"databaseVersion"`
	DatabaseSchemaVersion string `json:"databaseSchemaVersion"`
	FileDriver            string `json:"fileDriver"`
	ActiveUsers           int64  `json:"activeUsers"`
	DailyActiveUsers      int64  `json:"dailyActiveUsers"`
	MonthlyActiveUsers    int64  `json:"monthlyActiveUsers"`
	InactiveUserCount     int64  `json:"inactiveUserCount"`
	LicenseSupportedUsers int64  `json:"licenseSupportedUsers"`
	TotalPosts            int64  `json:"totalPosts"`
	TotalChannels         int64  `json:"totalChannels"`
	TotalTeams            int64  `json:"totalTeams"`
}

type CustomerPluginValues struct {
	PluginID string `json:"pluginID"`
	Version  string `json:"version"`
	IsActive bool   `json:"isActive"`
	Name     string `json:"name"`
}

type CustomerService interface {
	// Get retrieves a customer based on id
	Get(id string) (Customer, error)
}

type CustomerStore interface {
	// GetCustomers returns filtered customers and the total count before paging.
	Get(id string) (Customer, error)
}
