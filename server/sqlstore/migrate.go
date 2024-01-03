package sqlstore

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/morph"
	"github.com/mattermost/morph/drivers"
	ps "github.com/mattermost/morph/drivers/postgres"
	"github.com/mattermost/morph/sources"
	"github.com/mattermost/morph/sources/embedded"
	"github.com/sirupsen/logrus"
)

//go:embed migrations
var assets embed.FS

const driverName = "postgres"

// RunMigrations will run the migrations (if any). The caller should hold a cluster mutex if there
// is a danger of this being run on multiple servers at once.
func (sqlStore *SQLStore) RunMigrations() error {
	if err := sqlStore.runMigrationsWithMorph(); err != nil {
		return fmt.Errorf("failed to complete migrations (with morph): %w", err)
	}

	return nil
}

func (sqlStore *SQLStore) createSource() (sources.Source, error) {
	driverName := sqlStore.db.DriverName()
	assetsList, err := assets.ReadDir(filepath.Join("migrations", driverName))
	if err != nil {
		return nil, err
	}

	assetNamesForDriver := make([]string, len(assetsList))
	for i, entry := range assetsList {
		assetNamesForDriver[i] = entry.Name()
	}

	src, err := embedded.WithInstance(&embedded.AssetSource{
		Names: assetNamesForDriver,
		AssetFunc: func(name string) ([]byte, error) {
			return assets.ReadFile(filepath.Join("migrations", driverName, name))
		},
	})

	return src, err
}

func (sqlStore *SQLStore) createDriver() (drivers.Driver, error) {
	driverName := sqlStore.db.DriverName()
	if driverName != model.DatabaseDriverPostgres {
		return nil, fmt.Errorf("unsupported database type %s for migration", driverName)
	}
	return ps.WithInstance(sqlStore.db.DB)
}

func (sqlStore *SQLStore) createMorphEngine() (*morph.Morph, error) {
	src, err := sqlStore.createSource()
	if err != nil {
		return nil, err
	}

	driver, err := sqlStore.createDriver()
	if err != nil {
		return nil, err
	}

	opts := []morph.EngineOption{
		morph.WithLock("mm-customers-lock-key"),
		morph.SetMigrationTableName("CRM_db_migrations"),
		morph.SetStatementTimeoutInSeconds(30),
	}
	engine, err := morph.New(context.Background(), driver, src, opts...)

	return engine, err
}

// WARNING: We don't use morph migration until proper testing
func (sqlStore *SQLStore) runMigrationsWithMorph() error {
	engine, err := sqlStore.createMorphEngine()
	if err != nil {
		return err
	}
	defer func() {
		logrus.Debug("Closing morph engine")
		err = engine.Close()
		if err != nil {
			logrus.Errorf("Error closing morph engine. err: '%v'", err)
			return
		}
		logrus.Debug("Closed morph engine")
	}()
	if err := engine.ApplyAll(); err != nil {
		return fmt.Errorf("could not apply migrations: %w", err)
	}
	return nil
}
