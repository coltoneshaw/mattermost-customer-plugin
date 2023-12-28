package sqlstore

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"

	"github.com/isacikgoz/morph"
	"github.com/isacikgoz/morph/drivers"
	ps "github.com/isacikgoz/morph/drivers/postgres"
	"github.com/isacikgoz/morph/sources/embedded"
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

func (sqlStore *SQLStore) runMigrationsWithMorph() error {
	assetsList, err := assets.ReadDir(filepath.Join("migrations", driverName))
	if err != nil {
		return err
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
	if err != nil {
		return err
	}

	config := drivers.Config{
		StatementTimeoutInSecs: 100000,
		MigrationsTable:        "CRM_db_migrations",
	}

	var driver drivers.Driver

	driver, err = ps.WithInstance(sqlStore.db.DB, &ps.Config{
		Config: config,
	})

	if err != nil {
		return err
	}

	opts := []morph.EngineOption{
		morph.WithLock("mm-customers-lock-key"),
	}
	engine, err := morph.New(context.Background(), driver, src, opts...)
	if err != nil {
		return err
	}
	defer engine.Close()

	if err := engine.ApplyAll(); err != nil {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	return nil
}
