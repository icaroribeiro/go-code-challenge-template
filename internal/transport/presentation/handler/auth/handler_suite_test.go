package auth_test

// import (
// 	"database/sql"
// 	"log"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// type Case struct {
// 	Context    string
// 	SetUp      func(t *testing.T)
// 	StatusCode int
// 	WantError  bool
// }

// type Cases []Case

// type ReturnArgs [][]interface{}

// type TestSuite struct {
// 	suite.Suite
// 	DB    *gorm.DB
// 	Cases Cases
// }

// func (ts *TestSuite) SetupSuite() {
// 	var sqlDB *sql.DB
// 	var err error

// 	sqlDB, _, err = sqlmock.New()
// 	if err != nil {
// 		log.Panicf("failed to create a sqlmock database connection and a mock to manage expectations: %s", err.Error())
// 	}

// 	if sqlDB == nil {
// 		log.Panicf("The sqlDB is null")
// 	}

// 	ts.DB, err = gorm.Open(postgres.New(postgres.Config{
// 		Conn: sqlDB,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		log.Panicf("failed to open gorm db: %s", err.Error())
// 	}

// 	if ts.DB == nil {
// 		log.Panicf("The database is null")
// 	}

// 	if err = ts.DB.Error; err != nil {
// 		log.Panicf("%s", err.Error())
// 	}
// }
