package api

import (
	"encoding/json"
	"net/http"
)

// CreateSignatureDevice godoc
// @Title CreateSignatureDevice
// @Summary Create a new signature device
// @Description Creates a new signature device with the specified parameters
// @Tags SignatureDevices
// @Accept json
// @Produce json
// @Param id query string true "Device ID"
// @Param algorithm query string true "Algorithm (ECC or RSA)"
// @Param label query string false "Optional label for the device"
// @Success 200 {object} CreateSignatureDeviceResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /signature-device [post]
func (s *Server) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	algorithm := r.URL.Query().Get("algorithm")
	label := r.URL.Query().Get("label")

	// Validate required parameters
	if id == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: id"})
		return
	}
	if algorithm == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: algorithm"})
		return
	}

	// Validate algorithm value
	if algorithm != "ECC" && algorithm != "RSA" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid algorithm. Must be 'ECC' or 'RSA'"})
		return
	}
	// Prepare response data
	data := map[string]any{
		"id":        id,
		"algorithm": algorithm,
		"label":     label,
	}

	// ctx := r.Context()

	// Set response headers and encode the data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to encode response"})
	}
}

// SignTransaction godoc
// @Title SignTransaction
// @Summary Sign a transaction
// @Description Signs a transaction using the specified device ID and data payload.
// @Tags Device
// @Accept json
// @Produce json
// @Param deviceId path string true "Device ID"
// @Param data body string true "Data to be signed"
// @Success 200 {object} SignatureResponse "Signature successfully generated"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 404 {object} ErrorResponse "Device not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /sign/{deviceId} [post]
func (s *Server) SignTransaction(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("deviceId")
	data := r.URL.Query().Get("data")

	// Validate required parameters
	if deviceId == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: deviceId"})
		return
	}
	if data == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: data"})
		return
	}

	// Simulate signing logic (replace with actual implementation)
	signature := "signed_" + data

	// Prepare response data
	response := map[string]any{
		"deviceId":  deviceId,
		"signature": signature,
	}

	// Set response headers and encode the data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to encode response"})
	}
}

func (s *Server) GetSignatureDevice(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) ListSignatureDevices(w http.ResponseWriter, r *http.Request) {
}
