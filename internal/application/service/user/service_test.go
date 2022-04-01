package user_test

// import (
// 	"fmt"
// 	"testing"

// 	userservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/application/service/user"
// 	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
// 	userdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/infrastructure/storage/datastore/mockrepository/user"
// 	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
// 	domainfactorymodel "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/core/domain/model"
// 	"github.com/icaroribeiro/new-go-code-challenge-template/tests/mocks/pkg/mockvalidator"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// func TestService(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestGetAll() {
// 	user := domainmodel.User{}

// 	returnArgs := ReturnArgs{}

// 	errorType := customerror.NoType

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingAllUsers",
// 			SetUp: func(t *testing.T) {
// 				user = domainfactorymodel.NewUser(nil)

// 				returnArgs = ReturnArgs{
// 					{domainmodel.Users{user}, nil},
// 				}
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfItIsNotPossibleToGetAllUsers",
// 			SetUp: func(t *testing.T) {
// 				returnArgs = ReturnArgs{
// 					{domainmodel.Users{}, customerror.New("failed")},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			userDatastoreRepository := new(userdatastoremockrepository.Repository)
// 			userDatastoreRepository.On("GetAll").Return(returnArgs[0]...)

// 			validator := new(mockvalidator.Validator)

// 			userService := userservice.New(userDatastoreRepository, validator)

// 			returnedUsers, err := userService.GetAll()

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, user.ID, returnedUsers[0].ID)
// 				assert.Equal(t, user.Username, returnedUsers[0].Username)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedUsers)
// 			}
// 		})
// 	}
// }