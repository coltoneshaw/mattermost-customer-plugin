CREATE TABLE IF NOT EXISTS crm_customers (
	ID TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	customerSuccessManager TEXT NOT NULL,
	accountExecutive TEXT NOT NULL,
	technicalAccountManager TEXT NOT NULL,
	salesforceId TEXT,
	zendeskId TEXT,
	licensedTo TEXT,
	siteUrl TEXT,
	type TEXT NOT NULL
);
