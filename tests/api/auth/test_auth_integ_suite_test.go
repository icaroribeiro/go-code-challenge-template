package auth_test

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt"
	authpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/auth"
	datastorepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/datastore"
	envpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/env"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template/pkg/security"
	validatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator"
	passwordvalidatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator/password"
	usernamevalidatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator/username"
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
	DB                *gorm.DB
	RSAKeys           authpkg.RSAKeys
	TokenExpTimeInSec int
	Validator         validatorpkg.IValidator
	Security          securitypkg.ISecurity
	Cases             Cases
}

var (
	publicKeyPath  = envpkg.GetEnvWithDefaultValue("RSA_PUBLIC_KEY_INTEG_PATH", "./../../configs/auth/rsa_keys/rsa.public")
	privateKeyPath = envpkg.GetEnvWithDefaultValue("RSA_PRIVATE_KEY_INTEG_PATH", "./../../configs/auth/rsa_keys/rsa.private")

	tokenExpTimeInSecStr = envpkg.GetEnvWithDefaultValue("TOKEN_EXP_TIME_IN_SEC", "120")

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
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	ts.RSAKeys = authpkg.RSAKeys{
		PublicKey:  rsaPublicKey,
		PrivateKey: rsaPrivateKey,
	}

	ts.TokenExpTimeInSec, err = strconv.Atoi(tokenExpTimeInSecStr)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

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
		"uuid":     uuidvalidatorpkg.Validate,
		"username": usernamevalidatorpkg.Validate,
		"password": passwordvalidatorpkg.Validate,
	}

	ts.Validator, err = validatorpkg.New(validationFuncs)
	if err != nil {
		log.Panicf("Got error when setting up the validator: %s", err.Error())
	}

	ts.Security = securitypkg.New()
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

func TestAuthIntegSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
