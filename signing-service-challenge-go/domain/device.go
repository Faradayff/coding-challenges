package domain

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
)

type DeviceService struct {
	repo       *persistence.DeviceRepository
	devicesMus map[uuid.UUID]*sync.Mutex // map to avoid signning from the same device at the same time
	mu         sync.Mutex                // mutex to avoid concurrent access to the mutexes map
}

// NewUserService creates a new UserService instance with the provided repository and initializes the mutex map
func NewUserService(repo *persistence.DeviceRepository) *DeviceService {
	return &DeviceService{
		repo:       repo,
		devicesMus: make(map[uuid.UUID]*sync.Mutex),
	}
}

// CreateSignatureDevice creates a new signature device with the specified algorithm and label
func (s *DeviceService) CreateSignatureDevice(ctx context.Context, algorithm, label string) (model.Device, error) {
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

	// Save it in the database
	err := s.repo.Create(device)
	if err != nil {
		return model.Device{}, fmt.Errorf("failed to save device: %w", err)
	}

	// Create the mutex for this device
	s.mu.Lock()
	defer s.mu.Unlock()

	s.devicesMus[id] = &sync.Mutex{}

	return device, nil
}

// SignTransaction signs the provided data using the device's private key and returns the signed data
func (s *DeviceService) SignTransaction(ctx context.Context, id uuid.UUID, data string) (model.SignaturedData, error) {
	// Retrieve the device from the persistence layer using the ID
	device, err := s.repo.FindByID(id)
	if err != nil {
		return model.SignaturedData{}, fmt.Errorf("device not found: %w", err)
	}

	// Blocking device from be modified or accessed until we have finished
	s.devicesMus[id].Lock()
	defer s.devicesMus[id].Unlock()

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

	// Updating signature counter and last signature of the device
	err = s.repo.AfterSignUpdateDevice(device.ID, base64.StdEncoding.EncodeToString(signature))
	if err != nil {
		return model.SignaturedData{}, fmt.Errorf("failed to update device after signing: %w", err)
	}

	return signaturedData, nil
}

// GetDevice retrieves a device by its ID
func (s *DeviceService) GetDevice(ctx context.Context, ID uuid.UUID) (model.Device, error) {
	device, err := s.repo.FindByID(ID)
	if err != nil {
		return model.Device{}, fmt.Errorf("device not found: %w", err)
	}

	return *device, nil
}

// GetAllDevices retrieves all devices
func (s *DeviceService) GetAllDevices(ctx context.Context) ([]model.Device, error) {
	devices, err := s.repo.GetAll()
	if err != nil {
		return []model.Device{}, fmt.Errorf("error retrieving all the devices: %w", err)
	}

	return devices, nil
}
