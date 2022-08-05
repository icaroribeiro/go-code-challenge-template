package user

import (
	"net/http"

	userservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/service/user"
	presentationentity "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/entity"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
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
		responsehttputilpkg.RespondErrorWithJson(w, err)
		return
	}

	users := presentationentity.Users{}
	users.FromDomain(domainUsers)

	responsehttputilpkg.RespondWithJson(w, http.StatusOK, users)
}
