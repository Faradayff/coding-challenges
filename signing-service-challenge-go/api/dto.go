package api

// CreateSignatureDeviceResponse defines the response structure for CreateSignatureDevice
type CreateSignatureDeviceResponse struct {
	ID        ID     `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label,omitempty"`
}

// SignatureResponse defines the response structure for SignTransaction
type SignatureResponse struct {
	DeviceID  ID     `json:"deviceId"`
	Algorithm string `json:"algorithm"`
	Signature string `json:"signature"`
}

type ID string
