package user_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Case struct {
	Context    string
	SetUp      func(t *testing.T)
	StatusCode int
	WantError  bool
	TearDown   func(t *testing.T)
}

type Cases []Case

type ReturnArgs [][]interface{}

type TestSuite struct {
	suite.Suite
	Cases Cases
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
