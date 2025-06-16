package api

import "github.com/google/uuid"

type CreateDeviceResponse struct {
	ID         uuid.UUID `json:"id"`
	Algorithm  string    `json:"algorithm"`
	Label      string    `json:"label"`
	PublicKey  string    `json:"publicKey"`
	PrivateKey string    `json:"privateKey"`
}

type SignaturedDataResponse struct {
	Signature  []byte `json:"signature"`
	SignedData string `json:"signed_data"`
}

type GetDeviceResponse struct {
	ID               uuid.UUID `json:"id"`
	Algorithm        string    `json:"algorithm"`
	Label            string    `json:"label"`
	PublicKey        string    `json:"publicKey"`
	PrivateKey       string    `json:"privateKey"`
	SignatureCounter int       `json:"signatureCounter"`
	LastSignature    string    `json:"lastSignature,omitempty"`
}

// GetAllDevicesResponse created for possible future use of pagination
type GetAllDevicesResponse struct {
	Devices []GetDeviceResponse `json:"devices"`
	Total   int                 `json:"total"`
}

type SignTransactionRequest struct {
	Data string `json:"data"`
}
