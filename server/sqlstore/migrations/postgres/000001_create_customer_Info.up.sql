CREATE TABLE IF NOT EXISTS crm_customers (
	INSTEAD TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	customerSuccessManager TEXT NOT NULL,
	accountExecutive TEXTm NOT NULL,
	technicalAccountManager TEXT NOT NULL,
	salesforceId TEXT,
	zendeskId TEXT,
	licensedTo TEXT
	siteUrl TEXT,
	type TEXT NOT NULL,
);
