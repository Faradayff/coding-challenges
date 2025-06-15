package domain

import "github.com/google/uuid"

type Device struct {
	ID                uuid.UUID `json:"id"`
	Algorithm         string    `json:"algorithm"`
	Label             string    `json:"label"`
	PublicKey         any       `json:"publicKey"`
	PrivateKey        any       `json:"privateKey"`
	Signature_counter int       `json:"signatureCounter"`
}

type Signature struct {
	ID                uuid.UUID `json:"id"`
	Signature_counter int       `json:"signatureCounter"`
}
