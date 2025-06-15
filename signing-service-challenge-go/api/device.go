package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/utils"
	"github.com/google/uuid"
)

type Api struct {
	service domain.UserService
}

// CreateSignatureDevice godoc
// @Title CreateSignatureDevice
// @Summary Create a new signature device
// @Description Creates a new signature device with the specified parameters
// @Tags Devices
// @Accept json
// @Produce json
// @Param algorithm query string true "Algorithm (ECC or RSA)"
// @Param label query string true "Label for the device"
// @Success 200 {object} CreateSignatureDeviceResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /signature-device [post]
func (a *Api) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
	algorithm := r.URL.Query().Get("algorithm")
	label := r.URL.Query().Get("label")

	// Validate required parameters
	if algorithm == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: algorithm"})
		return
	}

	// Validate algorithm value
	if algorithm != "ECC" && algorithm != "RSA" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid algorithm. Must be 'ECC' or 'RSA'"})
		return
	}

	ctx := r.Context()

	device, err := a.service.CreateSignatureDevice(ctx, algorithm, label)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to create signature device", err.Error()})
		return
	}

	var publicKey, privateKey string
	if algorithm == "ECC" {

		publicKey, err = utils.ECCPublicKeyToString(device.PublicKey)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}

		privateKey, err = utils.ECCPrivateKeyToString(device.PrivateKey)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}

	} else if algorithm == "RSA" {

		publicKey, err = utils.RSAPublicKeyToString(device.PublicKey)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}

		privateKey, err = utils.RSAPrivateKeyToString(device.PrivateKey)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}
	}

	createSignatureDeviceResponse := CreateSignatureDeviceResponse{
		ID:         device.ID,
		Algorithm:  device.Algorithm,
		Label:      device.Label,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	// Set response headers and encode the data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(createSignatureDeviceResponse); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to encode response"})
	}
}

// SignTransaction godoc
// @Title SignTransaction
// @Summary Sign a transaction
// @Description Signs a transaction using the specified device ID and data payload.
// @Tags Devices
// @Accept json
// @Produce json
// @Param deviceId path string true "Device ID"
// @Param data body string true "Data to be signed"
// @Success 200 {object} SignatureResponse "Signature successfully generated"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 404 {object} ErrorResponse "Device not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /sign/{deviceId} [post]
func (a *Api) SignTransaction(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("deviceId")
	data := r.URL.Query().Get("data")

	// Validate required parameters
	if deviceId == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: deviceId"})
		return
	}
	uuid, err := uuid.Parse(deviceId)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid deviceId. Must be a valid UUID"})
		return
	}
	if data == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: data"})
		return
	}

	ctx := r.Context()

	signature, err := a.service.SignTransaction(ctx, uuid, data)
	if err != nil {
		if err.Error() == "device not found" {
			WriteErrorResponse(w, http.StatusNotFound, []string{"Device not found"})
			return
		} else {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to sign transaction", err.Error()})
			return
		}
	}

	response := SignatureResponse{
		ID: signature.ID,
	}

	// Set response headers and encode the data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to encode response"})
	}
}

func (a *Api) GetSignatureDevice(w http.ResponseWriter, r *http.Request) {
}

func (a *Api) ListSignatureDevices(w http.ResponseWriter, r *http.Request) {
}
