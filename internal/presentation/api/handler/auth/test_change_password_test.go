package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/auth"
	authhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/presentation/api/handler/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/security"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestChangePassword() {
	passwords := securitypkg.Passwords{}

	body := ""

	auth := domainmodel.Auth{}

	authDetailsCtxValue := domainmodel.Auth{}

	dbTrx := &gorm.DB{}
	dbTrx = nil

	message := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInResettingThePassword",
			SetUp: func(t *testing.T) {
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
				{
					"current_password":"%s",
					"new_password":"%s"
				}`,
					passwords.CurrentPassword, passwords.NewPassword)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				authDetailsCtxValue = auth

				message = "the password has been updated successfully"

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
			SetUp: func(t *testing.T) {
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				authDetailsCtxValue = domainmodel.Auth{}

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
			SetUp: func(t *testing.T) {
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
					"current_password":"%s",
					"new_password":"%s"
				`,
					passwords.CurrentPassword, passwords.NewPassword)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				authDetailsCtxValue = auth

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			StatusCode: http.StatusBadRequest,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenResettingThePassword",
			SetUp: func(t *testing.T) {
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
				{
					"current_password":"%s",
					"new_password":"%s"
				}`,
					passwords.CurrentPassword, passwords.NewPassword)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				authDetailsCtxValue = auth

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authService := new(authmockservice.Service)
			authService.On("WithDBTrx", dbTrx).Return(authService)
			authService.On("ModifyPassword", auth.UserID.String(), passwords).Return(returnArgs[0]...)

			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "ChangePassword",
				Method:      http.MethodPost,
				Path:        "/change_password",
				HandlerFunc: authHandler.ChangePassword,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
				Body:   body,
			}

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			ctx := req.Context()
			ctx = authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)
			req = req.WithContext(ctx)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedMessage := messagehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, returnedMessage.Text)
				assert.Equal(t, messagehttputilpkg.Message{Text: message}, returnedMessage)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}
