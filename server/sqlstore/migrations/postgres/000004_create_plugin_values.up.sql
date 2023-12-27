CREATE TABLE IF NOT EXISTS crm_pluginValues (
	updateDataId TEXT NOT NULL,
	customerId TEXT NOT NULL,
	pluginId TEXT NOT NULL,
	version TEXT NOT NULL,
	isActive BOOLEAN NOT NULL,
	name TEXT NOT NULL,
	current BOOLEAN NOT NULL
);
