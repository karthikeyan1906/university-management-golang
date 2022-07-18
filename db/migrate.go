package migrations

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/pkg/errors"
)

func MigrationsUp(username, password, host, port, dbName, schema string) error {
	assetFunc := func(name string) ([]byte, error) {
		asset, err := Asset(name) // TODO: Generated code
		if err != nil {
			return nil, fmt.Errorf("error while creating driver: %v", err)
		}
		return asset, nil
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=%s&sslmode=%s",
		username, password, host, port,
		dbName, schema, "disable")
	err := runMigrations(AssetNames(), assetFunc, "university-management", dsn)
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println(migrate.ErrNoChange.Error())
		} else {
			return errors.Errorf("migration failed with error - %v", err)
		}
	}

	return nil
}

func runMigrations(assets []string, asset func(name string) ([]byte, error), appName string, dsn string) error {
	assetSource := bindata.Resource(assets, func(name string) ([]byte, error) {
		asset, err := asset(name)
		if err != nil {
			return nil, errors.Errorf("error while creating driver: %v", err)
		}
		return asset, nil
	})

	driver, err := bindata.WithInstance(assetSource)
	if err != nil {
		return errors.Errorf("error while creating driver: %v", err)
	}

	sourceName := fmt.Sprintf("%s-migrations", appName)

	migrateInstance, err := migrate.NewWithSourceInstance(sourceName, driver, dsn)
	if err != nil {
		return errors.Errorf("error while creating Migrate instance: %v", err)
	}

	return migrateInstance.Up()
}
