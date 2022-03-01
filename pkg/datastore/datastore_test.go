package datastore_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/suite"
// )

// func TestDatabase(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestNew() {
// 	ts.Cases = Cases{
// 		{
// 			Context:   "ItShouldSucceedInInitializingTheDatabase",
// 			SetUp:     func(t *testing.T) {},
// 			WantError: false,
// 		},
// 		{
// 			Context:   "ItShouldFailIfAnErrorIsReturnedWhenSettingTheDialector",
// 			SetUp:     func(t *testing.T) {},
// 			WantError: false,
// 		},
// 		{
// 			Context:   "ItShouldFailIfAnErrorIsReturnedWhenOpeningTheDatabaseSessionBasedOnTheDialector",
// 			SetUp:     func(t *testing.T) {},
// 			WantError: false,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			if !tc.WantError {
// 			} else {
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestGetInstance() {
// 	ts.Cases = Cases{
// 		{
// 			Context:   "ItShouldSucceedInGettingTheDatabaseInstance",
// 			SetUp:     func(t *testing.T) {},
// 			WantError: false,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			if !tc.WantError {
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestClose() {
// 	ts.Cases = Cases{
// 		{
// 			Context:   "ItShouldSucceedInClosingTheDatabase",
// 			SetUp:     func(t *testing.T) {},
// 			WantError: false,
// 		},
// 		{
// 			Context:   "ItShouldFailIfAnErrorIsReturnedWhenGettingTheSQLDatabase",
// 			SetUp:     func(t *testing.T) {},
// 			WantError: true,
// 		},
// 		{
// 			Context:   "ItShouldFailIfAnErrorIsReturnedWhenClosingTheSQLDatabase",
// 			SetUp:     func(t *testing.T) {},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			if !tc.WantError {
// 			} else {
// 			}
// 		})
// 	}
// }
