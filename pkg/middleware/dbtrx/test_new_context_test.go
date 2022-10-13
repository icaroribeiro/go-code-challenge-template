package dbtrx_test

import (
	"context"
	"testing"

	dbtrxmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/dbtrx"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestNewContext() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrxCtxValue := &gorm.DB{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingACopyOfAContextWithAnAssociatedValue",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedCtx := dbtrxmiddlewarepkg.NewContext(ctx, dbTrxCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedDBTrxCtxValue, ok := dbtrxmiddlewarepkg.FromContext(returnedCtx)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.Equal(t, dbTrxCtxValue, returnedDBTrxCtxValue)
			}
		})
	}
}
