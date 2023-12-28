package sqlstore

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/store/storetest"

	sq "github.com/mattermost/squirrel"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t testing.TB) *sqlx.DB {
	t.Helper()

	sqlSettings := storetest.MakeSqlSettings("postgres", false)

	origDB, err := sql.Open(*sqlSettings.DriverName, *sqlSettings.DataSource)
	require.NoError(t, err)

	db := sqlx.NewDb(origDB, driverName)
	if driverName == model.DatabaseDriverMysql {
		db.MapperFunc(func(s string) string { return s })
	}

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
		storetest.CleanupSqlSettings(sqlSettings)
	})

	return db
}

func setupTables(t *testing.T, db *sqlx.DB) *SQLStore {
	t.Helper()

	driverName := db.DriverName()

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Question)
	if driverName == model.DatabaseDriverPostgres {
		builder = builder.PlaceholderFormat(sq.Dollar)
	}

	sqlStore := &SQLStore{
		db,
		builder,
	}

	return sqlStore
}

func setupSQLStore(t *testing.T, db *sqlx.DB) *SQLStore {
	sqlStore := setupTables(t, db)

	err := sqlStore.RunMigrations()
	require.NoError(t, err)

	return sqlStore
}
