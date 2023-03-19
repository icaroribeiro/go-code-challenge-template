package auth_test

import (
	"context"
	"testing"

	domainentity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	domainentityfactory "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestFromContext() {
	authDetailsCtxValue := domainentity.Auth{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAnAssociatedValueFromAContext",
			SetUp: func(t *testing.T) {
				authDetailsCtxValue = domainentityfactory.AuthFactory(nil)
				ctx = authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedAuthDetailsCtxValue, ok := authmiddlewarepkg.FromContext(ctx)

			if !tc.WantError {
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.NotEmpty(t, returnedAuthDetailsCtxValue)
				assert.Equal(t, authDetailsCtxValue, returnedAuthDetailsCtxValue)
			}
		})
	}
}
