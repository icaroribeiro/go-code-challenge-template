package auth_test

// import (
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"time"

// 	fake "github.com/brianvoe/gofakeit/v5"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/application/customerror"
// 	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/auth"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/application/validatormock"
// 	authmodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/auth"
// 	loginmodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/login"
// 	authdbrepositorymock "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/repositorymock/database/auth"
// 	logindbrepositorymock "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/repositorymock/database/login"
// 	userdbrepositorymock "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/repositorymock/database/user"
// 	authmock "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/authmock"
// 	authdbmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/database/model/auth"
// 	logindbmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/database/model/login"
// 	datastoremodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/database/model/user"
// 	authdbrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/database/repository/auth"
// 	logindbrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/database/repository/login"
// 	userdbrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/database/repository/user"
// 	security "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/security"
// 	securitymock "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/securitymock"
// 	authdbmodelfactory "github.com/icaroribeiro/go-code-challenge-template/tests/factory/infrastructure/persistence/database/model/auth"
// 	logindbmodelfactory "github.com/icaroribeiro/go-code-challenge-template/tests/factory/infrastructure/persistence/database/model/login"
// 	datastoremodelfactory "github.com/icaroribeiro/go-code-challenge-template/tests/factory/infrastructure/persistence/database/model/user"
// 	securityfactory "github.com/icaroribeiro/go-code-challenge-template/tests/factory/infrastructure/security"
// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// func TestService(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestRegister() {
// 	var credentials security.Credentials

// 	var userdb datastoremodel.User

// 	var logindb logindbmodel.Login

// 	var authdb authdbmodel.Auth

// 	var newAuth authmodel.Auth

// 	tokenExpTimeInSec := fake.Number(2, 10)

// 	token := ""

// 	errorType := customerror.NoType

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInRegisteringAUser",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				args := map[string]interface{}{
// 					"username": username,
// 					"password": password,
// 				}

// 				credentials = securityfactory.NewCredentials(args)

// 				args = map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"username":  credentials.Username,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				userdb = datastoremodelfactory.New(args)

// 				id := uuid.NewV4()

// 				args = map[string]interface{}{
// 					"id":        id,
// 					"username":  credentials.Username,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				newUserdb := datastoremodelfactory.New(args)

// 				args = map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    newUserdb.ID,
// 					"username":  credentials.Username,
// 					"password":  credentials.Password,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				logindb = logindbmodelfactory.New(args)

// 				id = uuid.NewV4()

// 				args = map[string]interface{}{
// 					"id":        id,
// 					"userID":    newUserdb.ID,
// 					"username":  credentials.Username,
// 					"password":  credentials.Password,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				newLogindb := logindbmodelfactory.New(args)

// 				args = map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    newUserdb.ID,
// 					"createdAt": time.Time{},
// 				}

// 				authdb = authdbmodelfactory.New(args)

// 				id = uuid.NewV4()

// 				args = map[string]interface{}{
// 					"id":        id,
// 					"userID":    newUserdb.ID,
// 					"createdAt": time.Time{},
// 				}

// 				newAuthdb := authdbmodelfactory.New(args)

// 				newAuth = newAuthdb.ToDomain()

// 				token = fake.Word()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 					{newUserdb, nil},
// 					{newLogindb, nil},
// 					{newAuthdb, nil},
// 					{token, nil},
// 				}
// 			},
// 			WantError: false,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheCredentialsAreNotValid",
// 			SetUp: func(t *testing.T) {
// 				args := map[string]interface{}{
// 					"username": "",
// 					"password": "",
// 				}

// 				credentials = securityfactory.NewCredentials(args)

// 				returnArgs = ReturnArgs{
// 					{errors.New("failed")},
// 					{logindbmodel.Login{}, nil},
// 					{datastoremodel.User{}, nil},
// 					{logindbmodel.Login{}, nil},
// 					{authdbmodel.Auth{}, nil},
// 					{"", nil},
// 				}

