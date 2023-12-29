CREATE TYPE customer_type AS ENUM ('cloud', 'enterprise', 'professional', 'trial', 'free', 'nonprofit', 'other', '');

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
	type customer_type NOT NULL,
	customerChannel TEXT DEFAULT '',
	gdriveLink TEXT DEFAULT ''
);
