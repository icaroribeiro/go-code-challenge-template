package user_test

import (
	"log"
	"testing"

	"github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/env"
	databasepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/datastore"
	validatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator"
	"github.com/icaroribeiro/go-code-challenge-template/pkg/validator/validationfunc"
	"github.com/stretchr/testify/suite"
	validatorv2 "gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type Case struct {
	Context    string
	SetUp      func(t *testing.T)
	StatusCode int
	WantError  bool
	TearDown   func(t *testing.T)
}

type Cases []Case

type TestSuite struct {
	suite.Suite
	DB        *gorm.DB
	Validator validatorpkg.IValidator
	Cases     Cases
}

var (
	dbDriver   = env.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
	dbUser     = env.GetEnvWithDefaultValue("DB_USER", "postgres")
	dbPassword = env.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
	dbHost     = env.GetEnvWithDefaultValue("DB_HOST", "localhost")
	dbPort     = env.GetEnvWithDefaultValue("DB_PORT", "5432")
	dbName     = env.GetEnvWithDefaultValue("DB_NAME", "testdb")
)

func (ts *TestSuite) SetupSuite() {
	dbConfig := setDBConfig()

	database, err := databasepkg.New(dbConfig)
	if err != nil {
		log.Panic(err.Error())
	}

	ts.DB = database.Instance
	if ts.DB == nil {
		log.Panicf("The database is null")
	}

	if err = ts.DB.Error; err != nil {
		log.Panicf("%s", err.Error())
	}

	validationFunc := validationfunc.New()

	validationFuncMap := map[string]validatorv2.ValidationFunc{
		"uuid":     validationFunc.ValidateUUID,
		"name":     validationFunc.ValidateName,
		"username": validationFunc.ValidateUsername,
		"password": validationFunc.ValidatePassword,
	}

	ts.Validator, err = validatorpkg.New(validationFuncMap)
	if err != nil {
		log.Panicf("%s", err.Error())
	}
}

func (ts *TestSuite) TearDownSuite() {
	db, err := ts.DB.DB()
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	if err = db.Close(); err != nil {
		log.Panicf("%s", err.Error())
	}
}

func setDBConfig() map[string]string {
	dbConfig := map[string]string{
		"driver":   dbDriver,
		"user":     dbUser,
		"password": dbPassword,
		"host":     dbHost,
		"port":     dbPort,
		"name":     dbName,
	}

	return dbConfig
}
