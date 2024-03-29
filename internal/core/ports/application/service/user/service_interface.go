package user

import (
	domainentity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	"gorm.io/gorm"
)

// IService interface is a collection of function signatures that represents the user's service contract.
type IService interface {
	GetAll() (domainentity.Users, error)
	WithDBTrx(dbTrx *gorm.DB) IService
}
