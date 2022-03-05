package datastore_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestPostgresDriver(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// func (ts *TestSuite) TestPostgresDriverNew() {
// 	_, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "11",
// 		Env: []string{
// 			"POSTGRES_DB=test",
// 			"POSTGRES_USER=test",
// 			"POSTGRES_PASSWORD=test",
// 			"listen_addresses = '*'",
// 		},
// 	}, func(config *docker.HostConfig) {
// 		config.AutoRemove = true
// 		config.RestartPolicy = docker.RestartPolicy{
// 			Name: "no",
// 		}
// 	})

// 	fmt.Print(resource.GetPort("5432"))

// 	dbConfig := map[string]string{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedUsingAURL",
// 			SetUp: func(t *testing.T) {
// 				dbConfig["url"] = "postgres://test:test@postgres:5432/test?sslmode=disable"
// 			},
// 			WantError: true,
// 		},
// 		// {
// 		// 	Context: "ItShouldSucceedUsingAConnectionString",
// 		// 	SetUp: func(t *testing.T) {
// 		// 		dbConfig["url"] = "testing"
// 		// 	},
// 		// 	WantError: false,
// 		// },
// 		// {
// 		// 	Context: "ItShouldFailIfItIsNotPossibleToParseDSN",
// 		// 	SetUp: func(t *testing.T) {
// 		// 		dbConfig["url"] = fake.Word()
// 		// 	},
// 		// 	WantError: true,
// 		// },
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			i, err := datastorepkg.NewPostgresDriver(dbConfig)

// 			if !tc.WantError {
// 				_, ok := i.(datastorepkg.IDatastore)
// 				assert.True(t, ok, "Unexpected type assertion fail")
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost")
// 			}
// 		})
//	}
// }

func (ts *TestSuite) TestPostgresDriverClose() {
	db := &gorm.DB{}
	var mock sqlmock.Sqlmock
	provider := datastorepkg.Provider{}

	var connPool gorm.ConnPool

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInClosingTheDatabase",
			SetUp: func(t *testing.T) {
				db, mock = NewMock()
				provider = datastorepkg.Provider{DB: db}
				mock.ExpectClose()
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingTheSQLDatabase",
			SetUp: func(t *testing.T) {
				db, _ = NewMock()
				connPool = db.ConnPool
				db.ConnPool = nil
				provider = datastorepkg.Provider{DB: db}
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

			postgresDriver := datastorepkg.PostgresDriver{Provider: provider}
			err := postgresDriver.Close()

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
