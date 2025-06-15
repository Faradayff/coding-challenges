package api

import "github.com/google/uuid"

type CreateSignatureDeviceResponse struct {
	ID         uuid.UUID `json:"id"`
	Algorithm  string    `json:"algorithm"`
	Label      string    `json:"label"`
	PublicKey  string    `json:"publicKey"`
	PrivateKey string    `json:"privateKey"`
}

type SignatureResponse struct {
	ID uuid.UUID `json:"deviceId"`
}

type GetSignatureDeviceResponse struct {
	ID                uuid.UUID `json:"id"`
	Algorithm         string    `json:"algorithm"`
	Label             string    `json:"label"`
	PublicKey         string    `json:"publicKey"`
	PrivateKey        string    `json:"privateKey"`
	Signature_counter int       `json:"signatureCounter"`
}
