package mock

import (
	"errors"
	"net/http"
	"time"
)

// MockAPIClient para testes
type MockAPIClient struct {
	responses map[string]*http.Response
	errors    map[string]error
	delays    map[string]time.Duration
}

// NewMockAPIClient cria uma nova instância do mock
func NewMockAPIClient() *MockAPIClient {
	return &MockAPIClient{
		responses: make(map[string]*http.Response),
		errors:    make(map[string]error),
		delays:    make(map[string]time.Duration),
	}
}

// Do implementa a interface APIClient
func (m *MockAPIClient) Do(req *http.Request) (*http.Response, error) {
	url := req.URL.String()

	// Simular delay se configurado
	if delay, exists := m.delays[url]; exists {
		time.Sleep(delay)
	}

	if err, exists := m.errors[url]; exists {
		return nil, err
	}

	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}

	return nil, errors.New("URL não encontrada no mock")
}

// SetResponse configura uma resposta para uma URL específica
func (m *MockAPIClient) SetResponse(url string, response *http.Response) {
	m.responses[url] = response
}

// SetError configura um erro para uma URL específica
func (m *MockAPIClient) SetError(url string, err error) {
	m.errors[url] = err
}

// SetDelay configura um delay para uma URL específica
func (m *MockAPIClient) SetDelay(url string, delay time.Duration) {
	m.delays[url] = delay
}
