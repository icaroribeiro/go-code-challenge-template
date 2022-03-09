package healthcheck

import (
	"net/http"

	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/service/healthcheck"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
)

type Handler struct {
	HealthCheckService healthcheckservice.IService
}

// New is the factory function that encapsulates the implementation related to healthcheck handler.
func New(healthCheckService healthcheckservice.IService) IHandler {
	return &Handler{
		HealthCheckService: healthCheckService,
	}
}

// GetStatus godoc
// @tags health check
// @Summary API endpoint used to verify if the service has started up correctly and is ready to accept requests.
// @Description
// @ID GetStatus
// @Produce json
// @Success 200 {object} message.Message
// @Failure 500 {object} error.Error
// @Router /status [GET]
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	text := "everything is up and running"

	if err := h.HealthCheckService.GetStatus(); err != nil {
		responsehttputilpkg.RespondErrorWithJson(w, err)
	}

	responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: text})
}
