package healthcheck_test

import (
	"fmt"
	"testing"

	healthcheckservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/healthcheck"
	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestGetStatus() {
	driver := "postgres"
	db, mock := NewMockDB(driver)
	connPool := db.ConnPool

	errorType := customerror.NoType

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingTheStatus",
			SetUp: func(t *testing.T) {
				mock.ExpectPing()
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDBFunctionEvaluatesToAnError",
			SetUp: func(t *testing.T) {
				db.ConnPool = nil

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				db.ConnPool = connPool
			},
		},
		{
			Context: "ItShouldFailIfThePingCommandEvaluatesToAnError",
			SetUp: func(t *testing.T) {
				mock.ExpectPing().
					WillReturnError(customerror.New("failed"))

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckService := healthcheckservice.New(db)

			err := healthCheckService.GetStatus()

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))

			tc.TearDown(t)
		})
	}
}
