package persistence

import (
	"crypto/rsa"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeviceRepo", func() {
	var (
		deviceRepo *DeviceRepository
	)

	BeforeEach(func() {
		deviceRepo = NewDeviceRepository()
	})

	Describe("Create", func() {
		Context("when creating a new device", func() {
			It("should create a device with ECC algorithm", func() {
				device := model.Device{
					ID:               uuid.New(),
					Algorithm:        "ECC",
					Label:            "Test device",
					PublicKey:        nil,
					PrivateKey:       nil,
					SignatureCounter: 0,
				}

				err := deviceRepo.Create(device)
				Expect(err).To(BeNil(), "Failed to create device")
			})
		})
	})

	Describe("FindByID", func() {
		var deviceID uuid.UUID
		var algorithm, label string

		BeforeEach(func() {
			// Create a device before each context
			algorithm = "RSA"
			label = "Test Device"
			device := model.Device{
				ID:               uuid.New(),
				Algorithm:        algorithm,
				Label:            label,
				PublicKey:        &rsa.PublicKey{},
				PrivateKey:       &rsa.PrivateKey{},
				SignatureCounter: 0,
			}
			err := deviceRepo.Create(device)
			if err != nil {
				Fail(fmt.Sprintf("Failed setting up the device: %v", err))
			}
			deviceID = device.ID
		})

		Context("when getting a device", func() {
			It("should return the device with same id", func() {
				createdDevice, err := deviceRepo.FindByID(deviceID)
				Expect(err).To(BeNil(), "Failed to find created device")
				Expect(createdDevice).ToNot(BeNil(), "Created device should not be nil")
				Expect(createdDevice.ID).To(Equal(deviceID), "Device ID should match the created device ID")
				Expect(createdDevice.Algorithm).To(Equal(algorithm), "Device algorithm should match the created device algorithm")
				Expect(createdDevice.Label).To(Equal(label), "Device label should match the created device label")
				Expect(createdDevice.PublicKey).ToNot(BeNil(), "Public key should not be nil")
				Expect(createdDevice.PrivateKey).ToNot(BeNil(), "Private key should not be nil")
				Expect(createdDevice.SignatureCounter).To(Equal(0), "Signature counter should be initialized to 0")
				Expect(createdDevice.LastSignature).To(BeEmpty(), "Last signature should be empty")
			})
		})

		Context("when getting a non existent device", func() {
			It("should return the device with same id", func() {
				_, err := deviceRepo.FindByID(uuid.New())
				Expect(err).To(HaveOccurred(), "Expected an error when finding a non-existent device")
				Expect(err.Error()).To(ContainSubstring("device not found"), "Error message should indicate that the device was not found")
			})
		})
	})

	Describe("GetAll", func() {
		var devicesIDs []uuid.UUID

		BeforeEach(func() {
			// Create a few devices before each context
			algorithm := "RSA"
			devicesIDs = make([]uuid.UUID, 0)
			for i := range 3 {
				id := uuid.New()
				device := model.Device{
					ID:               id,
					Algorithm:        algorithm,
					Label:            fmt.Sprintf("Test Device %d", i),
					PublicKey:        &rsa.PublicKey{},
					PrivateKey:       &rsa.PrivateKey{},
					SignatureCounter: 0,
				}
				err := deviceRepo.Create(device)
				if err != nil {
					Fail(fmt.Sprintf("Failed setting up the device: %v", err))
				}
				devicesIDs = append(devicesIDs, id)
			}
		})

		Context("when getting all devices", func() {
			It("should return a slice with all the devices", func() {
				devices, err := deviceRepo.GetAll()
				Expect(err).To(BeNil(), "Failed to get all devices")
				Expect(devices).ToNot(BeEmpty(), "Devices slice should not be empty")
				Expect(devices).To(BeAssignableToTypeOf([]model.Device{}), "Devices should be of type []model.Device")
				Expect(devices).To(HaveLen(len(devicesIDs)), "Devices slice should contain same number of devices than devicesIDs")

				for _, device := range devices {
					Expect(device).To(BeAssignableToTypeOf(model.Device{}), "Each device should be of type model.Device")
					Expect(device.ID).To(BeElementOf(devicesIDs), "Device ID should be one of the created devices")
					Expect(device.Algorithm).To(Equal("RSA"), "Algorithm should be RSA")
					Expect(device.Label).To(ContainSubstring("Test Device"), "Label should contain 'Test Device'")
					Expect(device.PublicKey).ToNot(BeNil(), "Public key should not be nil")
					Expect(device.PrivateKey).ToNot(BeNil(), "Private key should not be nil")
					Expect(device.SignatureCounter).To(Equal(0), "Signature counter should be 0")
					Expect(device.LastSignature).To(BeEmpty(), "Last signature should be empty")
				}
			})
		})
	})

	Describe("AfterSignUpdateDevice", func() {
		var deviceID uuid.UUID

		BeforeEach(func() {
			// Create a device before each context
			algorithm := "RSA"

			device := model.Device{
				ID:               uuid.New(),
				Algorithm:        algorithm,
				Label:            "Test Device",
				PublicKey:        &rsa.PublicKey{},
				PrivateKey:       &rsa.PrivateKey{},
				SignatureCounter: 0,
			}
			err := deviceRepo.Create(device)
			if err != nil {
				Fail(fmt.Sprintf("Failed setting up the device: %v", err))
			}
			deviceID = device.ID
		})

		Context("when updating a device after signning", func() {
			It("should update only the counter and the last signature", func() {
				lastSignature := "test_signature"
				err := deviceRepo.AfterSignUpdateDevice(deviceID, lastSignature)
				Expect(err).To(BeNil(), "Failed to update device after signing")

				updatedDevice, err := deviceRepo.FindByID(deviceID)
				Expect(err).To(BeNil(), "Failed to find updated device")
				Expect(updatedDevice.SignatureCounter).To(Equal(1), "Signature counter should be incremented to 1")
				Expect(updatedDevice.LastSignature).To(Equal(lastSignature), "Last signature should match the provided signature")
			})
		})
	})
})
