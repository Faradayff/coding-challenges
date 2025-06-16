package api

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/model"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/utils"
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeviceService", func() {
	var (
		mockService *domain.MockDeviceService
		mockUtils   *utils.MockUtils
		deviceApi   *DeviceApi
	)

	BeforeEach(func() {
		// Inicitialize the service mock
		mockService = &domain.MockDeviceService{}
		mockUtils = &utils.MockUtils{}
		deviceApi = NewDeviceApi(mockService, mockUtils)
	})

	Describe("CreateSignatureDevice", func() {
		Context("when both params are correct", func() {
			It("should return a CreateDeviceResponse", func() {
				// Mock the CreateSignatureDevice function
				mockService.CreateSignatureDeviceFunc = func(ctx context.Context, algorithm, label string) (model.Device, error) {
					return model.Device{
						ID:        uuid.New(),
						Algorithm: algorithm,
						Label:     label,
					}, nil
				}

				// Prepare the request
				algorithm := "RSA"
				label := "TestDevice"
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/new-device?algorithm=%s&label=%s", algorithm, label), nil)

				// Call the handler
				deviceApi.CreateSignatureDevice(w, r)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusCreated), "Expected status code 201 Created")

				var wrapper struct {
					Data CreateDeviceResponse `json:"data"`
				}

				Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")

				// Check the device returned
				Expect(wrapper.Data).ToNot(BeNil(), "Expected non-nil data in response")
				Expect(wrapper.Data).To(BeAssignableToTypeOf(CreateDeviceResponse{}), "Expected data to be of type CreateDeviceResponse")
				Expect(wrapper.Data.ID).ToNot(BeEmpty(), "Expected non-empty device ID")
				Expect(wrapper.Data.Algorithm).To(Equal("RSA"), "Expected algorithm to be RSA")
				Expect(wrapper.Data.Label).To(Equal("TestDevice"), "Expected label to match input")
			})
		})
		Context("when the param algorithm is empty", func() {
			It("should return an error", func() {
				// Prepare the request
				algorithm := ""
				label := "TestDevice"
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/new-device?algorithm=%s&label=%s", algorithm, label), nil)

				// Call the handler
				deviceApi.CreateSignatureDevice(w, r)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusBadRequest), "Expected status code 400 Bad Request")
				Expect(w.Body.String()).To(ContainSubstring("Missing required parameter: algorithm"), "Expected error message for missing algorithm")
			})
		})

		Context("when the algorithm is invalid", func() {
			It("should return an error", func() {
				// Prepare the request
				algorithm := "INVALID"
				label := "TestDevice"
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/new-device?algorithm=%s&label=%s", algorithm, label), nil)

				// Call the handler
				deviceApi.CreateSignatureDevice(w, r)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusBadRequest), "Expected status code 400 Bad Request")
				Expect(w.Body.String()).To(ContainSubstring("Invalid algorithm. Must be 'ECC' or 'RSA'"), "Expected error message for invalid algorithm")
			})
		})

		Context("when the label is empty", func() {
			It("should return an error", func() {
				// Prepare the request
				algorithm := "RSA"
				label := ""
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/new-device?algorithm=%s&label=%s", algorithm, label), nil)

				// Call the handler
				deviceApi.CreateSignatureDevice(w, r)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusBadRequest), "Expected status code 400 Bad Request")
				Expect(w.Body.String()).To(ContainSubstring("Missing required parameter: label"), "Expected error message for missing label")
			})
		})
	})

	Describe("SignTransaction", func() {
		Context("when all the params are correct", func() {
			It("should return the data signed", func() {
				id := uuid.New()
				// Mock the SignTransaction function with proper key values
				mockService.SignTransactionFunc = func(ctx context.Context, id uuid.UUID, data string) (model.SignaturedData, error) {
					return model.SignaturedData{
						Signature:  []byte("mock-signature"),
						SignedData: "this is signed",
					}, nil
				}

				// Prepare the request
				w := httptest.NewRecorder()
				body := `{"data":"data to sign"}`
				r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/sign?deviceId=%s", id), strings.NewReader(body))
				r.Header.Set("Content-Type", "application/json")

				// Call the handler
				deviceApi.SignTransaction(w, r)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusOK), "Expected status code 200 OK")

				var wrapper struct {
					Data SignaturedDataResponse `json:"data"`
				}

				Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")

				// Check the data returned
				Expect(wrapper.Data).ToNot(BeNil(), "Expected non-nil data in response")
				Expect(wrapper.Data).To(BeAssignableToTypeOf(SignaturedDataResponse{}), "Expected data to be of type SignaturedDataResponse")
				Expect(wrapper.Data.Signature).ToNot(BeEmpty(), "Expected non-empty signature")
				Expect(wrapper.Data.SignedData).ToNot(BeEmpty(), "Expected non-empty signed data")
				Expect(wrapper.Data.SignedData).To(Equal("this is signed"), "Expected signed data to match the mock response")
			})
		})
	})

	Describe("GetDevice", func() {
		Context("when the device exists", func() {
			It("should return the device", func() {
				id := uuid.New()
				// Mock the GetDevice function to return a device with the given ID
				mockService.GetDeviceFunc = func(ctx context.Context, id uuid.UUID) (model.Device, error) {
					return model.Device{
						ID:        id,
						Algorithm: "RSA",
						Label:     "Test Device",
					}, nil
				}

				// Prepare the request
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/device?deviceId=%s", id), nil)

				// Call the handler
				deviceApi.GetDevice(w, r)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusOK), "Expected status code 200 OK")

				var wrapper struct {
					Data GetDeviceResponse `json:"data"`
				}

				Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")

				// Check the device returned
				Expect(wrapper.Data).ToNot(BeNil(), "Expected non-nil data in response")
				Expect(wrapper.Data).To(BeAssignableToTypeOf(GetDeviceResponse{}), "Expected data to be of type DeviceResponse")
				Expect(wrapper.Data.ID).To(Equal(id), "Expected device ID to match the requested ID")
				Expect(wrapper.Data.Algorithm).To(Equal("RSA"), "Expected algorithm to be RSA")
				Expect(wrapper.Data.Label).To(Equal("Test Device"), "Expected label to match input")
			})
		})

		Context("when the device id is not valid", func() {
			It("should return error", func() {
				// Prepare the request with an invalid UUID
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/device?deviceId=invalid-uuid", nil)

				// Call the handler
				deviceApi.GetDevice(w, r)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusBadRequest), "Expected status code 400 Bad Request")
				Expect(w.Body.String()).To(ContainSubstring("Invalid id. Must be a valid UUID"), "Expected error message for invalid UUID")
			})
		})
	})

	Describe("GetAllDevices", func() {
		It("should return all the devices", func() {
			// Mock the GetAllDevices function to return two devices
			mockService.GetAllDevicesFunc = func(ctx context.Context) ([]model.Device, error) {
				return []model.Device{
					{
						ID:         uuid.New(),
						Algorithm:  "RSA",
						Label:      "Test Device",
						PublicKey:  &rsa.PublicKey{},
						PrivateKey: &rsa.PrivateKey{},
					},
					{
						ID:         uuid.New(),
						Algorithm:  "ECC",
						Label:      "Test Device 2",
						PublicKey:  &ecdsa.PublicKey{},
						PrivateKey: &ecdsa.PrivateKey{},
					},
				}, nil
			}

			// Prepare the request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/all", nil)

			// Call the handler
			deviceApi.GetAllDevices(w, r)

			// Verify response code
			Expect(w.Code).To(Equal(http.StatusOK), "Expected status code 200 OK")

			var wrapper struct {
				Data GetAllDevicesResponse `json:"data"`
			}

			Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")

			// Check the device returned
			Expect(wrapper.Data).ToNot(BeNil(), "Expected non-nil data in response")
			Expect(wrapper.Data).To(BeAssignableToTypeOf(GetAllDevicesResponse{}), "Expected data to be of type GetAllDevicesResponse")
			Expect(wrapper.Data.Devices).ToNot(BeEmpty(), "Expected non-empty devices list")
			Expect(len(wrapper.Data.Devices)).To(Equal(2), "Expected two devices in the response")
		})
	})
})
