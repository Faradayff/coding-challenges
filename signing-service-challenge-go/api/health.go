package api

import "net/http"

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Health godoc
// @Title Health Check
// @Summary Check the health of the service
// @Description Evaluates the health of the service and returns a standardized response.
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "Service is healthy"
// @Failure 405 {object} ErrorResponse "Method not allowed"
// @Router /api/v0/health

// Health evaluates the health of the service and writes a standardized response.
func (s *Server) Health(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	health := HealthResponse{
		Status:  "pass",
		Version: "v0",
	}

	WriteAPIResponse(response, http.StatusOK, health)
}
