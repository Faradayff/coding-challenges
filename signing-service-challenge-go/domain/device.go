package domain

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
)

type UserService struct {
	repo *persistence.DeviceRepository
}

// NewUserService crea una nueva instancia de UserService con el repositorio inyectado
func NewUserService(repo *persistence.DeviceRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateSignatureDevice(ctx context.Context, algorithm, label string) (model.Device, error) {
	id := uuid.New()

	// Creating new public and private keys
	var publicKey, privateKey any
	if algorithm == "ECC" {
		eccGenerator := crypto.ECCGenerator{}

		eccKeys, err := eccGenerator.Generate()
		if err != nil {
			return model.Device{}, err
		}

		publicKey = eccKeys.Public
		privateKey = eccKeys.Private
	} else if algorithm == "RSA" {
		rsaGenerator := crypto.RSAGenerator{}

		rsaKeys, err := rsaGenerator.Generate()
		if err != nil {
			return model.Device{}, err
		}

		publicKey = rsaKeys.Public
		privateKey = rsaKeys.Private
	}

	// Create the new device
	device := model.Device{
		ID:               id,
		Algorithm:        algorithm,
		Label:            label,
		PublicKey:        publicKey,
		PrivateKey:       privateKey,
		SignatureCounter: 0,
	}

	// Save it in memory
	err := s.repo.Create(device)
	if err != nil {
		return model.Device{}, fmt.Errorf("failed to save device: %w", err)
	}

	return device, nil
}

func (s *UserService) SignTransaction(ctx context.Context, ID uuid.UUID, data string) (model.SignaturedData, error) {
	// Retrieve the device from the persistence layer using the ID
	device, err := s.repo.FindByID(ID)
	if err != nil {
		return model.SignaturedData{}, fmt.Errorf("device not found: %w", err)
	}

	// Preparing data to be signed
	header := device.SignatureCounter
	body := data
	var end string
	if device.SignatureCounter == 0 {
		idBytes, err := device.ID.MarshalBinary()
		if err != nil {
			return model.SignaturedData{}, err
		}
		end = base64.StdEncoding.EncodeToString(idBytes)
	} else {
		end = device.LastSignature
	}
	preparedData := fmt.Sprintf("%d_%s_%s", header, body, end)

	// Signing the data
	var signature []byte

	if device.Algorithm == "ECC" {
		eccSigner := crypto.ECCSigner{}
		signature, err = eccSigner.Sign(preparedData, device.PrivateKey, device.PublicKey)
		if err != nil {
			return model.SignaturedData{}, fmt.Errorf("failed to sign data: %w", err)
		}
	} else if device.Algorithm == "RSA" {
		rsaSigner := crypto.RSASigner{}
		signature, err = rsaSigner.Sign(preparedData, device.PrivateKey, device.PublicKey)
		if err != nil {
			return model.SignaturedData{}, fmt.Errorf("failed to sign data: %w", err)
		}
	}

	// Creating returning data
	signaturedData := model.SignaturedData{
		Signature:  signature,
		SignedData: preparedData,
	}

	// Increasing the counter
	device.SignatureCounter++

	// Setting the new last signature
	device.LastSignature = base64.StdEncoding.EncodeToString(signature)

	// No need to save the changes in the device since we have a pointer to it

	return signaturedData, nil
}
