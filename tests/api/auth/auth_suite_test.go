package auth_test

// import (
// 	"io/ioutil"
// 	"log"
// 	"strconv"
// 	"testing"

// 	"github.com/dgrijalva/jwt-go"
// 	authinfra "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/auth"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/env"
// 	securitypkg "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/security"
// 	databasepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/datastore"
// 	validatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator"
// 	"github.com/icaroribeiro/go-code-challenge-template/pkg/validator/validationfunc"
// 	"github.com/stretchr/testify/suite"
// 	validatorv2 "gopkg.in/validator.v2"
// 	"gorm.io/gorm"
// )

// type Case struct {
// 	Context    string
// 	SetUp      func(t *testing.T)
// 	StatusCode int
// 	WantError  bool
// 	TearDown   func(t *testing.T)
// }

// type Cases []Case

// type TestSuite struct {
// 	suite.Suite
// 	DB                *gorm.DB
// 	RSAKeys           authinfra.RSAKeys
// 	TokenExpTimeInSec int
// 	Validator         validatorpkg.IValidator
// 	Security          securitypkg.ISecurity
// 	Cases             Cases
// }

// var (
// 	publicKeyPath  string
// 	privateKeyPath string

// 	tokenExpTimeInSecStr string

// 	dbDriver   string
// 	dbUser     string
// 	dbPassword string
// 	dbHost     string
// 	dbPort     string
// 	dbName     string
// )

// func init() {
// 	publicKeyPath = env.GetEnvWithDefaultValue("RSA_PUBLIC_KEY_PATH", "../../../tests/configs/auth/rsa_keys/rsa.public")
// 	privateKeyPath = env.GetEnvWithDefaultValue("RSA_PRIVATE_KEY_PATH", "../../../tests/configs/auth/rsa_keys/rsa.private")

// 	tokenExpTimeInSecStr = env.GetEnvWithDefaultValue("TOKEN_EXP_TIME_IN_SEC", "600")

// 	dbDriver = env.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
// 	dbUser = env.GetEnvWithDefaultValue("DB_USER", "postgres")
// 	dbPassword = env.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
// 	dbHost = env.GetEnvWithDefaultValue("DB_HOST", "localhost")
// 	dbPort = env.GetEnvWithDefaultValue("DB_PORT", "5432")
// 	dbName = env.GetEnvWithDefaultValue("DB_NAME", "testdb")
// }

// func (ts *TestSuite) SetupSuite() {
// 	publicKey, err := ioutil.ReadFile(publicKeyPath)
// 	if err != nil {
// 		log.Panicf("%s", err.Error())
// 	}

// 	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
// 	if err != nil {
// 		log.Panicf("%s", err.Error())
// 	}

// 	privateKey, err := ioutil.ReadFile(privateKeyPath)
// 	if err != nil {
// 		log.Panicf("%s", err.Error())
// 	}

// 	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
// 	if err != nil {
// 		log.Panicf("%s", err.Error())
// 	}

// 	ts.RSAKeys = authinfra.RSAKeys{
// 		PublicKey:  rsaPublicKey,
// 		PrivateKey: rsaPrivateKey,
// 	}

// 	ts.TokenExpTimeInSec, err = strconv.Atoi(tokenExpTimeInSecStr)
// 	if err != nil {
// 		log.Panicf("%s", err.Error())
// 	}

// 	dbConfig := setDBConfig()

// 	database, err := databasepkg.New(dbConfig)
// 	if err != nil {
// 		log.Panic(err.Error())
// 	}

// 	ts.DB = database.Instance
// 	if ts.DB == nil {
// 		log.Panicf("The database is null")
// 	}

// 	if err = ts.DB.Error; err != nil {
// 		log.Panicf("%s", err.Error())
// 	}

// 	validationFunc := validationfunc.New()

// 	validationFuncMap := map[string]validatorv2.ValidationFunc{
// 		"uuid":     validationFunc.ValidateUUID,
// 		"name":     validationFunc.ValidateName,
// 		"username": validationFunc.ValidateUsername,
// 		"password": validationFunc.ValidatePassword,
// 	}

// 	ts.Validator, err = validatorpkg.New(validationFuncMap)
// 	if err != nil {
// 		log.Panicf("%s", err.Error())
// 	}

// 	ts.Security = securitypkg.New()
// }

// func (ts *TestSuite) TearDownSuite() {
// 	db, err := ts.DB.DB()
// 	if err != nil {
// 		log.Panicf("%s", err.Error())
// 	}

// 	if err = db.Close(); err != nil {
// 		log.Panicf("%s", err.Error())
// 	}
// }

// func setDBConfig() map[string]string {
// 	dbConfig := map[string]string{
// 		"driver":   dbDriver,
// 		"user":     dbUser,
// 		"password": dbPassword,
// 		"host":     dbHost,
// 		"port":     dbPort,
// 		"name":     dbName,
// 	}

// 	return dbConfig
// }
