package auth_test

import (
	"log"
	"strconv"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/env"
	"github.com/stretchr/testify/suite"
)

type Case struct {
	Context   string
	SetUp     func(t *testing.T)
	WantError bool
	TearDown  func(t *testing.T)
}

type Cases []Case

type ReturnArgs [][]interface{}

type TestSuite struct {
	suite.Suite
	Cases             Cases
	TokenExpTimeInSec int
}

var (
	tokenExpTimeInSecStr = env.GetEnvWithDefaultValue("TOKEN_EXP_TIME_IN_SEC", "600")
)

func (ts *TestSuite) SetupSuite() {
	var err error

	ts.TokenExpTimeInSec, err = strconv.Atoi(tokenExpTimeInSecStr)
	if err != nil {
		log.Panicf("%s", err.Error())
	}
}
