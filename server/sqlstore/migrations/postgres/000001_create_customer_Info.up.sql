CREATE TABLE IF NOT EXISTS crm_customers (
	ID TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	customerSuccessManager TEXT DEFAULT '',
	accountExecutive TEXT DEFAULT '',
	technicalAccountManager TEXT DEFAULT '',
	salesforceId TEXT DEFAULT '',
	zendeskId TEXT DEFAULT '',
	licensedTo TEXT DEFAULT '',
	siteUrl TEXT DEFAULT '',
	type TEXT NOT NULL
);
