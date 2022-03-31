package auth

import (
	"strings"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/infrastructure/storage/datastore/repository/auth"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/model"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

var initDB *gorm.DB

// New is the factory function that encapsulates the implementation related to auth.
func New(db *gorm.DB) authdatastorerepository.IRepository {
	initDB = db
	return &Repository{
		DB: db,
	}
}

// Create is the function that creates an auth in the datastore.
func (r *Repository) Create(auth domainmodel.Auth) (domainmodel.Auth, error) {
	authDS := datastoremodel.Auth{}
	authDS.FromDomain(auth)

	result := r.DB.Create(&authDS)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "auths_user_id_key") {
			loginDS := datastoremodel.Login{}

			if result := r.DB.Find(&loginDS, "user_id=?", authDS.UserID); result.Error != nil {
				return domainmodel.Auth{}, result.Error
			}

			if result.RowsAffected == 0 && loginDS.IsEmpty() {
				return domainmodel.Auth{}, customerror.NotFound.Newf("the login record with user id %s was not found", authDS.UserID)
			}

			return domainmodel.Auth{}, customerror.Conflict.Newf("The user with id %s is already logged in", authDS.UserID)
		}

		return domainmodel.Auth{}, result.Error
	}

	return authDS.ToDomain(), nil
}

// GetByUserID is the function that gets an auth by user id from the datastore.
func (r *Repository) GetByUserID(userID string) (domainmodel.Auth, error) {
	authDS := datastoremodel.Auth{}

	if result := r.DB.Find(&authDS, "user_id=?", userID); result.Error != nil {
		return domainmodel.Auth{}, result.Error
	}

	return authDS.ToDomain(), nil
}

// Delete is the function that deletes an auth by id from the datastore.
func (r *Repository) Delete(id string) (domainmodel.Auth, error) {
	authDS := datastoremodel.Auth{}

	result := r.DB.Find(&authDS, "id=?", id)
	if result.Error != nil {
		return domainmodel.Auth{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainmodel.Auth{}, customerror.NotFound.Newf("the auth with id %s was not found", id)
	}

	if result = r.DB.Delete(&authDS); result.Error != nil {
		return domainmodel.Auth{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domainmodel.Auth{}, customerror.NotFound.Newf("the auth with id %s was not deleted", id)
	}

	return authDS.ToDomain(), nil
}

// WithDBTrx is the function that enables the repository with datastore transaction.
func (r *Repository) WithDBTrx(dbTrx *gorm.DB) authdatastorerepository.IRepository {
	if dbTrx == nil {
		r.DB = initDB
		return r
	}

	r.DB = dbTrx
	return r
}
