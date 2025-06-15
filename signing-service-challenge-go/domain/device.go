package domain

import (
	"context"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type UserService struct {
	repo any //UserRepository
}

func (s *UserService) CreateSignatureDevice(ctx context.Context, algorithm, label string) (Device, error) {
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
		ID:                id,
		Algorithm:         algorithm,
		Label:             label,
		PublicKey:         publicKey,
		PrivateKey:        privateKey,
		Signature_counter: 0,
	}
	return device, nil
}

func (s *UserService) SignTransaction(ctx context.Context, ID uuid.UUID, data string) (Signature, error) {
	// Retrieve the device from the persistence layer using the ID
	device := Device{}

	device.Signature_counter++
	return Signature{}, nil
}
