package dbtrx_test

import (
	"context"
	"testing"

	dbtrxmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/dbtrx"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestFromContext() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrxCtxValue := &gorm.DB{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAnAssociatedValueFromAContext",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db
				ctx = dbtrxmiddlewarepkg.NewContext(ctx, dbTrxCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedDBTrxCtxValue, ok := dbtrxmiddlewarepkg.FromContext(ctx)

			if !tc.WantError {
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.NotEmpty(t, returnedDBTrxCtxValue)
				assert.Equal(t, dbTrxCtxValue, returnedDBTrxCtxValue)
			}
		})
	}
}
