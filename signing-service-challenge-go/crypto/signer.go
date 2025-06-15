package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

type ECCSigner struct{}

func NewECCSigner() *ECCSigner {
	return &ECCSigner{}
}

func (s *ECCSigner) Sign(data string, privateKey, publicKey any) ([]byte, error) {
	// Cast private and public key from any to *rsa.PrivateKey
	privateKeyCasted, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return []byte{}, fmt.Errorf("failed to assert type of RSA private key")
	}
	publicKeyCasted, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return []byte{}, fmt.Errorf("failed to assert type of RSA public key")
	}

	// Calculate the SHA-256 hash of the data
	hashed := sha256.Sum256([]byte(data))

	// Sign data
	signature, err := ecdsa.SignASN1(rand.Reader, privateKeyCasted, hashed[:])
	if err != nil {
		return []byte{}, fmt.Errorf("failed signing data: %w", err)
	}

	// Verify if the signature is valid
	valid := ecdsa.VerifyASN1(publicKeyCasted, hashed[:], signature)
	if !valid {
		return []byte{}, fmt.Errorf("failed verifying the signed data: %w", err)
	}

	return signature, nil
}

type RSASigner struct{}

func NewRSASigner() *RSASigner {
	return &RSASigner{}
}

func (s *RSASigner) Sign(data string, privateKey, publicKey any) ([]byte, error) {
	// Cast private and public key from any to *rsa.PrivateKey
	privateKeyCasted, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return []byte{}, fmt.Errorf("failed to assert type of RSA private key")
	}
	publicKeyCasted, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return []byte{}, fmt.Errorf("failed to assert type of RSA public key")
	}

	// Calculate the SHA-256 hash of the data
	hashed := sha256.Sum256([]byte(data))

	// Sign data
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKeyCasted, crypto.SHA256, hashed[:])
	if err != nil {
		return []byte{}, fmt.Errorf("failed signing data: %w", err)
	}

	// Verify if the signature is valid
	err = rsa.VerifyPKCS1v15(publicKeyCasted, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return []byte{}, fmt.Errorf("failed verifying the signed data: %w", err)
	}

	return signature, nil
}
