package domain

import (
	"context"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/google/uuid"
)

// MockDeviceService is a mock implementation of DeviceServiceInterface for testing purposes
type MockDeviceService struct {
	CreateSignatureDeviceFunc func(ctx context.Context, algorithm, label string) (model.Device, error)
	SignTransactionFunc       func(ctx context.Context, id uuid.UUID, data string) (model.SignaturedData, error)
	GetDeviceFunc             func(ctx context.Context, id uuid.UUID) (model.Device, error)
	GetAllDevicesFunc         func(ctx context.Context) ([]model.Device, error)
}

func (m *MockDeviceService) CreateSignatureDevice(ctx context.Context, algorithm, label string) (model.Device, error) {
	return m.CreateSignatureDeviceFunc(ctx, algorithm, label)
}

func (m *MockDeviceService) SignTransaction(ctx context.Context, id uuid.UUID, data string) (model.SignaturedData, error) {
	return m.SignTransactionFunc(ctx, id, data)
}

func (m *MockDeviceService) GetDevice(ctx context.Context, id uuid.UUID) (model.Device, error) {
	return m.GetDeviceFunc(ctx, id)
}

func (m *MockDeviceService) GetAllDevices(ctx context.Context) ([]model.Device, error) {
	return m.GetAllDevicesFunc(ctx)
}
