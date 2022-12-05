// Code generated by mockery v2.10.0. DO NOT EDIT.

package mockauth

import (
	jwt "github.com/dgrijalva/jwt-go"
	entity "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// Auth is an autogenerated mock type for the IAuth type
type Auth struct {
	mock.Mock
}

// CreateToken provides a mock function with given fields: _a0, tokenExpTimeInSec
func (_m *Auth) CreateToken(_a0 entity.Auth, tokenExpTimeInSec int) (string, error) {
	ret := _m.Called(_a0, tokenExpTimeInSec)

	var r0 string
	if rf, ok := ret.Get(0).(func(entity.Auth, int) string); ok {
		r0 = rf(_a0, tokenExpTimeInSec)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(entity.Auth, int) error); ok {
		r1 = rf(_a0, tokenExpTimeInSec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecodeToken provides a mock function with given fields: tokenString
func (_m *Auth) DecodeToken(tokenString string) (*jwt.Token, error) {
	ret := _m.Called(tokenString)

	var r0 *jwt.Token
	if rf, ok := ret.Get(0).(func(string) *jwt.Token); ok {
		r0 = rf(tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExtractTokenString provides a mock function with given fields: authHeaderString
func (_m *Auth) ExtractTokenString(authHeaderString string) (string, error) {
	ret := _m.Called(authHeaderString)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(authHeaderString)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(authHeaderString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchAuthFromToken provides a mock function with given fields: token
func (_m *Auth) FetchAuthFromToken(token *jwt.Token) (entity.Auth, error) {
	ret := _m.Called(token)

	var r0 entity.Auth
	if rf, ok := ret.Get(0).(func(*jwt.Token) entity.Auth); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(entity.Auth)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*jwt.Token) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateTokenRenewal provides a mock function with given fields: token, timeBeforeTokenExpTimeInSec
func (_m *Auth) ValidateTokenRenewal(token *jwt.Token, timeBeforeTokenExpTimeInSec int) (*jwt.Token, error) {
	ret := _m.Called(token, timeBeforeTokenExpTimeInSec)

	var r0 *jwt.Token
	if rf, ok := ret.Get(0).(func(*jwt.Token, int) *jwt.Token); ok {
		r0 = rf(token, timeBeforeTokenExpTimeInSec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*jwt.Token, int) error); ok {
		r1 = rf(token, timeBeforeTokenExpTimeInSec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
