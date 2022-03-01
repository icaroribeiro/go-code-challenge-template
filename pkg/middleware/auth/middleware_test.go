package auth_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
	mockauthpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/mockauth"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/core/domain/model"
	datastoremodelfactory "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/infrastructure/storage/datastore/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestMiddlewareUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestAuth() {
	bearerToken := []string{"", ""}

	var token *jwt.Token
	timeBeforeExpTimeInSec := 0

	headers := make(map[string][]string)

	statusCode := 0
	payload := messagehttputilpkg.Message{}

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionAndApplyingAuthenticationToARequest",
			SetUp: func(t *testing.T) {
				bearerToken = []string{"Bearer", "token"}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusOK
				payload = messagehttputilpkg.Message{Text: "OK"}

				id := uuid.NewV4()
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     id,
					"userID": userID,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				authFromDB := datastoremodelfactory.NewAuth(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authFromDB.ID, authFromDB.UserID, authFromDB.CreatedAt)

				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheAuthorizationHeaderIsNotSetInTheRequestHeader",
			SetUp: func(t *testing.T) {
				bearerToken = []string{"", ""}

				headers = map[string][]string{}

				statusCode = http.StatusBadRequest

				returnArgs = ReturnArgs{
					{nil, nil},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthenticationTokenIsNotSetInAuthorizationHeader",
			SetUp: func(t *testing.T) {
				bearerToken = []string{"Bearer", ""}

				key := "Authorization"
				value := bearerToken[0]
				headers = map[string][]string{
					key: {value},
				}

				statusCode = http.StatusBadRequest

				returnArgs = ReturnArgs{
					{nil, nil},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheTokenIsNotDecoded",
			SetUp: func(t *testing.T) {
				bearerToken = []string{"Bearer", "token"}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				headers = map[string][]string{
					key: {value},
				}

				statusCode = http.StatusUnauthorized

				returnArgs = ReturnArgs{
					{nil, customerror.New("failed")},
					{domainmodel.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFetchedFromTheToken",
			SetUp: func(t *testing.T) {
				bearerToken = []string{"Bearer", "token"}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}
				statusCode = http.StatusInternalServerError

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodel.Auth{}, customerror.New("failed")},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenTryingToFindTheAuthInTheDatabase",
			SetUp: func(t *testing.T) {
				bearerToken = []string{"Bearer", "token"}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusInternalServerError

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnError(errors.New("failed"))
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFoundInTheDatabase",
			SetUp: func(t *testing.T) {
				bearerToken = []string{"Bearer", "token"}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusBadRequest

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(&sqlmock.Rows{})
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheUserIDFromTokenDoesNotMatchWithTheUserIDFromAuthRecordFromTheDatabase",
			SetUp: func(t *testing.T) {
				bearerToken = []string{"Bearer", "token"}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusBadRequest

				id := uuid.NewV4()

				args := map[string]interface{}{
					"id": id,
				}

				returnArgs = ReturnArgs{
					{token, nil},
					{domainmodelfactory.NewAuth(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				authFromDB := datastoremodelfactory.NewAuth(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authFromDB.ID, authFromDB.UserID, authFromDB.CreatedAt)

				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authN := new(mockauthpkg.MockAuth)
			authN.On("DecodeToken", bearerToken[1]).Return(returnArgs[0]...)
			authN.On("FetchAuthFromToken", token).Return(returnArgs[1]...)

			authMiddleware := authmiddlewarepkg.Auth(ts.DB, authN, timeBeforeExpTimeInSec)

			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: "OK"})
			}

			returnedHandlerFunc := adapterhttputilpkg.AdaptFunc(handlerFunc).With(authMiddleware)

			req := httptest.NewRequest(http.MethodGet, "/testing", nil)

			requesthttputilpkg.SetRequestHeaders(req, headers)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name("testing").
				Methods(http.MethodGet).
				Path("/testing").
				HandlerFunc(returnedHandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Result().Header.Get("Content-Type"), "application/json")
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
				message := messagehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&message)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.Equal(t, payload, message)
			} else {
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
			}
		})
	}
}
