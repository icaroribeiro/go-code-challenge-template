package datastore_test

import (
	"fmt"
	"testing"

	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDatastore(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestNew() {
	dbConfig := map[string]string{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInInitializingThePostgresDriver",
			SetUp: func(t *testing.T) {
				dbConfig["DB_DRIVER"] = "postgres"
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheSQLDatabaseDriverIsNotRecognized",
			SetUp: func(t *testing.T) {
				dbConfig["DB_DRIVER"] = "testing"
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			i, err := datastorepkg.New(dbConfig)

			if !tc.WantError {
				_, ok := i.(datastorepkg.IDatastore)
				assert.True(t, ok, "Unexpected type assertion fail")
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}

func (ts *TestSuite) TestClose() {
	db, mock := NewMock()
	connPool := db.ConnPool

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInClosingTheDatabase",
			SetUp: func(t *testing.T) {
				mock.ExpectClose()
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingTheSQLDatabase",
			SetUp: func(t *testing.T) {
				db.ConnPool = nil
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				db.ConnPool = connPool
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			provider := datastorepkg.Provider{DB: db}
			err := provider.Close()

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))

			tc.TearDown(t)
		})
	}
}
