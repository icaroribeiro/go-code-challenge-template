package datastore

import (
	"fmt"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDriver struct {
	Provider Provider
}

// NewPostgresDriver is the factory function that encapsulates the implementation related to postgres.
func NewPostgresDriver(dbConfig map[string]string) (IDatastore, error) {
	dsn := ""

	if dbConfig["url"] != "" {
		dsn = dbConfig["url"]
	} else {
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			dbConfig["user"],
			dbConfig["password"],
			dbConfig["host"],
			dbConfig["port"],
			dbConfig["name"],
		)
	}

	dialector := postgres.Open(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return &PostgresDriver{}, customerror.Newf("failed to establish a database connection: %s", err.Error())
	}

	return &PostgresDriver{
		Provider{
			DB: db,
		},
	}, nil
}

// Close is the function that closes the database connection, releasing any open resources.
func (d *PostgresDriver) Close() error {
	return d.Provider.Close()
}
