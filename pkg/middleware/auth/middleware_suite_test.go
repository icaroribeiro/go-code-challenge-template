package auth_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Case struct {
	Context   string
	SetUp     func(t *testing.T)
	WantError bool
	TearDown  func(t *testing.T)
}

type Cases []Case

type TestSuite struct {
	suite.Suite
	SQLMock sqlmock.Sqlmock
	DB      *gorm.DB
	Cases   Cases
}

func (ts *TestSuite) SetupSuite() {
	var sqlDB *sql.DB
	var err error

	sqlDB, ts.SQLMock, err = sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		log.Panicf("failed to create a sqlmock database connection and a mock to manage expectations: %s", err.Error())
	}

	if sqlDB == nil {
		log.Panicf("The sqlDB is null")
	}

	if ts.SQLMock == nil {
		log.Panicf("The mock is null")
	}

	ts.SQLMock.ExpectPing()

	ts.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to open gorm db: %s", err.Error())
	}

	if ts.DB == nil {
		log.Panicf("The database is null")
	}

	if err = ts.DB.Error; err != nil {
		log.Panicf("%s", err.Error())
	}
}
