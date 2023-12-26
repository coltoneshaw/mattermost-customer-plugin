CREATE TABLE IF NOT EXISTS crm_updateData (
	id TEXT PRIMARY KEY,
	customerId TEXT NOT NULL,
	updatedBy TEXT NOT NULL,
	updatedAt INT NOT NULL,
	updateType TEXT NOT NULL,
	path TEXT NOT NULL
);
