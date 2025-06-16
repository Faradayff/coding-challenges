package domain

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/utils"
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeviceService", func() {
	var (
		mockSigner     *crypto.MockSigner
		mockUtils      *utils.MockUtils
		mockDeviceRepo *persistence.MockDeviceRepo
		deviceService  *DeviceService
	)

	BeforeEach(func() {
		// Inicitialize the device repository, mock service and mock utils before each test
		mockUtils = &utils.MockUtils{}
		mockDeviceRepo = &persistence.MockDeviceRepo{}
		deviceService = NewDeviceService(mockDeviceRepo, mockUtils, mockSigner)
	})

	Describe("CreateSignatureDevice", func() {
		Context("when creating a new device", func() {
			It("should return a new device", func() {
				device, err := deviceService.CreateSignatureDevice(context.Background(), "ECC", "Test ECC Device")
				Expect(err).To(BeNil(), "Failed to create device")
				Expect(device).To(BeAssignableToTypeOf(model.Device{}), "The created device should be of type model.Device")
				Expect(device.ID).ToNot(BeEmpty(), "The device ID should not be empty")
				Expect(device.Algorithm).To(Equal("ECC"), "The device algorithm should be ECC")
				Expect(device.Label).To(Equal("Test ECC Device"), "The device label should match the input")
				Expect(device.PublicKey).To(BeAssignableToTypeOf(&ecdsa.PublicKey{}), "The public key should be of type *ecdsa.PublicKey")
				Expect(device.PrivateKey).To(BeAssignableToTypeOf(&ecdsa.PrivateKey{}), "The private key should be of type *ecdsa.PrivateKey")
				Expect(device.SignatureCounter).To(Equal(0), "The signature counter should be initialized to 0")
			})

			It("should create a device with RSA algorithm", func() {
				device, err := deviceService.CreateSignatureDevice(context.Background(), "RSA", "Test RSA Device")
				Expect(err).To(BeNil(), "Failed to create RSA device")
				Expect(device).To(BeAssignableToTypeOf(model.Device{}), "The created device should be of type model.Device")
				Expect(device.ID).ToNot(BeEmpty(), "The device ID should not be empty")
				Expect(device.Algorithm).To(Equal("RSA"), "The device algorithm should be RSA")
				Expect(device.Label).To(Equal("Test RSA Device"), "The device label should match the input")
				Expect(device.PublicKey).To(BeAssignableToTypeOf(&rsa.PublicKey{}), "The public key should be of type *rsa.PublicKey")
				Expect(device.PrivateKey).To(BeAssignableToTypeOf(&rsa.PrivateKey{}), "The private key should be of type *rsa.PrivateKey")
				Expect(device.SignatureCounter).To(Equal(0), "The signature counter should be initialized to 0")
			})
		})

		Context("when repo is not working", func() {
			It("should fail when creating the device", func() {
				faillingDeviceService := NewDeviceService(nil, mockUtils, mockSigner)
				Expect(func() {
					_, _ = faillingDeviceService.CreateSignatureDevice(context.Background(), "ECC", "Test ECC Device")
				}).To(Panic(), "The device service should panic when the repository is nil")
			})

		})
	})

	Describe("SignTransaction", func() {
		Context("when the device exists", func() {
			It("should increment the signature counter", func() {
				id := uuid.New()
				// Mock the device repository to return a device with the given ID
				mockDeviceRepo.FindByIDFunc = func(id uuid.UUID) (*model.Device, error) {
					return &model.Device{
						ID:               id,
						Algorithm:        "ECC",
						Label:            "Test Device",
						PublicKey:        &ecdsa.PublicKey{},
						PrivateKey:       &ecdsa.PrivateKey{},
						SignatureCounter: 0,
						LastSignature:    "",
					}, nil
				}
				// Sign data
				signaturedData, err := deviceService.SignTransaction(context.Background(), id, "test data to sign")
				Expect(err).To(BeNil(), "Failed to sign")
				Expect(signaturedData).To(BeAssignableToTypeOf(model.SignaturedData{}), "The signed data should be of type model.SignaturedData")
				Expect(signaturedData.Signature).To(Not(BeEmpty()), "The signature should not be empty")
				Expect(signaturedData.SignedData).To(Not(BeEmpty()), "The signed data should not be empty")
			})
		})

		Context("when the device does not exist", func() {
			It("should return an error", func() {
				// Mock the device repository to return a device with the given ID
				mockDeviceRepo.FindByIDFunc = func(id uuid.UUID) (*model.Device, error) {
					return nil, errors.New("device not found")
				}
				_, err := deviceService.SignTransaction(context.Background(), uuid.New(), "test data")
				Expect(err).To(HaveOccurred(), "Signing a transaction with a non-existent device should return an error")
				Expect(err.Error()).To(ContainSubstring("device not found"), "The error message should indicate that the device was not found")
			})
		})

		Context("when multiple transactions are signed concurrently", func() {
			It("should handle concurrent transactions correctly", func() {
				id := uuid.New()
				numTransactions := 10
				wg := sync.WaitGroup{}
				wg.Add(numTransactions)

				// Run concurrent sign petitions
				for range numTransactions {
					go func() {
						defer wg.Done()
						signaturedData, err := deviceService.SignTransaction(context.Background(), id, "test data to sign")
						Expect(err).To(BeNil(), "Failed to sign")
						Expect(signaturedData).To(BeAssignableToTypeOf(model.SignaturedData{}), "The signed data should be of type model.SignaturedData")
						Expect(signaturedData.Signature).To(Not(BeEmpty()), "The signature should not be empty")
						Expect(signaturedData.SignedData).To(Not(BeEmpty()), "The signed data should not be empty")
					}()
				}

				wg.Wait()
			})
		})
	})
})
