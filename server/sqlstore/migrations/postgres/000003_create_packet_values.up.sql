CREATE TABLE IF NOT EXISTS crm_packetValues (
	updateDataId TEXT NOT NULL,
	licensedTo TEXT NOT NULL,
	version TEXT NOT NULL,
	serverOS TEXT NOT NULL,
	serverArch TEXT NOT NULL,
	databaseType TEXT NOT NULL,
	databaseVersion TEXT NOT NULL,
	databaseSchemaVersion TEXT NOT NULL,
	fileDriver TEXT NOT NULL,
	activeUsers TEXT NOT NULL,
	dailyActiveUsers INT NOT NULL,
	monthlyActiveUsers TEXT NOT NULL,
	inactiveUserCount INT NOT NULL,
	licenseSupportedUsers TEXT NOT NULL,
	totalPosts INT NOT NULL,
	totalChannels INT NOT NULL,
	totalTeams INT NOT NULL
);
