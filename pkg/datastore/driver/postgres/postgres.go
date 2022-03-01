package postgres

import (
	"fmt"

	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Driver struct {
	Name string
	DB   *gorm.DB
}

// New is the factory function that encapsulates the implementation related to postgres.
func New(dbConfig map[string]string) (*Driver, error) {
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
		return &Driver{}, customerror.Newf("failed to establish a database connection: %s", err.Error())
	}

	return &Driver{
		Name: dbConfig["driver"],
		DB:   db,
	}, nil
}

// GetDB is the function that gets the "client" for peforming database operations.
func (d *Driver) GetDB() interface{} {
	return d.DB
}

// SetDB is the function that gets the "client" for peforming database operations.
func (d *Driver) SetDB(db interface{}) {
	d.DB = db.(*gorm.DB)
}

// CheckStatus is the function that checks the database status.
func (d *Driver) CheckStatus() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}

// Close is the function that closes the database connection, releasing any open resources.
func (d *Driver) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
