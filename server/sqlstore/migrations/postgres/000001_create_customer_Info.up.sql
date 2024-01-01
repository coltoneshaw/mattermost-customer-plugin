DO
$$
BEGIN
  IF NOT EXISTS (SELECT * FROM pg_type typ
                            INNER JOIN pg_namespace nsp ON nsp.oid = typ.typnamespace
                        WHERE nsp.nspname = current_schema()
                            AND typ.typname = 'customer_type') THEN
    CREATE TYPE customer_type AS ENUM ('cloud', 'enterprise', 'professional', 'trial', 'free', 'nonprofit', 'other', '');
  END IF;
END;
$$
LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS crm_customers (
	ID TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	lastUpdated bigint NOT NULL,
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
