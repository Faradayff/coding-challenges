package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/google/uuid"
)

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

// Delete removes a device by its ID
func (r *DeviceRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return errors.New("device not found")
	}

	delete(r.data, id)
	return nil
}
