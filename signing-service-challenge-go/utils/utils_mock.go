package utils

import (
	"crypto/ecdsa"
	"crypto/rsa"
)

type MockUtils struct {
	ECCPublicKeyToStringFunc  func(publicKey any) (string, error)
	ECCPrivateKeyToStringFunc func(privateKey any) (string, error)
	RSAPublicKeyToStringFunc  func(publicKey any) (string, error)
	RSAPrivateKeyToStringFunc func(privateKey any) (string, error)
	GenerateNewKeyPairFunc    func(algorithm string) (any, any, error)
}

func (m *MockUtils) ECCPublicKeyToString(publicKey any) (string, error) {
	if m.ECCPublicKeyToStringFunc != nil {
		return m.ECCPublicKeyToStringFunc(publicKey)
	}
	return "ECC Public Key String", nil
}

func (m *MockUtils) ECCPrivateKeyToString(privateKey any) (string, error) {
	if m.ECCPrivateKeyToStringFunc != nil {
		return m.ECCPrivateKeyToStringFunc(privateKey)
	}
	return "ECC Private Key String", nil
}

func (m *MockUtils) RSAPublicKeyToString(publicKey any) (string, error) {
	if m.RSAPublicKeyToStringFunc != nil {
		return m.RSAPublicKeyToStringFunc(publicKey)
	}
	return "RSA Public Key String", nil
}

func (m *MockUtils) RSAPrivateKeyToString(privateKey any) (string, error) {
	if m.RSAPrivateKeyToStringFunc != nil {
		return m.RSAPrivateKeyToStringFunc(privateKey)
	}
	return "RSA Private Key String", nil
}

func (m *MockUtils) GenerateNewKeyPair(algorithm string) (any, any, error) {
	if m.GenerateNewKeyPairFunc != nil {
		return m.GenerateNewKeyPairFunc(algorithm)
	}
	switch algorithm {
	case "ECC":
		return &ecdsa.PublicKey{}, &ecdsa.PrivateKey{}, nil
	case "RSA":
		return &rsa.PublicKey{}, &rsa.PrivateKey{}, nil
	default:
		return nil, nil, nil
	}
}
