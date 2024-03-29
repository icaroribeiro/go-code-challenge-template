package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	domainentity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	persistententity "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/perentity"
	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
	adapterhttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/adapter"
	requesthttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/response"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/auth"
	mockauthpkg "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockauth"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestAuth() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	bearerToken := []string{"", ""}

	tokenString := ""

	authHeaderString := ""

	token := &jwt.Token{}

	headers := make(map[string][]string)

	statusCode := 0
	payload := responsehttputilpkg.Message{}

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithAuthenticationMiddleware",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				authHeaderString = value
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusOK

				payload = responsehttputilpkg.Message{Text: "ok"}

				id := uuid.NewV4()
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     id,
					"userID": userID,
				}

				returnArgs = ReturnArgs{
					{tokenString, nil},
					{token, nil},
					{domainentity.AuthFactory(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				persistentAuth := persistententity.AuthFactory(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(persistentAuth.ID, persistentAuth.UserID, persistentAuth.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheAuthorizationHeaderIsNotSetInTheRequestHeader",
			SetUp: func(t *testing.T) {
				tokenString = ""

				bearerToken = []string{"", tokenString}

				authHeaderString = ""

				headers = map[string][]string{}

				token = &jwt.Token{}

				statusCode = http.StatusInternalServerError

				returnArgs = ReturnArgs{
					{tokenString, customerror.New("failed")},
					{token, nil},
					{domainentity.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthenticationTokenIsNotSetInAuthorizationHeader",
			SetUp: func(t *testing.T) {
				tokenString = ""

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := bearerToken[0]
				authHeaderString = value
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusInternalServerError

				returnArgs = ReturnArgs{
					{tokenString, customerror.New("failed")},
					{token, nil},
					{domainentity.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheTokenIsNotDecoded",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				authHeaderString = value
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusUnauthorized

				returnArgs = ReturnArgs{
					{tokenString, nil},
					{token, customerror.New("failed")},
					{domainentity.Auth{}, nil},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFetchedFromTheToken",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				authHeaderString = value
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusInternalServerError

				returnArgs = ReturnArgs{
					{tokenString, nil},
					{token, nil},
					{domainentity.Auth{}, customerror.New("failed")},
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenTryingToFindTheAuthInTheDatabase",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				authHeaderString = value
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
					{tokenString, nil},
					{token, nil},
					{domainentity.AuthFactory(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnError(customerror.New("failed"))
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheAuthIsNotFoundInTheDatabase",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				authHeaderString = value
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
					{tokenString, nil},
					{token, nil},
					{domainentity.AuthFactory(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(&sqlmock.Rows{})
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheUserIDFromTokenDoesNotMatchWithTheUserIDFromAuthRecordFromTheDatabase",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				authHeaderString = value
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
					{tokenString, nil},
					{token, nil},
					{domainentity.AuthFactory(args), nil},
				}

				sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

				persistentAuth := persistententity.AuthFactory(args)

				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(persistentAuth.ID, persistentAuth.UserID, persistentAuth.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authN := new(mockauthpkg.Auth)
			authN.On("ExtractTokenString", authHeaderString).Return(returnArgs[0]...)
			authN.On("DecodeToken", bearerToken[1]).Return(returnArgs[1]...)
			authN.On("FetchAuthFromToken", token).Return(returnArgs[2]...)

			authMiddleware := authmiddlewarepkg.Auth(db, authN)

			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				responsehttputilpkg.RespondWithJSON(w, http.StatusOK, responsehttputilpkg.Message{Text: "ok"})
			}

			returnedHandlerFunc := adapterhttputilpkg.AdaptFunc(handlerFunc).With(authMiddleware)

			route := routehttputilpkg.Route{
				Name:        "Testing",
				Method:      http.MethodGet,
				Path:        "/testing",
				HandlerFunc: returnedHandlerFunc,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			requesthttputilpkg.SetRequestHeaders(req, headers)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Result().Header.Get("Content-Type"), "application/json")
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
				returnedMessage := responsehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.Equal(t, payload, returnedMessage)
			} else {
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
			}

			err := mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}
