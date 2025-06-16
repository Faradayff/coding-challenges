package e2etests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"
)

var _ = Describe("Device API End-to-End", func() {
	var (
		deviceRepo    persistence.DeviceRepoInterface
		realUtils     utils.UtilsInterface
		deviceService domain.DeviceServiceInterface
		deviceApi     *api.DeviceApi
		w             *httptest.ResponseRecorder
		deviceID      uuid.UUID
	)

	BeforeEach(func() {
		realUtils = &utils.RealUtils{}
		deviceRepo = persistence.NewDeviceRepository()
		deviceService = domain.NewDeviceService(deviceRepo, realUtils, nil)
		deviceApi = api.NewDeviceApi(deviceService, realUtils)
		w = httptest.NewRecorder()
	})

	Describe("CreateSignatureDevice", func() {
		Context("when all the parameters are valid", func() {
			It("should create a device", func() {
				// Prepare the request
				r := httptest.NewRequest("POST", "/new-device?algorithm=ECC&label=testlabel", nil)

				// Call the handler
				deviceApi.CreateSignatureDevice(w, r)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusCreated), "Expected status code 201 Created")

				var wrapper struct {
					Data api.CreateDeviceResponse `json:"data"`
				}
				Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")
				resp := wrapper.Data

				// Check the response
				Expect(resp.Algorithm).To(Equal("ECC"), "Expected algorithm to be ECC")
				Expect(resp.Label).To(Equal("testlabel"), "Expected label to be 'testlabel'")
				Expect(resp.ID).ToNot(BeEmpty(), "Expected ID to be generated")
				Expect(resp.ID).To(BeAssignableToTypeOf(uuid.UUID{}), "Expected ID to be of type uuid.UUID")
				Expect(resp.PublicKey).ToNot(BeEmpty(), "Expected PublicKey to be generated")
				Expect(resp.PrivateKey).ToNot(BeEmpty(), "Expected PrivateKey to be generated")
			})
		})
	})

	Describe("SignTransaction", func() {

		BeforeEach(func() {
			// Creating a new device
			// Prepare the request
			r := httptest.NewRequest("POST", "/new-device?algorithm=ECC&label=testlabel", nil)

			// Call the handler
			deviceApi.CreateSignatureDevice(w, r)

			// Verify response code
			Expect(w.Code).To(Equal(http.StatusCreated), "Failed during set up")

			var wrapper struct {
				Data api.CreateDeviceResponse `json:"data"`
			}
			Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")
			resp := wrapper.Data

			deviceID = resp.ID

			w = httptest.NewRecorder()
		})

		Context("when the device exists", func() {
			It("should sign data", func() {
				// Prepare the request
				text := "hello Fiskaly!"
				payload := map[string]string{"data": text}
				body, _ := json.Marshal(payload)
				url := fmt.Sprintf("/sign?deviceId=%s", deviceID.String())
				req := httptest.NewRequest("POST", url, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				// Call the handler
				deviceApi.SignTransaction(w, req)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusOK))

				var wrapper struct {
					Data api.SignaturedDataResponse `json:"data"`
				}
				Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")
				resp := wrapper.Data

				// Check the response
				Expect(resp).To(BeAssignableToTypeOf(api.SignaturedDataResponse{}), "Expected data to be of type SignaturedDataResponse")
				Expect(resp.Signature).ToNot(BeEmpty(), "Expected non-empty signature")
				Expect(resp.SignedData).To(ContainSubstring(text), "Expected signed data to match input data")
			})
		})

		Context("when the device does not exist", func() {
			It("should return a 404 error", func() {
				// Prepare the request
				text := "hello Fiskaly!"
				payload := map[string]string{"data": text}
				body, _ := json.Marshal(payload)
				url := fmt.Sprintf("/sign?deviceId=%s", uuid.New().String())
				req := httptest.NewRequest("POST", url, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				// Call the handler
				deviceApi.SignTransaction(w, req)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusNotFound), "Expected status code 404 Not Found")

				var errWrap struct {
					Errors []string `json:"errors"`
				}
				Expect(json.NewDecoder(w.Body).Decode(&errWrap)).To(Succeed(), "Expected to decode error response body without error")

				// Check the response
				Expect(errWrap.Errors).To(ContainElement("Device not found"), "Expected error message to indicate device not found")
			})
		})
	})

	Describe("GetDevice", func() {

		BeforeEach(func() {
			// Creating a new device
			// Prepare the request
			r := httptest.NewRequest("POST", "/new-device?algorithm=RSA&label=testlabel", nil)

			// Call the handler
			deviceApi.CreateSignatureDevice(w, r)

			// Verify response code
			Expect(w.Code).To(Equal(http.StatusCreated), "Failed during set up")

			var wrapper struct {
				Data api.CreateDeviceResponse `json:"data"`
			}
			Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")
			resp := wrapper.Data

			deviceID = resp.ID

			w = httptest.NewRecorder()
		})

		Context("when the device exists", func() {
			It("should retrieve the device", func() {
				// Prepare the request
				url := fmt.Sprintf("/device?deviceId=%s", deviceID.String())
				req := httptest.NewRequest("GET", url, nil)

				// Call the handler
				deviceApi.GetDevice(w, req)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusOK), "Expected status code 200 OK")

				var wrapper struct {
					Data api.GetDeviceResponse `json:"data"`
				}
				Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")
				resp := wrapper.Data

				// Check the response
				Expect(resp.Algorithm).To(Equal("RSA"), "Expected algorithm to be RSA")
				Expect(resp.Label).To(Equal("testlabel"), "Expected label to be 'testlabel'")
				Expect(resp.SignatureCounter).To(Equal(0), "Expected signature counter to be 0")
			})
		})

		Context("when the device does not exist", func() {
			It("should return a 404 error", func() {
				// Prepare the request
				url := fmt.Sprintf("/device?deviceId=%s", uuid.New().String())
				req := httptest.NewRequest("GET", url, nil)

				// Call the handler
				deviceApi.GetDevice(w, req)

				// Verify response code
				Expect(w.Code).To(Equal(http.StatusNotFound), "Expected status code 404 Not Found")

				var errWrap struct {
					Errors []string `json:"errors"`
				}
				Expect(json.NewDecoder(w.Body).Decode(&errWrap)).To(Succeed(), "Expected to decode error response body without error")

				// Check the response
				Expect(errWrap.Errors).To(ContainElement("Device not found"), "Expected error message to indicate device not found")
			})
		})
	})

	Describe("GetAllDevices", func() {
		BeforeEach(func() {
			req1 := httptest.NewRequest("POST", "/new-device?algorithm=ECC&label=first", nil)
			deviceApi.CreateSignatureDevice(w, req1)
			w = httptest.NewRecorder()

			req2 := httptest.NewRequest("POST", "/new-device?algorithm=RSA&label=second", nil)
			deviceApi.CreateSignatureDevice(w, req2)
			w = httptest.NewRecorder()
		})

		It("should list all devices", func() {
			// Prepare the request
			req := httptest.NewRequest("GET", "/all", nil)

			// Call the handler
			deviceApi.GetAllDevices(w, req)

			// Verify response code
			Expect(w.Code).To(Equal(http.StatusOK), "Expected status code 200 OK")

			var wrapper struct {
				Data api.GetAllDevicesResponse `json:"data"`
			}
			Expect(json.NewDecoder(w.Body).Decode(&wrapper)).To(Succeed(), "Expected to decode response body without error")
			resp := wrapper.Data

			// Check the response
			// It should return the devices always in the same order
			Expect(len(resp.Devices)).To(Equal(2), "Expected two devices to be returned")
			Expect(resp.Total).To(Equal(2), "Expected total count of devices to be 2")
			Expect(resp.Devices[0].Algorithm).To(Equal("ECC"), "Expected first device to be ECC")
			Expect(resp.Devices[0].Label).To(Equal("first"), "Expected first device label to be 'first'")
			Expect(resp.Devices[1].Algorithm).To(Equal("RSA"), "Expected second device to be RSA")
			Expect(resp.Devices[1].Label).To(Equal("second"), "Expected second device label to be 'second'")
		})
	})
})
