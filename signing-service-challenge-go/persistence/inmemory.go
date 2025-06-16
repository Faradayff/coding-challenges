package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/google/uuid"
)

type DeviceRepoInterface interface {
	Create(device model.Device) error
	FindByID(id uuid.UUID) (*model.Device, error)
	GetAll() ([]model.Device, error)
	AfterSignUpdateDevice(id uuid.UUID, lastSignature string) error
}

type DeviceRepository struct {
	data map[uuid.UUID]model.Device
	mu   sync.RWMutex
}

// Initialize
func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{
		data: make(map[uuid.UUID]model.Device),
	}
}

// Create stores a new device
func (r *DeviceRepository) Create(device model.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[device.ID] = device
	return nil
}

// FindByID retrieves a device by its ID
func (r *DeviceRepository) FindByID(id uuid.UUID) (*model.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.data[id]
	if !exists {
		return nil, errors.New("device not found")
	}

	return &device, nil
}

// GetAll retrieves all the devices
func (r *DeviceRepository) GetAll() ([]model.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]model.Device, 0, len(r.data))
	for _, device := range r.data {
		devices = append(devices, device)
	}

	return devices, nil
}

// AfterSignUpdateDevice increments the signature counter and update the last signature checking multiple accesses
func (r *DeviceRepository) AfterSignUpdateDevice(id uuid.UUID, lastSignature string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	device, exists := r.data[id]
	if !exists {
		return errors.New("device not found")
	}

	device.SignatureCounter++
	device.LastSignature = lastSignature

	r.data[id] = device
	return nil
}
