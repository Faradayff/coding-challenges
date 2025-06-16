package crypto

type MockSigner struct {
	SignFunc func(data string, privateKey, publicKey any) ([]byte, error)
}

func (m *MockSigner) Sign(data string, privateKey, publicKey any) ([]byte, error) {
	if m.SignFunc != nil {
		return m.SignFunc(data, privateKey, publicKey)
	}
	return []byte("mock-signature"), nil
}
