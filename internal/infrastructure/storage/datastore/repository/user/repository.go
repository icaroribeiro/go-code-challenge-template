package user

// import (
// 	"strings"

// 	usermodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/user"
// 	userdsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/persistence/datastore/repository/user"
// 	userdsmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/model/user"
// 	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB *gorm.DB
// }

// var initDB *gorm.DB

// // New is the factory function that encapsulates the implementation related to user repository.
// func New(db *gorm.DB) userdsrepository.IRepository {
// 	initDB = db
// 	return &Repository{
// 		DB: db,
// 	}
// }

// // Create is the function that creates a user in the database.
// func (r *Repository) Create(user usermodel.User) (usermodel.User, error) {
// 	userDB := userdsmodel.User{}
// 	userDB.FromDomain(user)

// 	if result := r.DB.Create(&userDB); result.Error != nil {
// 		if strings.Contains(result.Error.Error(), "duplicate key value") {
// 			return usermodel.User{}, customerror.Conflict.New(result.Error.Error())
// 		}

// 		return usermodel.User{}, result.Error
// 	}

// 	return userDB.ToDomain(), nil
// }

// // GetAll is the function that gets the list of all users from the database.
// func (r *Repository) GetAll() (usermodel.Users, error) {
// 	usersDB := userdsmodel.Users{}

// 	if result := r.DB.Find(&usersDB); result.Error != nil {
// 		return usermodel.Users{}, result.Error
// 	}

// 	return usersDB.ToDomain(), nil
// }

// // WithDBTrx is the function that enables the repository with database transaction.
// func (r *Repository) WithDBTrx(dbTrx *gorm.DB) userdsrepository.IRepository {
// 	if dbTrx == nil {
// 		r.DB = initDB
// 		return r
// 	}

// 	r.DB = dbTrx
// 	return r
// }
