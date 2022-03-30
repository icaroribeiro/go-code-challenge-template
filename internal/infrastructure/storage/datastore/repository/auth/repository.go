package auth

import (
	"strings"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	authdsrepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/infrastructure/storage/datastore/repository/auth"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/model"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

var initDB *gorm.DB

// New is the factory function that encapsulates the implementation related to auth.
func New(db *gorm.DB) authdsrepository.IRepository {
	initDB = db
	return &Repository{
		DB: db,
	}
}

// Create is the function that creates an auth in the datastore.
func (r *Repository) Create(auth domainmodel.Auth) (domainmodel.Auth, error) {
	authDB := datastoremodel.Auth{}
	authDB.FromDomain(auth)

	result := r.DB.Create(&authDB)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "auths_user_id_key") {
			loginDB := datastoremodel.Login{}

			if result := r.DB.Find(&loginDB, "user_id=?", authDB.UserID); result.Error != nil {
				return domainmodel.Auth{}, result.Error
			}

			if result.RowsAffected == 0 && loginDB.IsEmpty() {
				return domainmodel.Auth{}, customerror.NotFound.Newf("the login record with user id %s was not found", authDB.UserID)
			}

			return domainmodel.Auth{}, customerror.Conflict.Newf("The user with id %s is already logged in", authDB.UserID)
		}

		return domainmodel.Auth{}, result.Error
	}

	return authDB.ToDomain(), nil
}

// GetByUserID is the function that gets an auth by user id from the datastore.
func (r *Repository) GetByUserID(userID string) (domainmodel.Auth, error) {
	authDB := datastoremodel.Auth{}

	if result := r.DB.Find(&authDB, "user_id=?", userID); result.Error != nil {
		return domainmodel.Auth{}, result.Error
	}

	return authDB.ToDomain(), nil
}

// Delete is the function that deletes an auth by id from the datastore.
func (r *Repository) Delete(id string) (domainmodel.Auth, error) {
	authDB := datastoremodel.Auth{}

	result := r.DB.Find(&authDB, "id=?", id)
	if result.Error != nil {
		return domainmodel.Auth{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainmodel.Auth{}, customerror.NotFound.Newf("the auth with id %s was not found", id)
	}

	if result = r.DB.Delete(&authDB); result.Error != nil {
		return domainmodel.Auth{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainmodel.Auth{}, customerror.NotFound.Newf("the auth with id %s was not deleted", id)
	}

	return authDB.ToDomain(), nil
}

// WithDBTrx is the function that enables the repository with datastore transaction.
func (r *Repository) WithDBTrx(dbTrx *gorm.DB) authdsrepository.IRepository {
	if dbTrx == nil {
		r.DB = initDB
		return r
	}

	r.DB = dbTrx
	return r
}