// 				errorType = customerror.BadRequest
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		// {
// 		// 	Context: "ItShouldFailIfAnErrorOccursWhenGettingALoginByUsername",
// 		// 	SetUp: func(t *testing.T) {
// 		// 		username := fake.Username()
// 		// 		password := fake.Password(true, true, true, false, false, 8)

// 		// 		credentials = security.Credentials{
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		returnArgs = ReturnArgs{
// 		// 			{nil},
// 		// 			{logindbmodel.Login{}, errors.New("failed")},
// 		// 			{datastoremodel.User{}, nil},
// 		// 			{logindbmodel.Login{}, nil},
// 		// 			{authdbmodel.Auth{}, nil},
// 		// 			{"", nil},
// 		// 		}

// 		// 		errorType = customerror.NoType
// 		// 	},
// 		// 	WantError: true,
// 		// 	TearDown:  func(t *testing.T) {},
// 		// },
// 		// {
// 		// 	Context: "ItShouldFailIfTheUsernameIsAlreadyRegistered",
// 		// 	SetUp: func(t *testing.T) {
// 		// 		username := fake.Username()
// 		// 		password := fake.Password(true, true, true, false, false, 8)

// 		// 		credentials = security.Credentials{
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		id := uuid.NewV4()
// 		// 		userID := uuid.NewV4()

// 		// 		logindb = logindbmodel.Login{
// 		// 			ID:       id,
// 		// 			UserID:   userID,
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		returnArgs = ReturnArgs{
// 		// 			{nil},
// 		// 			{logindb, nil},
// 		// 			{datastoremodel.User{}, nil},
// 		// 			{logindbmodel.Login{}, nil},
// 		// 			{authdbmodel.Auth{}, nil},
// 		// 			{"", nil},
// 		// 		}

// 		// 		errorType = customerror.Conflict
// 		// 	},
// 		// 	WantError: true,
// 		// 	TearDown:  func(t *testing.T) {},
// 		// },
// 		// {
// 		// 	Context: "ItShouldFailIfAnErrorOccursWhenCreatingAUser",
// 		// 	SetUp: func(t *testing.T) {
// 		// 		username := fake.Username()
// 		// 		password := fake.Password(true, true, true, false, false, 8)

// 		// 		credentials = security.Credentials{
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		userdb = datastoremodel.User{
// 		// 			Username: username,
// 		// 		}

// 		// 		returnArgs = ReturnArgs{
// 		// 			{nil},
// 		// 			{logindbmodel.Login{}, nil},
// 		// 			{datastoremodel.User{}, errors.New("failed")},
// 		// 			{logindbmodel.Login{}, nil},
// 		// 			{authdbmodel.Auth{}, nil},
// 		// 			{"", nil},
// 		// 		}

// 		// 		errorType = customerror.NoType
// 		// 	},
// 		// 	WantError: true,
// 		// 	TearDown:  func(t *testing.T) {},
// 		// },
// 		// {
// 		// 	Context: "ItShouldFailIfAnErrorOccursWhenCreatingALogin",
// 		// 	SetUp: func(t *testing.T) {
// 		// 		username := fake.Username()
// 		// 		password := fake.Password(true, true, true, false, false, 8)

// 		// 		credentials = security.Credentials{
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		userdb = datastoremodel.User{
// 		// 			Username: username,
// 		// 		}

// 		// 		id := uuid.NewV4()

// 		// 		newUserdb := datastoremodel.User{
// 		// 			ID:       id,
// 		// 			Username: username,
// 		// 		}

// 		// 		logindb = logindbmodel.Login{
// 		// 			UserID:   newUserdb.ID,
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		returnArgs = ReturnArgs{
// 		// 			{nil},
// 		// 			{logindbmodel.Login{}, nil},
// 		// 			{newUserdb, nil},
// 		// 			{logindbmodel.Login{}, errors.New("failed")},
// 		// 			{authdbmodel.Auth{}, nil},
// 		// 			{"", nil},
// 		// 		}

// 		// 		errorType = customerror.NoType
// 		// 	},
// 		// 	WantError: true,
// 		// 	TearDown:  func(t *testing.T) {},
// 		// },
// 		// {
// 		// 	Context: "ItShouldFailIfAnErrorOccursWhenCreatingAnAuth",
// 		// 	SetUp: func(t *testing.T) {
// 		// 		username := fake.Username()
// 		// 		password := fake.Password(true, true, true, false, false, 8)

// 		// 		credentials = security.Credentials{
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		userdb = datastoremodel.User{
// 		// 			Username: username,
// 		// 		}

// 		// 		id := uuid.NewV4()

// 		// 		newUserdb := datastoremodel.User{
// 		// 			ID:       id,
// 		// 			Username: username,
// 		// 		}

// 		// 		logindb = logindbmodel.Login{
// 		// 			UserID:   newUserdb.ID,
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		id = uuid.NewV4()

// 		// 		newLogindb := logindbmodel.Login{
// 		// 			ID:       id,
// 		// 			UserID:   newUserdb.ID,
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		authdb = authdbmodel.Auth{
// 		// 			UserID: newUserdb.ID,
// 		// 		}

// 		// 		returnArgs = ReturnArgs{
// 		// 			{nil},
// 		// 			{logindbmodel.Login{}, nil},
// 		// 			{newUserdb, nil},
// 		// 			{newLogindb, nil},
// 		// 			{authdbmodel.Auth{}, errors.New("failed")},
// 		// 			{"", nil},
// 		// 		}

// 		// 		errorType = customerror.NoType
// 		// 	},
// 		// 	WantError: true,
// 		// 	TearDown:  func(t *testing.T) {},
// 		// },
// 		// {
// 		// 	Context: "ItShouldFailIfAnErrorOccursWhenCreatingAToken",
// 		// 	SetUp: func(t *testing.T) {
// 		// 		username := fake.Username()
// 		// 		password := fake.Password(true, true, true, false, false, 8)

// 		// 		credentials = security.Credentials{
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		userdb = datastoremodel.User{
// 		// 			Username: username,
// 		// 		}

// 		// 		id := uuid.NewV4()

// 		// 		newUserdb := datastoremodel.User{
// 		// 			ID:       id,
// 		// 			Username: username,
// 		// 		}

// 		// 		logindb = logindbmodel.Login{
// 		// 			UserID:   newUserdb.ID,
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		id = uuid.NewV4()

// 		// 		newLogindb := logindbmodel.Login{
// 		// 			ID:       id,
// 		// 			UserID:   newUserdb.ID,
// 		// 			Username: username,
// 		// 			Password: password,
// 		// 		}

// 		// 		authdb = authdbmodel.Auth{
// 		// 			UserID: newUserdb.ID,
// 		// 		}

// 		// 		id = uuid.NewV4()

// 		// 		newAuthdb = authdbmodel.Auth{
// 		// 			ID:     id,
// 		// 			UserID: newUserdb.ID,
// 		// 		}

// 		// 		newAuth = newAuthdb.ToDomain()

// 		// 		returnArgs = ReturnArgs{
// 		// 			{nil},
// 		// 			{logindbmodel.Login{}, nil},
// 		// 			{newUserdb, nil},
// 		// 			{newLogindb, nil},
// 		// 			{newAuthdb, nil},
// 		// 			{"", errors.New("failed")},
// 		// 		}

// 		// 		errorType = customerror.NoType
// 		// 	},
// 		// 	WantError: true,
// 		// 	TearDown:  func(t *testing.T) {},
// 		// },
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			validator := new(validatormock.ValidatorMock)
// 			validator.On("Validate", credentials, "").Return(returnArgs[0]...)

// 			loginDatastoreRepository := new(logindbrepositorymock.RepositoryMock)
// 			loginDatastoreRepository.On("GetByUsername", credentials.Username).Return(returnArgs[1]...)

// 			userDatastoreRepository := new(userdbrepositorymock.RepositoryMock)
// 			userDatastoreRepository.On("Create", userdb).Return(returnArgs[2]...)

// 			loginDatastoreRepository.On("Create", logindb).Return(returnArgs[3]...)

// 			authDatastoreRepository := new(authdbrepositorymock.RepositoryMock)
// 			authDatastoreRepository.On("Create", authdb).Return(returnArgs[4]...)

// 			authN := new(authmock.AuthMock)
// 			authN.On("CreateToken", newAuth, tokenExpTimeInSec).Return(returnArgs[5]...)

// 			security := new(securitymock.SecurityMock)

// 			authService := authservice.New(authDatastoreRepository, userDatastoreRepository, loginDatastoreRepository,
// 				authN, tokenExpTimeInSec, security, validator)

// 			returnedToken, err := authService.Register(credentials)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, token, returnedToken)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedToken)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestLogIn() {
// 	credentials := security.Credentials{}

// 	logindb := logindbmodel.Login{}

// 	login := loginmodel.Login{}

// 	authdb := authdbmodel.Auth{}

// 	newAuthdb := authdbmodel.Auth{}

// 	newAuth := authmodel.Auth{}

// 	tokenExpTimeInSec := fake.Number(2, 10)

// 	token := ""

// 	errorType := customerror.NoType

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInLoggingIn",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				logindb = logindbmodel.Login{
// 					ID:       id,
// 					UserID:   userID,
// 					Username: username,
// 					Password: password,
// 				}

// 				login = logindb.ToDomain()

// 				authdb = authdbmodel.Auth{
// 					UserID: logindb.UserID,
// 				}

// 				id = uuid.NewV4()

// 				newAuthdb = authdbmodel.Auth{
// 					ID:     id,
// 					UserID: logindb.UserID,
// 				}

// 				newAuth = newAuthdb.ToDomain()

// 				token = fake.Word()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{logindb, nil},
// 					{nil},
// 					{newAuthdb, nil},
// 					{token, nil},
// 				}
// 			},
// 			WantError: false,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheEvaluatedCredentialsValuesAreNotValid",
// 			SetUp: func(t *testing.T) {
// 				credentials = security.Credentials{}

// 				returnArgs = ReturnArgs{
// 					{errors.New("failed")},
// 					{logindbmodel.Login{}, nil},
// 					{nil},
// 					{authdbmodel.Auth{}, nil},
// 					{"", nil},
// 				}

// 				errorType = customerror.BadRequest
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenGettingALoginByUsername",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{logindbmodel.Login{}, errors.New("failed")},
// 					{nil},
// 					{authdbmodel.Auth{}, nil},
// 					{"", nil},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheUsernameIsNotRegistered",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 					{nil},
// 					{authdbmodel.Auth{}, nil},
// 					{"", nil},
// 				}

// 				errorType = customerror.NotFound
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenVerifyingThePasswords",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				logindb = logindbmodel.Login{
// 					ID:       id,
// 					UserID:   userID,
// 					Username: credentials.Username,
// 					Password: credentials.Password,
// 				}

// 				login = logindb.ToDomain()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{logindb, nil},
// 					{errors.New("failed")},
// 					{authdbmodel.Auth{}, nil},
// 					{"", nil},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAnAuth",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				logindb = logindbmodel.Login{
// 					ID:       id,
// 					UserID:   userID,
// 					Username: credentials.Username,
// 					Password: credentials.Password,
// 				}

// 				login = logindb.ToDomain()

// 				authdb = authdbmodel.Auth{
// 					UserID: login.UserID,
// 				}

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{logindb, nil},
// 					{nil},
// 					{authdbmodel.Auth{}, errors.New("failed")},
// 					{"", nil},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAToken",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				logindb = logindbmodel.Login{
// 					ID:       id,
// 					UserID:   userID,
// 					Username: credentials.Username,
// 					Password: credentials.Password,
// 				}

// 				login = logindb.ToDomain()

// 				authdb = authdbmodel.Auth{
// 					UserID: logindb.UserID,
// 				}

// 				id = uuid.NewV4()

// 				newAuthdb = authdbmodel.Auth{
// 					ID:     id,
// 					UserID: login.UserID,
// 				}

// 				newAuth = newAuthdb.ToDomain()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{logindb, nil},
// 					{nil},
// 					{newAuthdb, nil},
// 					{"", errors.New("failed")},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			validator := new(validatormock.ValidatorMock)
// 			validator.On("Validate", credentials, "").Return(returnArgs[0]...)

// 			loginDatastoreRepository := new(logindbrepositorymock.RepositoryMock)
// 			loginDatastoreRepository.On("GetByUsername", credentials.Username).Return(returnArgs[1]...)

// 			security := new(securitymock.SecurityMock)
// 			security.On("VerifyPasswords", login.Password, credentials.Password).Return(returnArgs[2]...)

// 			authDatastoreRepository := new(authdbrepositorymock.RepositoryMock)
// 			authDatastoreRepository.On("Create", authdb).Return(returnArgs[3]...)

// 			authN := new(authmock.AuthMock)
// 			authN.On("CreateToken", newAuth, tokenExpTimeInSec).Return(returnArgs[4]...)

// 			userDatastoreRepository := new(userdbrepositorymock.RepositoryMock)

// 			authService := authservice.New(authDatastoreRepository, userDatastoreRepository, loginDatastoreRepository,
// 				authN, tokenExpTimeInSec, security, validator)

// 			returnedToken, err := authService.LogIn(credentials)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, token, returnedToken)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedToken)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestRenewToken() {
// 	auth := authmodel.Auth{}

// 	tokenExpTimeInSec := fake.Number(2, 10)

// 	token := ""

// 	errorType := customerror.NoType

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInRenewingTheToken",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				token = fake.Word()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{token, nil},
// 				}
// 			},
// 			WantError: false,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAToken",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{"", errors.New("failed")},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			validator := new(validatormock.ValidatorMock)
// 			validator.On("Validate", auth, "").Return(returnArgs[0]...)

// 			authN := new(authmock.AuthMock)
// 			authN.On("CreateToken", auth, tokenExpTimeInSec).Return(returnArgs[1]...)

// 			authDatastoreRepository := new(authdbrepositorymock.RepositoryMock)
// 			userDatastoreRepository := new(userdbrepositorymock.RepositoryMock)
// 			loginDatastoreRepository := new(logindbrepositorymock.RepositoryMock)
// 			security := new(securitymock.SecurityMock)

// 			authService := authservice.New(authDatastoreRepository, userDatastoreRepository, loginDatastoreRepository,
// 				authN, tokenExpTimeInSec, security, validator)

// 			returnedToken, err := authService.RenewToken(auth)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, token, returnedToken)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedToken)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestModifyPassword() {
// 	id := ""

// 	passwords := security.Passwords{}

// 	logindb := logindbmodel.Login{}

// 	login := loginmodel.Login{}

// 	updatedLogindb := logindbmodel.Login{}

// 	errorType := customerror.NoType

// 	tokenExpTimeInSec := fake.Number(2, 10)

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInModifyingThePassword",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				loginID := uuid.NewV4()
// 				userID := uuid.NewV4()
// 				username := fake.Username()

// 				logindb = logindbmodel.Login{
// 					ID:       loginID,
// 					UserID:   userID,
// 					Username: username,
// 					Password: currentPassword,
// 				}

// 				login = logindb.ToDomain()

// 				updatedLogindb = logindb
// 				updatedLogindb.Password = newPassword

// 				newLogindb := updatedLogindb

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{nil},
// 					{logindb, nil},
// 					{nil},
// 					{newLogindb, nil},
// 				}
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheIDIsNotValid",
// 			SetUp: func(t *testing.T) {
// 				id = ""

// 				returnArgs = ReturnArgs{
// 					{errors.New("failed")},
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 				}

// 				errorType = customerror.BadRequest
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheEvaluatedPasswordsValuesAreNotValid",
// 			SetUp: func(t *testing.T) {
// 				passwords = security.Passwords{}

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{errors.New("failed")},
// 					{logindbmodel.Login{}, nil},
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 				}

// 				errorType = customerror.BadRequest
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenGettingALoginByUsername",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{nil},
// 					{logindbmodel.Login{}, errors.New("failed")},
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheIDIsNotRegistered",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 				}

// 				errorType = customerror.NotFound
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOfInvalidPasswordHappensWhenVerifyingThePasswords",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				loginID := uuid.NewV4()
// 				userID := uuid.NewV4()
// 				username := fake.Username()

// 				logindb = logindbmodel.Login{
// 					ID:       loginID,
// 					UserID:   userID,
// 					Username: username,
// 					Password: currentPassword,
// 				}

// 				login = logindb.ToDomain()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{nil},
// 					{logindb, nil},
// 					{errors.New("the password is invalid")},
// 					{logindbmodel.Login{}, nil},
// 				}

// 				errorType = customerror.Unauthorized
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnotherErrorHappensWhenVerifyingThePasswords",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				loginID := uuid.NewV4()
// 				userID := uuid.NewV4()
// 				username := fake.Username()

// 				logindb = logindbmodel.Login{
// 					ID:       loginID,
// 					UserID:   userID,
// 					Username: username,
// 					Password: currentPassword,
// 				}

// 				login = logindb.ToDomain()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{nil},
// 					{logindb, nil},
// 					{errors.New("failed")},
// 					{logindbmodel.Login{}, nil},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheNewPasswordISTheSameAsTheOneCurrentlyRegistered",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := currentPassword

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				loginID := uuid.NewV4()
// 				userID := uuid.NewV4()
// 				username := fake.Username()

// 				logindb = logindbmodel.Login{
// 					ID:       loginID,
// 					UserID:   userID,
// 					Username: username,
// 					Password: currentPassword,
// 				}

// 				login = logindb.ToDomain()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{nil},
// 					{logindb, nil},
// 					{nil},
// 					{logindbmodel.Login{}, nil},
// 				}

// 				errorType = customerror.BadRequest
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenUpdatingTheLogin",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				loginID := uuid.NewV4()
// 				userID := uuid.NewV4()
// 				username := fake.Username()

// 				logindb = logindbmodel.Login{
// 					ID:       loginID,
// 					UserID:   userID,
// 					Username: username,
// 					Password: currentPassword,
// 				}

// 				login = logindb.ToDomain()

// 				updatedLogindb = logindb
// 				updatedLogindb.Password = newPassword

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{nil},
// 					{logindb, nil},
// 					{nil},
// 					{logindbmodel.Login{}, errors.New("failed")},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			validator := new(validatormock.ValidatorMock)
// 			validator.On("Valid", id, "nonzero, uuid").Return(returnArgs[0]...)
// 			validator.On("Validate", passwords, "").Return(returnArgs[1]...)

// 			loginDatastoreRepository := new(logindbrepositorymock.RepositoryMock)
// 			loginDatastoreRepository.On("GetByUserID", id).Return(returnArgs[2]...)

// 			security := new(securitymock.SecurityMock)
// 			security.On("VerifyPasswords", login.Password, passwords.CurrentPassword).Return(returnArgs[3]...)

// 			loginDatastoreRepository.On("Update", updatedLogindb.ID.String(), updatedLogindb).Return(returnArgs[4]...)

// 			authDatastoreRepository := new(authdbrepositorymock.RepositoryMock)
// 			userDatastoreRepository := new(userdbrepositorymock.RepositoryMock)
// 			authN := new(authmock.AuthMock)

// 			authService := authservice.New(authDatastoreRepository, userDatastoreRepository, loginDatastoreRepository,
// 				authN, tokenExpTimeInSec, security, validator)

// 			err := authService.ModifyPassword(id, passwords)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestLogOut() {
// 	var id string

// 	errorType := customerror.NoType

// 	tokenExpTimeInSec := fake.Number(2, 10)

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInLoggingOut",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()

// 				authID, err := uuid.FromString(id)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				userID := uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":     authID,
// 					"userID": userID,
// 				}

// 				authdb := authdbmodelfactory.New(args)

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{authdb, nil},
// 				}
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheIDIsNotValid",
// 			SetUp: func(t *testing.T) {
// 				id = ""

// 				returnArgs = ReturnArgs{
// 					{errors.New("failed")},
// 					{authdbmodel.Auth{}, nil},
// 				}

// 				errorType = customerror.BadRequest
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenDeletingTheAuth",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4().String()

// 				returnArgs = ReturnArgs{
// 					{nil},
// 					{authdbmodel.Auth{}, errors.New("failed")},
// 				}

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			validator := new(validatormock.ValidatorMock)
// 			validator.On("Valid", id, "nonzero, uuid").Return(returnArgs[0]...)

// 			authDatastoreRepository := new(authdbrepositorymock.RepositoryMock)
// 			authDatastoreRepository.On("Delete", id).Return(returnArgs[1]...)

// 			userDatastoreRepository := new(userdbrepositorymock.RepositoryMock)
// 			loginDatastoreRepository := new(logindbrepositorymock.RepositoryMock)
// 			authN := new(authmock.AuthMock)
// 			security := new(securitymock.SecurityMock)

// 			authService := authservice.New(authDatastoreRepository, userDatastoreRepository, loginDatastoreRepository,
// 				authN, tokenExpTimeInSec, security, validator)

// 			err := authService.LogOut(id)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestWithDBTrx() {
// 	dbTrx := &gorm.DB{}

// 	authDatastoreRepositoryWithDBTrx := &authdbrepository.Repository{}
// 	userDatastoreRepositoryWithDBTrx := &userdbrepository.Repository{}
// 	loginDatastoreRepositoryWithDBTrx := &logindbrepository.Repository{}

// 	tokenExpTimeInSec := fake.Number(2, 10)

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSettingRepositoriesWithDatabaseTransaction",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()

// 				authDatastoreRepositoryWithDBTrx = &authdbrepository.Repository{}
// 				userDatastoreRepositoryWithDBTrx = &userdbrepository.Repository{}
// 				loginDatastoreRepositoryWithDBTrx = &logindbrepository.Repository{}

// 				returnArgs = ReturnArgs{
// 					{authDatastoreRepositoryWithDBTrx},
// 					{userDatastoreRepositoryWithDBTrx},
// 					{loginDatastoreRepositoryWithDBTrx},
// 				}
// 			},
// 			WantError: false,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDatastoreRepository := new(authdbrepositorymock.RepositoryMock)
// 			authDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[0]...)

// 			userDatastoreRepository := new(userdbrepositorymock.RepositoryMock)
// 			userDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[1]...)

// 			loginDatastoreRepository := new(logindbrepositorymock.RepositoryMock)
// 			loginDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[2]...)

// 			authN := new(authmock.AuthMock)
// 			security := new(securitymock.SecurityMock)
// 			validator := new(validatormock.ValidatorMock)

// 			authService := authservice.New(authDatastoreRepository, userDatastoreRepository, loginDatastoreRepository,
// 				authN, tokenExpTimeInSec, security, validator)

// 			returnedAuthService := authService.WithDBTrx(dbTrx)

// 			if !tc.WantError {
// 				assert.NotEmpty(t, returnedAuthService, "Service interface is empty.")
// 				assert.Equal(t, authService, returnedAuthService, "Service interfaces are not the same.")
// 			}
// 		})
// 	}
// }
