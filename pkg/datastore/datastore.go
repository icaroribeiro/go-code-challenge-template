package datastore

import (
	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
	postgrespkg "github.com/icaroribeiro/go-code-challenge-template/pkg/datastore/driver/postgres"
	envpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/env"
)

type Provider struct{}

var (
	dbDriver = envpkg.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
)

// New is the factory function that encapsulates the implementation related to datastore.
func New(dbConfig map[string]string) (IDatastore, error) {
	switch dbDriver {
	case "postgres":
		return postgrespkg.New(dbConfig)
	}

	return nil, customerror.Newf("database driver %s is not recognized", dbDriver)
}
