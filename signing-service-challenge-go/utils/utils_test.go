package utils

import (
	"crypto/ecdsa"
	"crypto/rsa"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	var (
		u *RealUtils
	)

	BeforeEach(func() {
		u = &RealUtils{}
	})
	Describe("Generate keys", func() {

		Context("when creating new ECC key pair", func() {
			It("should return both keys", func() {
				algorithm := "ECC"
				publicKey, privateKey, err := u.GenerateNewKeyPair(algorithm)
				Expect(err).To(BeNil(), "Failed to generate ECC key pair")
				Expect(publicKey).ToNot(BeNil(), "Public key should not be nil")
				Expect(privateKey).ToNot(BeNil(), "Private key should not be nil")
				Expect(publicKey).To(BeAssignableToTypeOf(&ecdsa.PublicKey{}), "Public key should be of type *ecdsa.PublicKey")
				Expect(privateKey).To(BeAssignableToTypeOf(&ecdsa.PrivateKey{}), "Private key should be of type *ecdsa.PrivateKey")
			})
		})

		Context("when creating new RSA key pair", func() {
			It("should return both keys", func() {
				algorithm := "RSA"
				publicKey, privateKey, err := u.GenerateNewKeyPair(algorithm)
				Expect(err).To(BeNil(), "Failed to generate ECC key pair")
				Expect(publicKey).ToNot(BeNil(), "Public key should not be nil")
				Expect(privateKey).ToNot(BeNil(), "Private key should not be nil")
				Expect(publicKey).To(BeAssignableToTypeOf(&rsa.PublicKey{}), "Public key should be of type *ecdsa.PublicKey")
				Expect(privateKey).To(BeAssignableToTypeOf(&rsa.PrivateKey{}), "Private key should be of type *ecdsa.PrivateKey")
			})
		})
	})

	Describe("Keys to string", func() {

		Context("when converting ECC keys to string", func() {
			var publicKey, privateKey any

			BeforeEach(func() {
				var err error
				publicKey, privateKey, err = u.GenerateNewKeyPair("ECC")
				Expect(err).To(BeNil(), "Setup failed in BeforeEach")
			})

			Context("when converting public key to string", func() {
				It("should return a string with the key", func() {
					publicKeyStr, err := u.ECCPublicKeyToString(publicKey)
					Expect(err).To(BeNil(), "Failed to convert ECC public key to string")
					Expect(publicKeyStr).To(BeAssignableToTypeOf(string("")), "Public key string should be of type string")
					Expect(publicKeyStr).ToNot(BeEmpty(), "Public key string should not be empty")
					Expect(publicKeyStr).To(ContainSubstring("BEGIN EC PUBLIC KEY"), "Public key string should contain 'BEGIN EC PUBLIC KEY'")
					Expect(publicKeyStr).To(ContainSubstring("END EC PUBLIC KEY"), "Public key string should contain 'END EC PUBLIC KEY'")
				})
			})
			Context("when converting private key to string", func() {
				It("should return a string with the key", func() {
					privateKeyStr, err := u.ECCPrivateKeyToString(privateKey)
					Expect(err).To(BeNil(), "Failed to convert ECC private key to string")
					Expect(privateKeyStr).To(BeAssignableToTypeOf(string("")), "Private key string should be of type string")
					Expect(privateKeyStr).ToNot(BeEmpty(), "Private key string should not be empty")
					Expect(privateKeyStr).To(ContainSubstring("BEGIN EC PRIVATE KEY"), "Private key string should contain 'BEGIN EC PRIVATE KEY'")
					Expect(privateKeyStr).To(ContainSubstring("END EC PRIVATE KEY"), "Private key string should contain 'END EC PRIVATE KEY'")
				})
			})
		})

		Context("when converting RSA keys to string", func() {
			var publicKey, privateKey any

			BeforeEach(func() {
				var err error
				publicKey, privateKey, err = u.GenerateNewKeyPair("RSA")
				Expect(err).To(BeNil(), "Setup failed in BeforeEach")
			})

			Context("when converting public key to string", func() {
				It("should return a string with the key", func() {
					publicKeyStr, err := u.RSAPublicKeyToString(publicKey)
					Expect(err).To(BeNil(), "Failed to convert RSA public key to string")
					Expect(publicKeyStr).To(BeAssignableToTypeOf(string("")), "Public key string should be of type string")
					Expect(publicKeyStr).ToNot(BeEmpty(), "Public key string should not be empty")
					Expect(publicKeyStr).To(ContainSubstring("BEGIN RSA PUBLIC KEY"), "Public key string should contain 'BEGIN RSA PUBLIC KEY'")
					Expect(publicKeyStr).To(ContainSubstring("END RSA PUBLIC KEY"), "Public key string should contain 'END RSA PUBLIC KEY'")
				})
			})
			Context("when converting private key to string", func() {
				It("should return a string with the key", func() {
					privateKeyStr, err := u.RSAPrivateKeyToString(privateKey)
					Expect(err).To(BeNil(), "Failed to convert RSA private key to string")
					Expect(privateKeyStr).To(BeAssignableToTypeOf(string("")), "Private key string should be of type string")
					Expect(privateKeyStr).ToNot(BeEmpty(), "Private key string should not be empty")
					Expect(privateKeyStr).To(ContainSubstring("BEGIN RSA PRIVATE KEY"), "Private key string should contain 'BEGIN RSA PRIVATE KEY'")
					Expect(privateKeyStr).To(ContainSubstring("END RSA PRIVATE KEY"), "Private key string should contain 'END RSA PRIVATE KEY'")
				})
			})
		})
	})
})
