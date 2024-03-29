// Code generated by mockery v2.16.0. DO NOT EDIT.

package user

import (
	entity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	user "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/application/service/user"
)

// Service is an autogenerated mock type for the IService type
type Service struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *Service) GetAll() (entity.Users, error) {
	ret := _m.Called()

	var r0 entity.Users
	if rf, ok := ret.Get(0).(func() entity.Users); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(entity.Users)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithDBTrx provides a mock function with given fields: dbTrx
func (_m *Service) WithDBTrx(dbTrx *gorm.DB) user.IService {
	ret := _m.Called(dbTrx)

	var r0 user.IService
	if rf, ok := ret.Get(0).(func(*gorm.DB) user.IService); ok {
		r0 = rf(dbTrx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(user.IService)
		}
	}

	return r0
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
