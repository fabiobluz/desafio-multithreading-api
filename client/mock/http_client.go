package mock

import (
	"errors"
	"net/http"
)

// MockHTTPClient para testes
type MockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
}

// NewMockHTTPClient cria uma nova instância do mock
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		responses: make(map[string]*http.Response),
		errors:    make(map[string]error),
	}
}

// Get implementa a interface HTTPClient
func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if err, exists := m.errors[url]; exists {
		return nil, err
	}

	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}

	return nil, errors.New("URL não encontrada no mock")
}

// SetResponse configura uma resposta para uma URL específica
func (m *MockHTTPClient) SetResponse(url string, response *http.Response) {
	m.responses[url] = response
}

// SetError configura um erro para uma URL específica
func (m *MockHTTPClient) SetError(url string, err error) {
	m.errors[url] = err
}
