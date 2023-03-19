package healthcheck

import (
	"net/http"

	healthcheckservice "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/application/service/healthcheck"
	responsehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/response"
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
// @summary API endpoint used to verify if the service has started up correctly and is ready to accept requests.
// @description
// @id GetStatus
// @produce json
// @success 200 {object} message.Message
// @failure 500 {object} error.Error
// @router /status [GET]
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if err := h.HealthCheckService.GetStatus(); err != nil {
		responsehttputilpkg.RespondErrorWithJSON(w, err)
	}

	text := "everything is up and running"
	responsehttputilpkg.RespondWithJSON(w, http.StatusOK, responsehttputilpkg.Message{Text: text})
}
