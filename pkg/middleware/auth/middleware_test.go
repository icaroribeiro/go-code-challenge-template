package auth_test

// import (
// 	"testing"

// 	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
// 	mockauthpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/mockauth"
// 	"github.com/stretchr/testify/suite"
// )

// func TestMiddlewareUnit(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestAuth() {
// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInWrappingAFunctionAndApplyingAuthenticationToARequest",
// 			SetUp:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authN := new(mockauthpkg.MockAuth)
// 			authMiddleware := authmiddlewarepkg.Auth(ts.DB)

// 		})
// 	}
// }
