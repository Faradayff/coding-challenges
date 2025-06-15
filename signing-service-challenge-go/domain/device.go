package domain

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type UserService struct {
	repo any //UserRepository
}

func (s *UserService) CreateSignatureDevice(ctx context.Context, algorithm, label string) (Device, error) {
	id := uuid.New()

	// Creating new public and private keys
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
		ID:               id,
		Algorithm:        algorithm,
		Label:            label,
		PublicKey:        publicKey,
		PrivateKey:       privateKey,
		SignatureCounter: 0,
	}
	return device, nil
}

func (s *UserService) SignTransaction(ctx context.Context, ID uuid.UUID, data string) (SignaturedData, error) {
	// Retrieve the device from the persistence layer using the ID
	device := Device{}

	// Preparing data to be signed
	header := device.SignatureCounter
	body := data
	var end string
	if device.SignatureCounter == 0 {
		idBytes, err := device.ID.MarshalBinary()
		if err != nil {
			return SignaturedData{}, err
		}
		end = base64.StdEncoding.EncodeToString(idBytes)
	} else {
		end = device.LastSignature
	}
	preparedData := fmt.Sprintf("%d_%s_%s", header, body, end)

	// Signing the data
	var signature []byte
	var err error

	if device.Algorithm == "ECC" {
		eccSigner := crypto.ECCSigner{}
		signature, err = eccSigner.Sign(preparedData, device.PrivateKey, device.PublicKey)
		if err != nil {
			return SignaturedData{}, fmt.Errorf("failed to sign data: %w", err)
		}
	} else if device.Algorithm == "RSA" {
		rsaSigner := crypto.RSASigner{}
		signature, err = rsaSigner.Sign(preparedData, device.PrivateKey, device.PublicKey)
		if err != nil {
			return SignaturedData{}, fmt.Errorf("failed to sign data: %w", err)
		}
	}

	// Creating returning data
	signaturedData := SignaturedData{
		Signature:  signature,
		SignedData: preparedData,
	}

	// Increasing the counter
	device.SignatureCounter++

	// Setting the new last signature
	device.LastSignature = base64.StdEncoding.EncodeToString(signature)

	return signaturedData, nil
}
