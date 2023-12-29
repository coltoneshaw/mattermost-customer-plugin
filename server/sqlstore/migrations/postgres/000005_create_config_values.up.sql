CREATE TABLE IF NOT EXISTS crm_configValues (
	id TEXT NOT NULL PRIMARY KEY,
	updateDataId TEXT NOT NULL,
	current BOOLEAN NOT NULL,
	config jsonb NOT NULL,
	customerId TEXT NOT NULL
);
