package user

import (
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
)

// IService interface is a collection of function signatures that represents the user's service contract.
type IService interface {
	GetAll() (domainmodel.Users, error)
}
