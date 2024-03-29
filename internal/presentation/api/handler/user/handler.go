package user

import (
	"net/http"

	userservice "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/application/service/user"
	presentableentity "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/presentity"
	responsehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/response"
)

type Handler struct {
	UserService userservice.IService
}

// New is the factory function that encapsulate the implementation related to user handler.
func New(userService userservice.IService) IHandler {
	return &Handler{
		UserService: userService,
	}
}

// GetAllUsers godoc
// @tags user
// @summary API endpoint to get the list of all users.
// @description
// @id GetAllUsers
// @produce json
// @success 200 {array} model.User
// @failure 401 {object} error.Error
// @failure 500 {object} error.Error
// @router /users [GET]
// @security ApiKeyAuth
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	domainUsers, err := h.UserService.WithDBTrx(nil).GetAll()
	if err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, err)
		return
	}

	users := presentableentity.Users{}
	users.FromDomain(domainUsers)

	responsehttputilpkg.RespondWithJSON(w, http.StatusOK, users)
}
