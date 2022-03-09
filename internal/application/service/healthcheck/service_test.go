package healthcheck_test

// import (
// 	"errors"
// 	"fmt"
// 	"testing"

// 	healthcheckservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/healthcheck"
// 	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// func TestService(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestGetStatus() {
// 	connPool := ts.DB.ConnPool

// 	errorType := customerror.NoType

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingTheStatus",
// 			SetUp: func(t *testing.T) {
// 				ts.SQLMock.ExpectPing()
// 			},
// 			WantError: false,
// 			TearDown:  func() {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDBFunctionEvaluatesToAnError",
// 			SetUp: func(t *testing.T) {
// 				ts.DB.ConnPool = nil

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown: func(t *testing.T) {
// 				ts.DB.ConnPool = connPool
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfThePingCommandEvaluatesToAnError",
// 			SetUp: func(t *testing.T) {
// 				ts.SQLMock.ExpectPing().
// 					WillReturnError(errors.New("failed"))

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			healthCheckService := healthcheckservice.New(ts.DB)

// 			err := healthCheckService.GetStatus()

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) AfterTest(_, _ string) {
// 	err := ts.SQLMock.ExpectationsWereMet()
// 	assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
// }
