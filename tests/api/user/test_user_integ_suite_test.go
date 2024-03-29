package user_test

import (
	"log"
	"testing"

	datastorepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/datastore"
	envpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/env"
	validatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator"
	uuidvalidatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator/uuid"
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
	dbDriver   = envpkg.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
	dbUser     = envpkg.GetEnvWithDefaultValue("DB_USER", "postgres")
	dbPassword = envpkg.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
	dbHost     = envpkg.GetEnvWithDefaultValue("DB_HOST", "localhost")
	dbPort     = envpkg.GetEnvWithDefaultValue("DB_PORT", "5434")
	dbName     = envpkg.GetEnvWithDefaultValue("DB_NAME", "testdb")
)

func setupDBConfig() (map[string]string, error) {
	dbConfig := map[string]string{
		"DRIVER":   dbDriver,
		"USER":     dbUser,
		"PASSWORD": dbPassword,
		"HOST":     dbHost,
		"PORT":     dbPort,
		"NAME":     dbName,
	}

	return dbConfig, nil
}

func (ts *TestSuite) SetupSuite() {
	dbConfig, err := setupDBConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	datastore, err := datastorepkg.New(dbConfig)
	if err != nil {
		log.Panic(err.Error())
	}

	ts.DB = datastore.GetInstance()
	if ts.DB == nil {
		log.Panicf("The database instance is null")
	}

	if err = ts.DB.Error; err != nil {
		log.Panicf("Got error when acessing the database instance: %s", err.Error())
	}

	validationFuncs := map[string]validatorv2.ValidationFunc{
		"uuid": uuidvalidatorpkg.Validate,
	}

	ts.Validator, err = validatorpkg.New(validationFuncs)
	if err != nil {
		log.Panicf("Got error when setting up the validator: %s", err.Error())
	}
}

func (ts *TestSuite) TearDownSuite() {
	db, err := ts.DB.DB()
	if err != nil {
		log.Panicf("Got error when acessing *sql.DB from database instance: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		log.Panicf("Got error when closing the database instance: %s", err.Error())
	}
}

func TestUserIntegSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
