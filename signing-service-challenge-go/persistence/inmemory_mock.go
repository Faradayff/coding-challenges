package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/google/uuid"
)

type MockDeviceRepo struct {
	CreateFunc                func(device model.Device) error
	FindByIDFunc              func(id uuid.UUID) (*model.Device, error)
	GetAllFunc                func() ([]model.Device, error)
	AfterSignUpdateDeviceFunc func(id uuid.UUID, lastSignature string) error
}

func (m *MockDeviceRepo) Create(device model.Device) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(device)
	}
	return nil
}

func (m *MockDeviceRepo) FindByID(id uuid.UUID) (*model.Device, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return &model.Device{
		ID:               id,
		Algorithm:        "ECC",
		Label:            "Mock Device",
		PublicKey:        nil,
		PrivateKey:       nil,
		SignatureCounter: 0,
		LastSignature:    "",
	}, nil
}

func (m *MockDeviceRepo) GetAll() ([]model.Device, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockDeviceRepo) AfterSignUpdateDevice(id uuid.UUID, lastSignature string) error {
	if m.AfterSignUpdateDeviceFunc != nil {
		return m.AfterSignUpdateDeviceFunc(id, lastSignature)
	}
	return nil
}
