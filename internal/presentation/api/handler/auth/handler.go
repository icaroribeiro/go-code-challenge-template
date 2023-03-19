package auth

import (
	"encoding/json"
	"net/http"

	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/application/service/auth"
	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
	responsehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/response"
	tokenhttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/token"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/auth"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/dbtrx"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template/pkg/security"
)

type Handler struct {
	AuthService authservice.IService
}

// New is the factory function that encapsulate the implementation related to auth handler.
func New(authService authservice.IService) IHandler {
	return &Handler{
		AuthService: authService,
	}
}

// SignUp godoc
// @tags authentication
// @summary API endpoint to perform sign up.
// @description
// @id SignUp
// @accept json
// @produce json
// @param credentials body securitypkg.Credentials true "SignUp"
// @success 200 {object} token.Token
// @failure 400 {object} error.Error
// @failure 404 {object} error.Error
// @failure 409 {object} error.Error
// @failure 500 {object} error.Error
// @router /sign_up [POST]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	dbTrx, ok := dbtrxmiddlewarepkg.FromContext(r.Context())
	if !ok || dbTrx == nil {
		responsehttputilpkg.RespondErrorWithJSON(w, customerror.New("failed to get the db_trx value from the request context"))
		return
	}

	credentials := securitypkg.Credentials{}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, customerror.BadRequest.New(err.Error()))
		return
	}

	token, err := h.AuthService.WithDBTrx(dbTrx).Register(credentials)
	if err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, err)
		return
	}

	responsehttputilpkg.RespondWithJSON(w, http.StatusOK, tokenhttputilpkg.Token{Text: token})
}

// SignIn godoc
// @tags authentication
// @summary API endpoint to perform sign in.
// @description
// @id SignIn
// @accept json
// @produce json
// @param credentials body securitypkg.Credentials true "SignIn"
// @success 200 {object} token.Token
// @failure 400 {object} error.Error
// @failure 401 {object} error.Error
// @failure 404 {object} error.Error
// @failure 409 {object} error.Error
// @failure 500 {object} error.Error
// @router /sign_in [POST]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	dbTrx, ok := dbtrxmiddlewarepkg.FromContext(r.Context())
	if !ok || dbTrx == nil {
		responsehttputilpkg.RespondErrorWithJSON(w, customerror.New("failed to get the db_trx value from the request context"))
		return
	}

	credentials := securitypkg.Credentials{}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, customerror.BadRequest.New(err.Error()))
		return
	}

	token, err := h.AuthService.WithDBTrx(dbTrx).LogIn(credentials)
	if err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, err)
		return
	}

	responsehttputilpkg.RespondWithJSON(w, http.StatusOK, tokenhttputilpkg.Token{Text: token})
}

// RefreshToken godoc
// @tags authentication
// @summary API endpoint to refresh the access token.
// @description
// @id RefreshToken
// @produce json
// @success 200 {object} token.Token
// @failure 400 {object} error.Error
// @failure 401 {object} error.Error
// @failure 500 {object} error.Error
// @router /refresh_token [POST]
// @security ApiKeyAuth
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	auth, ok := authmiddlewarepkg.FromContext(r.Context())
	if !ok || auth.IsEmpty() {
		responsehttputilpkg.RespondErrorWithJSON(w, customerror.New("failed to get the auth_details value from the request context"))
		return
	}

	token, err := h.AuthService.WithDBTrx(nil).RenewToken(auth)
	if err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, err)
	}

	responsehttputilpkg.RespondWithJSON(w, http.StatusOK, tokenhttputilpkg.Token{Text: token})
}

// ChangePassword godoc
// @tags authentication
// @summary API endpoint to reset the user's password.
// @description
// @id ChangePassword
// @accept json
// @produce json
// @param passwords body securitypkg.Passwords true "Reset Password"
// @success 200 {object} message.Message
// @failure 400 {object} error.Error
// @failure 401 {object} error.Error
// @failure 404 {object} error.Error
// @failure 500 {object} error.Error
// @router /change_password [POST]
// @security ApiKeyAuth
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	auth, ok := authmiddlewarepkg.FromContext(r.Context())
	if !ok || auth.IsEmpty() {
		responsehttputilpkg.RespondErrorWithJSON(w, customerror.New("failed to get the auth_details value from the request context"))
		return
	}

	passwords := securitypkg.Passwords{}

	if err := json.NewDecoder(r.Body).Decode(&passwords); err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, customerror.BadRequest.New(err.Error()))
		return
	}

	err := h.AuthService.WithDBTrx(nil).ModifyPassword(auth.UserID.String(), passwords)
	if err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, err)
		return
	}

	responsehttputilpkg.RespondWithJSON(w, http.StatusOK, responsehttputilpkg.Message{Text: "the password has been updated successfully"})
}

// SignOut godoc
// @tags authentication
// @summary API endpoint to perform sign out.
// @description
// @id SignOut
// @produce json
// @success 200 {object} message.Message
// @failure 400 {object} error.Error
// @failure 401 {object} error.Error
// @failure 404 {object} error.Error
// @failure 500 {object} error.Error
// @router /sign_out [POST]
// @security ApiKeyAuth
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	auth, ok := authmiddlewarepkg.FromContext(r.Context())
	if !ok || auth.IsEmpty() {
		responsehttputilpkg.RespondErrorWithJSON(w, customerror.New("failed to get the auth_details value from the request context"))
		return
	}

	err := h.AuthService.WithDBTrx(nil).LogOut(auth.ID.String())
	if err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, err)
		return
	}

	responsehttputilpkg.RespondWithJSON(w, http.StatusOK, responsehttputilpkg.Message{Text: "you have logged out successfully"})
}
