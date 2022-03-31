package user

import (
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/service/user"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/infrastructure/storage/datastore/repository/user"
	validatorpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator"
)

type Service struct {
	UserDatastoreRepository userdatastorerepository.IRepository
	Validator               validatorpkg.IValidator
}

// New is the factory function that encapsulates the implementation related to user service.
func New(userDatastoreRepository userdatastorerepository.IRepository, validator validatorpkg.IValidator) userservice.IService {
	return &Service{
		UserDatastoreRepository: userDatastoreRepository,
	}
}

// GetAll is the function that deals with the user repository for getting all users.
func (u *Service) GetAll() (domainmodel.Users, error) {
	users, err := u.UserDatastoreRepository.GetAll()
	if err != nil {
		return domainmodel.Users{}, err
	}

	return users, nil
}
