package utils

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

type UtilsInterface interface {
	RSAPublicKeyToString(publicKey any) (string, error)
	ECCPublicKeyToString(publicKey any) (string, error)
	RSAPrivateKeyToString(privateKey any) (string, error)
	ECCPrivateKeyToString(privateKey any) (string, error)
	GenerateNewKeyPair(algorithm string) (any, any, error)
}

type RealUtils struct{}

// Convert ECC public key to PEM string
func (u *RealUtils) ECCPublicKeyToString(publicKey any) (string, error) {
	publicKeyECC, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to assert type of ECC public key")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKeyECC)
	if err != nil {
		return "", err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM), nil
}

// Convert ECC private key to PEM string
func (u *RealUtils) ECCPrivateKeyToString(privateKey any) (string, error) {
	privateKeyECC, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("failed to assert type of ECC private key")
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKeyECC)
	if err != nil {
		return "", err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return string(privateKeyPEM), nil
}

// Convert RSA public key to PEM string
func (u *RealUtils) RSAPublicKeyToString(publicKey any) (string, error) {
	publicKeyRSA, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to assert type of RSA public key")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKeyRSA)
	if err != nil {
		return "", err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM), nil
}

// Convert RSA private key to PEM string
func (u *RealUtils) RSAPrivateKeyToString(privateKey any) (string, error) {
	privateKeyRSA, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("failed to assert type of RSA private key")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKeyRSA)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return string(privateKeyPEM), nil
}

// GenerateNewKeyPair generates a new key pair based on the specified algorithm
func (u *RealUtils) GenerateNewKeyPair(algorithm string) (any, any, error) {
	var publicKey, privateKey any
	if algorithm == "ECC" {
		eccGenerator := crypto.ECCGenerator{}

		eccKeys, err := eccGenerator.Generate()
		if err != nil {
			return nil, nil, err
		}

		publicKey = eccKeys.Public
		privateKey = eccKeys.Private
	} else if algorithm == "RSA" {
		rsaGenerator := crypto.RSAGenerator{}

		rsaKeys, err := rsaGenerator.Generate()
		if err != nil {
			return nil, nil, err
		}

		publicKey = rsaKeys.Public
		privateKey = rsaKeys.Private
	} else {
		return nil, nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	return publicKey, privateKey, nil
}
