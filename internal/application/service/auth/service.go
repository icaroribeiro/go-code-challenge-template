package auth

// import (
// 	authmodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/auth"
// 	loginmodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/login"
// 	usermodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/user"
// 	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/application/service/auth"
// 	authdsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/persistence/datastore/repository/auth"
// 	logindsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/persistence/datastore/repository/login"
// 	userdsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/persistence/datastore/repository/user"
// 	authinfra "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/auth"
// 	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
// 	securitypkg "github.com/icaroribeiro/go-code-challenge-template/pkg/security"
// 	validatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator"
// 	"gorm.io/gorm"
// )

// type Service struct {
// 	AuthDatastoreRepository  authdsrepository.IRepository
// 	UserDatastoreRepository  userdsrepository.IRepository
// 	LoginDatastoreRepository logindsrepository.IRepository
// 	AuthInfra                authinfra.IAuth
// 	TokenExpTimeInSec        int
// 	Security                 securitypkg.ISecurity
// 	Validator                validatorpkg.IValidator
// }

// // New is the factory function that encapsulates the implementation related to auth.
// func New(authDatastoreRepository authdsrepository.IRepository, userDatastoreRepository userdsrepository.IRepository, loginDatastoreRepository logindsrepository.IRepository,
// 	authInfra authinfra.IAuth, tokenExpTimeInSec int, security securitypkg.ISecurity, validator validatorpkg.IValidator) authservice.IService {
// 	return &Service{
// 		AuthDatastoreRepository:  authDatastoreRepository,
// 		UserDatastoreRepository:  userDatastoreRepository,
// 		LoginDatastoreRepository: loginDatastoreRepository,
// 		AuthInfra:                authInfra,
// 		TokenExpTimeInSec:        tokenExpTimeInSec,
// 		Security:                 security,
// 		Validator:                validator,
// 	}
// }

// // Register is the function that registers the user to the system.
// func (a *Service) Register(credentials securitypkg.Credentials) (string, error) {
// 	// if err := a.Validator.Validate(credentials, ""); err != nil {
// 	// 	return "", customerror.BadRequest.New(err.Error())
// 	// }

// 	login, err := a.LoginDatastoreRepository.GetByUsername(credentials.Username)
// 	if err != nil {
// 		return "", err
// 	}

// 	if !login.IsEmpty() {
// 		return "", customerror.Conflict.Newf("the username %s already exists", credentials.Username)
// 	}

// 	user := usermodel.User{
// 		Username: credentials.Username,
// 	}

// 	newUser, err := a.UserDatastoreRepository.Create(user)
// 	if err != nil {
// 		return "", err
// 	}

// 	login = loginmodel.Login{
// 		UserID:   newUser.ID,
// 		Username: credentials.Username,
// 		Password: credentials.Password,
// 	}

// 	_, err = a.LoginDatastoreRepository.Create(login)
// 	if err != nil {
// 		return "", err
// 	}

// 	auth := authmodel.Auth{
// 		UserID: newUser.ID,
// 	}

// 	newAuth, err := a.AuthDatastoreRepository.Create(auth)
// 	if err != nil {
// 		return "", err
// 	}

// 	auth = newAuth

// 	token, err := a.AuthInfra.CreateToken(auth, a.TokenExpTimeInSec)
// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

// // LogIn is the function that initializes the user access to the system.
// func (a *Service) LogIn(credentials securitypkg.Credentials) (string, error) {
// 	// if err := a.Validator.Validate(credentials, ""); err != nil {
// 	// 	return "", customerror.BadRequest.New(err.Error())
// 	// }

// 	login, err := a.LoginDatastoreRepository.GetByUsername(credentials.Username)
// 	if err != nil {
// 		return "", err
// 	}

// 	if login.IsEmpty() {
// 		return "", customerror.NotFound.Newf("the username %s was not found", credentials.Username)
// 	}

// 	if err = a.Security.VerifyPasswords(login.Password, credentials.Password); err != nil {
// 		return "", err
// 	}

// 	auth, err := a.AuthDatastoreRepository.GetByUserID(login.UserID.String())
// 	if err != nil {
// 		return "", err
// 	}

// 	if !auth.IsEmpty() {
// 		return "", customerror.New("you are already logged in")
// 	}

// 	auth = authmodel.Auth{
// 		UserID: login.UserID,
// 	}

// 	newAuth, err := a.AuthDatastoreRepository.Create(auth)
// 	if err != nil {
// 		return "", err
// 	}

// 	auth = newAuth

// 	token, err := a.AuthInfra.CreateToken(auth, a.TokenExpTimeInSec)
// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

// // RenewToken is the function that renews the token.
// func (a *Service) RenewToken(auth authmodel.Auth) (string, error) {
// 	return a.AuthInfra.CreateToken(auth, a.TokenExpTimeInSec)
// }

// // ModifyPassword is the function that modifies the user's password.
// func (a *Service) ModifyPassword(id string, passwords securitypkg.Passwords) error {
// 	// if err := a.Validator.Valid(id, "nonzero, uuid"); err != nil {
// 	// 	return customerror.BadRequest.Newf("UserID: %s", err.Error())
// 	// }

// 	// if err := a.Validator.Validate(passwords, ""); err != nil {
// 	// 	return customerror.BadRequest.New(err.Error())
// 	// }

// 	login, err := a.LoginDatastoreRepository.GetByUserID(id)
// 	if err != nil {
// 		return err
// 	}

// 	if login.IsEmpty() {
// 		return customerror.NotFound.New("the user who owns this token was not found")
// 	}

// 	if err = a.Security.VerifyPasswords(login.Password, passwords.CurrentPassword); err != nil {
// 		if err.Error() == "the password is invalid" {
// 			return customerror.Unauthorized.New("the current password did not match the one already registered")
// 		}

// 		return err
// 	}

// 	if passwords.NewPassword == passwords.CurrentPassword {
// 		return customerror.BadRequest.New("the new password is the same as the one currently registered")
// 	}

// 	login.Password = passwords.NewPassword

// 	_, err = a.LoginDatastoreRepository.Update(login.ID.String(), login)

// 	return err
// }

// // LogOut is the function that concludes the user access to the system.
// func (a *Service) LogOut(id string) error {
// 	if err := a.Validator.Valid(id, "nonzero, uuid"); err != nil {
// 		return customerror.BadRequest.New(err.Error())
// 	}

// 	_, err := a.AuthDatastoreRepository.Delete(id)

// 	return err
// }

// // WithDBTrx is the function that enables the service with database transaction.
// func (a *Service) WithDBTrx(dbTrx *gorm.DB) authservice.IService {
// 	a.AuthDatastoreRepository = a.AuthDatastoreRepository.WithDBTrx(dbTrx)

// 	a.UserDatastoreRepository = a.UserDatastoreRepository.WithDBTrx(dbTrx)

// 	a.LoginDatastoreRepository = a.LoginDatastoreRepository.WithDBTrx(dbTrx)

// 	return a
// }
