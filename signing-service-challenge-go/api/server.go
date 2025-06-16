package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/utils"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Response is the generic API response container.
type Response struct {
	Data any `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	api           *DeviceApi
	service       *domain.DeviceService
	repo          *persistence.DeviceRepository
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string) *Server {
	// Initialize persistence layer
	repo := persistence.NewDeviceRepository()

	// Initialize utils
	utils := utils.RealUtils{}

	// Initialize user service
	service := domain.NewDeviceService(repo, &utils, nil)

	// Initialize device API
	api := NewDeviceApi(service, &utils)
	return &Server{
		listenAddress: listenAddress,
		api:           api,
		service:       service,
		repo:          repo,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	// Initialize health service
	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))

	// Initialize Swagger documentation
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Create a subrouter for device-related routes
	deviceMux := http.NewServeMux()
	deviceMux.Handle("POST /new-device", http.HandlerFunc(s.api.CreateSignatureDevice))
	deviceMux.Handle("GET /sign", http.HandlerFunc(s.api.SignTransaction))
	deviceMux.Handle("GET /", http.HandlerFunc(s.api.GetDevice))
	deviceMux.Handle("GET /all", http.HandlerFunc(s.api.GetAllDevices))

	// Add the device prefix
	mux.Handle("/api/v0/device/", http.StripPrefix("/api/v0/device", deviceMux))

	log.Printf("Server running at %s", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	if err != nil {
		WriteInternalError(w)
	}
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	_, err = w.Write(bytes)
	if err != nil {
		WriteInternalError(w)
	}
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data any) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	_, err = w.Write(bytes)
	if err != nil {
		WriteInternalError(w)
	}
}
