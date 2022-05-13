package user

import (
	"net/http"

	userservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/service/user"
	httppresentationmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/model"
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
// @param param1 query int false "Param 1"
// @param param2 query int true "Param 2"
// @success 200 {array} user.User
// @failure 500 {object} error.Error
// @router /users [GET]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	domainUsers, err := h.UserService.WithDBTrx(nil).GetAll()
	if err != nil {
		responsehttputilpkg.RespondErrorWithJson(w, err)
		return
	}

	users := httppresentationmodel.Users{}
	users.FromDomain(domainUsers)

	responsehttputilpkg.RespondWithJson(w, http.StatusOK, users)
}
