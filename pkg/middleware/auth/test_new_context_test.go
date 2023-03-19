package auth_test

import (
	"context"
	"testing"

	domainentity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	domainentityfactory "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNewContext() {
	authDetailsCtxValue := domainentity.Auth{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingACopyOfAContextWithAnAssociatedValue",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainentityfactory.AuthFactory(nil)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedCtx := authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedAuthDetailsCtxValue, ok := authmiddlewarepkg.FromContext(returnedCtx)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.Equal(t, authDetailsCtxValue, returnedAuthDetailsCtxValue)
			}
		})
	}
}
