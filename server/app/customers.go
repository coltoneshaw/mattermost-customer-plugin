package app

import "github.com/mattermost/mattermost/server/public/model"

type LicenseType string

const (
	Cloud        LicenseType = "cloud"
	Enterprise   LicenseType = "enterprise"
	Professional LicenseType = "professional"
	Free         LicenseType = "free"
	Trial        LicenseType = "trial"
	NonProfit    LicenseType = "nonprofit"
	Other        LicenseType = "other"
)

type Customer struct {
	ID                      string      `json:"id"`
	Name                    string      `json:"name"`
	LastUpdated             int64       `json:"lastUpdated"`
	SalesforceID            string      `json:"salesforceId"`
	ZendeskID               string      `json:"zendeskId"`
	CustomerSuccessManager  string      `json:"customerSuccessManager"`
	AccountExecutive        string      `json:"accountExecutive"`
	TechnicalAccountManager string      `json:"technicalAccountManager"`
	ProductManager          string      `json:"productManager"`
	LicensedTo              string      `json:"licensedTo"`
	SiteURL                 string      `json:"siteURL"`
	LicenseType             LicenseType `json:"licenseType"`
	CustomerChannel         string      `json:"customerChannel"`
	GDriveLink              string      `json:"GDriveLink"`
	AirGapped               bool        `json:"airGapped"`
	AirGappedReason         string      `json:"airGappedReason"`
	Region                  string      `json:"region"`      // emea, apac, amer
	Status                  string      `json:"status"`      // gold standard, onboarding, stable
	CompanyType             string      `json:"companyType"` // enterprise, federal, midmarket, smb,
	CodeWord                string      `json:"codeWord"`
}

// todo - modify the licnesedTo to match mattermost with licenseto
type CustomerPacketValues struct {
	// AuditID               string `json:"auditID"`
	LicensedTo            string `json:"licensedTo"`
	Version               string `json:"version"`
	ServerOS              string `json:"serverOS"`
	ServerArch            string `json:"serverArch"`
	DatabaseType          string `json:"databaseType"` // `mysql` or `postgres`
	DatabaseVersion       string `json:"databaseVersion"`
	DatabaseSchemaVersion string `json:"databaseSchemaVersion"`
	FileDriver            string `json:"fileDriver"`
	ActiveUsers           int    `json:"activeUsers"`
	DailyActiveUsers      int    `json:"dailyActiveUsers"`
	MonthlyActiveUsers    int    `json:"monthlyActiveUsers"`
	InactiveUserCount     int    `json:"inactiveUserCount"`
	LicenseSupportedUsers int    `json:"licenseSupportedUsers"`
	TotalPosts            int    `json:"totalPosts"`
	TotalChannels         int    `json:"totalChannels"`
	TotalTeams            int    `json:"totalTeams"`
	ElasticServerVersion  string `json:"elasticServerVersion"`
	Metrics               bool   `json:"metrics"` // This will be a true / false or string of what they're using.
	MetricService         string `json:"metricService"`
	HostingType           string `json:"hostingType"`    // aws, onprem, azure, cloud,
	DeploymentType        string `json:"deploymentType"` // kube, tar, docker, AKS, EKS, GKE, etc
	MobileApp             bool   `json:"mobileApp"`      // true, false, NA
	ProductsInUse         string `json:"productsInUse"`  // calls, playbooks
	SAMLProvider          string `json:"samlProvider"`   // Name of their SAML technology
	LDAPProvider          string `json:"ldapProvider"`   // Name of their LDAP technology
}

type CustomerPluginValues struct {
	PluginID    string `json:"pluginID"`
	Version     string `json:"version"`
	IsActive    bool   `json:"isActive"`
	Name        string `json:"name"`
	HomePageURL string `json:"homePageURL"`
}

type FullCustomerInfo struct {
	Customer
	PacketValues CustomerPacketValues   `json:"packet"`
	Plugins      []CustomerPluginValues `json:"plugins"`
	Config       model.Config           `json:"config"`
}

type CustomerService interface {

	// GetCustomers returns filtered customers
	GetCustomers(opts CustomerFilterOptions) (GetCustomersResult, error)

	// Get retrieves a customer based on id
	GetCustomerByID(id string) (FullCustomerInfo, error)

	// Checks to see if a customer exists based on the siteURL and licensedTo
	GetCustomerID(siteURL string, licensedTo string) (id string, err error)

	// This monitors the posts for a support packet and triggers actions based on that
	MessageHasBeenPosted(post *model.Post)

	GetPacket(customerID string) (CustomerPacketValues, error)
	// StorePacket(updateId string, packet CustomerPacketValues) error

	GetConfig(customerID string) (model.Config, error)
	GetPlugins(customerID string) ([]CustomerPluginValues, error)

	UpdateCustomer(customer Customer) error
	UpdateCustomerData(customerID string, userID string, packet *CustomerPacketValues, config *model.Config, plugins []CustomerPluginValues) error
}

type CustomerStore interface {
	GetCustomers(opts CustomerFilterOptions) (GetCustomersResult, error)

	// GetCustomers returns filtered customers and the total count before paging.
	GetCustomerByID(id string) (FullCustomerInfo, error)

	// Checks to see if a customer exists based on the siteURL and licensedTo
	GetCustomerID(siteURL string, licensedTo string) (id string, err error)

	GetPacket(customerID string) (CustomerPacketValues, error)
	// StorePacket(updateId string, packet CustomerPacketValues) error

	GetConfig(customerID string) (model.Config, error)
	GetPlugins(customerID string) ([]CustomerPluginValues, error)

	UpdateCustomer(customer Customer) error
	UpdateCustomerData(customerID string, userID string, packet *CustomerPacketValues, config *model.Config, plugins []CustomerPluginValues) error

	UpdateCustomerThroughUpload(customerID string, packet *model.SupportPacket, config *model.Config, plugins *model.PluginsResponse) error
}

type GetCustomersResult struct {
	TotalCount int        `json:"totalCount"`
	PageCount  int        `json:"pageCount"`
	HasMore    bool       `json:"hasMore"`
	Customers  []Customer `json:"customers"`
}

type CustomerFilterOptions struct {
	Sort       SortField
	Direction  SortDirection
	SearchTerm string

	// Pagination options.
	Page    int
	PerPage int
}
