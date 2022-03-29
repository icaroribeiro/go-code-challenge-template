package login

// import (
// 	"strings"

// 	loginmodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/login"
// 	logindsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/persistence/datastore/repository/login"
// 	logindsmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/model/login"
// 	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB *gorm.DB
// }

// var initDB *gorm.DB

// // New is the factory function that encapsulates the implementation related to login.
// func New(db *gorm.DB) logindsrepository.IRepository {
// 	initDB = db
// 	return &Repository{
// 		DB: db,
// 	}
// }

// // Create is the function that creates a login in the database.
// func (r *Repository) Create(login loginmodel.Login) (loginmodel.Login, error) {
// 	loginDB := logindsmodel.Login{}
// 	loginDB.FromDomain(login)

// 	if result := r.DB.Create(&loginDB); result.Error != nil {
// 		if strings.Contains(result.Error.Error(), "logins_user_id_key") {
// 			return loginmodel.Login{}, customerror.Conflict.Newf("The user with id %s is already logged in", login.Username)
// 		}

// 		return loginmodel.Login{}, result.Error
// 	}

// 	return loginDB.ToDomain(), nil
// }

// // GetByUsername is the function that gets a user by username from the database.
// func (r *Repository) GetByUsername(username string) (loginmodel.Login, error) {
// 	loginDB := logindsmodel.Login{}

// 	if result := r.DB.Find(&loginDB, "username=?", username); result.Error != nil {
// 		return loginmodel.Login{}, result.Error
// 	}

// 	return loginDB.ToDomain(), nil
// }

// // GetByUsername is the function that gets a user by username from the database.
// func (r *Repository) GetByUserID(userID string) (loginmodel.Login, error) {
// 	loginDB := logindsmodel.Login{}

// 	if result := r.DB.Find(&loginDB, "user_id=?", userID); result.Error != nil {
// 		return loginmodel.Login{}, result.Error
// 	}

// 	return loginDB.ToDomain(), nil
// }

// // Update is the function that updates a login by id in the database.
// func (r *Repository) Update(id string, login loginmodel.Login) (loginmodel.Login, error) {
// 	loginDB := logindsmodel.Login{}
// 	loginDB.FromDomain(login)

// 	result := r.DB.Model(&loginDB).Where("id=?", id).Updates(&loginDB)
// 	if result.Error != nil {
// 		return loginmodel.Login{}, result.Error
// 	}

// 	if result.RowsAffected == 0 {
// 		return loginmodel.Login{}, customerror.NotFound.Newf("the login with id %s was not updated", id)
// 	}

// 	if result = r.DB.Find(&loginDB, "id=?", id); result.Error != nil {
// 		return loginmodel.Login{}, result.Error
// 	}

// 	if result.RowsAffected == 0 {
// 		return loginmodel.Login{}, customerror.NotFound.Newf("the login id %s was not found", id)
// 	}

// 	return loginDB.ToDomain(), nil
// }

// // Delete is the function that deletes a login by id from the database.
// func (r *Repository) Delete(id string) (loginmodel.Login, error) {
// 	loginDB := logindsmodel.Login{}

// 	result := r.DB.Find(&loginDB, "id=?", id)
// 	if result.Error != nil {
// 		return loginmodel.Login{}, result.Error
// 	}

// 	if result.RowsAffected == 0 {
// 		return loginmodel.Login{}, customerror.NotFound.Newf("the login with id %s was not found", id)
// 	}

// 	if result = r.DB.Delete(&loginDB); result.Error != nil {
// 		return loginmodel.Login{}, result.Error
// 	}

// 	if result.RowsAffected == 0 {
// 		return loginmodel.Login{}, customerror.NotFound.Newf("the login with id %s was not deleted", id)
// 	}

// 	return loginDB.ToDomain(), nil
// }

// // WithDBTrx is the function that enables the repository with database transaction.
// func (r *Repository) WithDBTrx(dbTrx *gorm.DB) logindsrepository.IRepository {
// 	if dbTrx == nil {
// 		r.DB = initDB
// 		return r
// 	}

// 	r.DB = dbTrx
// 	return r
// }
