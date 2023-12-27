CREATE TABLE IF NOT EXISTS crm_configValues (
	updateDataId TEXT NOT NULL,
	current BOOLEAN NOT NULL,
	config jsonb NOT NULL,
	customerId TEXT NOT NULL
);
