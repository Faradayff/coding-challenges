package domain

import (
	"context"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

func CreateSignatureDevice(ctx context.Context, algorithm, label string) (Device, error) {
	id := uuid.New()

	var publicKey, privateKey any
	if algorithm == "ECC" {
		eccGenerator := crypto.ECCGenerator{}

		eccKeys, err := eccGenerator.Generate()
		if err != nil {
			return Device{}, err
		}

		publicKey = eccKeys.Public
		privateKey = eccKeys.Private
	} else if algorithm == "RSA" {
		rsaGenerator := crypto.RSAGenerator{}

		rsaKeys, err := rsaGenerator.Generate()
		if err != nil {
			return Device{}, err
		}

		publicKey = rsaKeys.Public
		privateKey = rsaKeys.Private
	}

	// Call the persistence layer to save the device

	device := Device{
		ID:         id,
		Algorithm:  algorithm,
		Label:      label,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	return device, nil
}

func SignTransaction(ctx context.Context, ID uuid.UUID, data string) (Signature, error) {
	return Signature{}, nil
}
