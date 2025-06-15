package domain

import "github.com/google/uuid"

type Device struct {
	ID               uuid.UUID `json:"id"`
	Algorithm        string    `json:"algorithm"`
	Label            string    `json:"label"`
	PublicKey        any       `json:"publicKey"`
	PrivateKey       any       `json:"privateKey"`
	SignatureCounter int       `json:"signatureCounter"`
	LastSignature    string    `json:"lastSignature,omitempty"`
}

type SignaturedData struct {
	Signature  []byte `json:"signature"`
	SignedData string `json:"signed_data"`
}
