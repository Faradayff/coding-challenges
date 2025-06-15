package api

import "github.com/google/uuid"

type CreateSignatureDeviceResponse struct {
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

type GetSignatureDeviceResponse struct {
	ID               uuid.UUID `json:"id"`
	Algorithm        string    `json:"algorithm"`
	Label            string    `json:"label"`
	PublicKey        string    `json:"publicKey"`
	PrivateKey       string    `json:"privateKey"`
	SignatureCounter int       `json:"signatureCounter"`
	LastSignature    string    `json:"lastSignature,omitempty"`
}
