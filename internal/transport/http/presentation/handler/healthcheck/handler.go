package healthcheck

import (
	"net/http"

	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
)

type Handler struct{}

// New is the factory function that encapsulates the implementation related to healthcheck handler.
func New() IHandler {
	return &Handler{}
}

// GetStatus godoc
// @tags health check
// @Summary API endpoint used to verify if the service has started up correctly and is ready to accept requests.
// @Description
// @ID GetStatus
// @Produce json
// @Success 200 {object} message.Message
// @Router /status [GET]
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	text := "everything is up and running"

	responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: text})
}
