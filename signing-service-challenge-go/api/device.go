package api

import (
	"net/http"

	"encoding/json"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/utils"
	"github.com/google/uuid"
)

type DeviceApi struct {
	service *domain.UserService
}

func NewDeviceApi(service *domain.UserService) *DeviceApi {
	return &DeviceApi{
		service: service,
	}
}

// CreateSignatureDevice godoc
// @Title CreateSignatureDevice
// @Summary Create a new signature device
// @Description Creates a new signature device with the specified parameters
// @Tags Devices
// @Produce json
// @Param algorithm query string true "Algorithm (ECC or RSA)"
// @Param label query string true "Label for the device"
// @Success 200 {object} CreateDeviceResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /new-device [post]
func (a *DeviceApi) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
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

	// Calling the service
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

	// Creating response
	createSignatureDeviceResponse := CreateDeviceResponse{
		ID:         device.ID,
		Algorithm:  device.Algorithm,
		Label:      device.Label,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	WriteAPIResponse(w, http.StatusCreated, createSignatureDeviceResponse)
}

// SignTransaction godoc
// @Title SignTransaction
// @Summary Sign a transaction
// @Description Signs a transaction using the specified device ID and data payload.
// @Tags Devices
// @Produce json
// @Param deviceId path string true "Device ID"
// @Param data body string true "Data to be signed"
// @Success 200 {object} SignaturedDataResponse "Signature successfully generated"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 404 {object} ErrorResponse "Device not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /sign/{deviceId} [post]
func (a *DeviceApi) SignTransaction(w http.ResponseWriter, r *http.Request) {
	// Get and validate deviceId
	deviceId := r.URL.Query().Get("deviceId")

	if deviceId == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: deviceId"})
		return
	}
	uuid, err := uuid.Parse(deviceId)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid deviceId. Must be a valid UUID"})
		return
	}

	// Get and validate data
	var req SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid request body"})
		return
	}

	// Validate data recieved
	if req.Data == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Field 'data' is required"})
		return
	}

	ctx := r.Context()

	// Calling the service
	signaturedData, err := a.service.SignTransaction(ctx, uuid, req.Data)
	if err != nil {
		if err.Error() == "device not found" {
			WriteErrorResponse(w, http.StatusNotFound, []string{"Device not found"})
			return
		} else {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to sign transaction", err.Error()})
			return
		}
	}

	// Creating response
	signaturedDataResponse := SignaturedDataResponse{
		Signature:  signaturedData.Signature,
		SignedData: signaturedData.SignedData,
	}

	WriteAPIResponse(w, http.StatusOK, signaturedDataResponse)
}

// GetDevice godoc
// @Title GetDevice
// @Summary Get a device
// @Description Retrieves a device by its ID and returns its details.
// @Tags Devices
// @Produce json
// @Param deviceId path string true "Device ID"
// @Success 200 {object} GetDeviceResponse "Device successfully retrieved"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 404 {object} ErrorResponse "Device not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /{deviceId} [get]
func (a *DeviceApi) GetDevice(w http.ResponseWriter, r *http.Request) {
	// Get and validate deviceId
	deviceId := r.URL.Query().Get("deviceId")

	if deviceId == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required parameter: deviceId"})
		return
	}
	uuid, err := uuid.Parse(deviceId)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid id. Must be a valid UUID"})
		return
	}

	ctx := r.Context()

	// Calling the service
	device, err := a.service.GetDevice(ctx, uuid)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	// Creating response
	getDeviceResponse, err := deviceToGetDeviceResponse(device)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to convert device to response", err.Error()})
	}

	WriteAPIResponse(w, http.StatusOK, getDeviceResponse)
}

// GetDevice godoc
// @Title GetAllDevices
// @Summary Get all the devices
// @Description Retrieves all the devices and its details.
// @Tags Devices
// @Produce json
// @Success 200 {object} GetAllDevicesResponse "Devices successfully retrieved"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /all [get]
func (a *DeviceApi) GetAllDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Calling the service
	devices, err := a.service.GetAllDevices(ctx)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	// Creating response
	getDeviceResponse, err := devicesToGetAllDevicesResponse(devices)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to convert device to response", err.Error()})
	}

	WriteAPIResponse(w, http.StatusOK, getDeviceResponse)
}

// Convert Device to GetDeviceResponse
func deviceToGetDeviceResponse(device model.Device) (GetDeviceResponse, error) {
	var publicKey, privateKey string
	var err error

	if device.Algorithm == "ECC" {
		publicKey, err = utils.ECCPublicKeyToString(device.PublicKey)
		if err != nil {
			return GetDeviceResponse{}, err
		}

		privateKey, err = utils.ECCPrivateKeyToString(device.PrivateKey)
		if err != nil {
			return GetDeviceResponse{}, err
		}
	} else if device.Algorithm == "RSA" {
		publicKey, err = utils.RSAPublicKeyToString(device.PublicKey)
		if err != nil {
			return GetDeviceResponse{}, err
		}

		privateKey, err = utils.RSAPrivateKeyToString(device.PrivateKey)
		if err != nil {
			return GetDeviceResponse{}, err
		}
	}

	return GetDeviceResponse{
		ID:               device.ID,
		Algorithm:        device.Algorithm,
		Label:            device.Label,
		PublicKey:        publicKey,
		PrivateKey:       privateKey,
		SignatureCounter: device.SignatureCounter,
		LastSignature:    device.LastSignature,
	}, nil
}

// Convert a slice of Devices to a GetAllDevicesResponse
func devicesToGetAllDevicesResponse(devices []model.Device) (GetAllDevicesResponse, error) {
	var deviceResponses []GetDeviceResponse
	for _, device := range devices {
		deviceResponse, err := deviceToGetDeviceResponse(device)
		if err != nil {
			return GetAllDevicesResponse{}, err
		}
		deviceResponses = append(deviceResponses, deviceResponse)
	}

	return GetAllDevicesResponse{
		Devices: deviceResponses,
		Total:   len(deviceResponses),
	}, nil
}
