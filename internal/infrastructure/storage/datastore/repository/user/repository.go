package user

import (
	"strings"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/infrastructure/storage/datastore/repository/user"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/model"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

var initDB *gorm.DB

// New is the factory function that encapsulates the implementation related to user repository.
func New(db *gorm.DB) userdatastorerepository.IRepository {
	initDB = db
	return &Repository{
		DB: db,
	}
}

// Create is the function that creates a user in the database.
func (r *Repository) Create(user domainmodel.User) (domainmodel.User, error) {
	userDS := datastoremodel.User{}
	userDS.FromDomain(user)

	if result := r.DB.Create(&userDS); result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			return domainmodel.User{}, customerror.Conflict.New(result.Error.Error())
		}

		return domainmodel.User{}, result.Error
	}

	return userDS.ToDomain(), nil
}

// GetAll is the function that gets the list of all users from the database.
func (r *Repository) GetAll() (domainmodel.Users, error) {
	usersDS := datastoremodel.Users{}

	if result := r.DB.Find(&usersDS); result.Error != nil {
		return domainmodel.Users{}, result.Error
	}

	return usersDS.ToDomain(), nil
}

// WithDBTrx is the function that enables the repository with database transaction.
func (r *Repository) WithDBTrx(dbTrx *gorm.DB) userdatastorerepository.IRepository {
	if dbTrx == nil {
		r.DB = initDB
		return r
	}

	r.DB = dbTrx
	return r
}
