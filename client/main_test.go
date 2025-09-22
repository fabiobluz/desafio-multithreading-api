package main

import (
	"bytes"
	"client/mock"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestConsultarCEP_Success(t *testing.T) {
	// Dados de resposta simulados
	resultado := Resultado{
		Source: "BrasilAPI",
		Data: map[string]interface{}{
			"cep":          "01310-100",
			"state":        "SP",
			"city":         "São Paulo",
			"neighborhood": "Bela Vista",
			"street":       "Avenida Paulista",
		},
	}

	resultadoJSON, _ := json.Marshal(resultado)

	mockClient := mock.NewMockHTTPClient()
	mockClient.SetResponse("http://localhost:8080/consulta?cep=01310100", &http.Response{
		StatusCode: 200,
		Body:       &mockBody{data: string(resultadoJSON)},
	})

	// Capturar output
	originalOutput := captureOutput(func() {
		err := consultarCEP("01310100", mockClient)
		if err != nil {
			t.Errorf("Não esperado erro: %v", err)
		}
	})

	// Verificar se a saída contém informações esperadas
	if !strings.Contains(originalOutput, "✅ Resposta recebida da: BrasilAPI") {
		t.Errorf("Output não contém mensagem de sucesso esperada")
	}

	if !strings.Contains(originalOutput, "cep: 01310-100") {
		t.Errorf("Output não contém dados do CEP")
	}
}

func TestConsultarCEP_HTTPError(t *testing.T) {
	mockClient := mock.NewMockHTTPClient()
	mockClient.SetError("http://localhost:8080/consulta?cep=01310100", errors.New("connection failed"))

	err := consultarCEP("01310100", mockClient)
	if err == nil {
		t.Error("Esperado erro, mas não ocorreu")
	}

	if !strings.Contains(err.Error(), "erro ao chamar o servidor") {
		t.Errorf("Erro inesperado: %v", err)
	}
}

func TestConsultarCEP_ServerError(t *testing.T) {
	mockClient := mock.NewMockHTTPClient()
	mockClient.SetResponse("http://localhost:8080/consulta?cep=01310100", &http.Response{
		StatusCode: 500,
		Body:       &mockBody{data: "Internal Server Error"},
	})

	err := consultarCEP("01310100", mockClient)
	if err == nil {
		t.Error("Esperado erro, mas não ocorreu")
	}

	if !strings.Contains(err.Error(), "erro do servidor") {
		t.Errorf("Erro inesperado: %v", err)
	}
}

func TestConsultarCEP_InvalidJSON(t *testing.T) {
	mockClient := mock.NewMockHTTPClient()
	mockClient.SetResponse("http://localhost:8080/consulta?cep=01310100", &http.Response{
		StatusCode: 200,
		Body:       &mockBody{data: "invalid json"},
	})

	err := consultarCEP("01310100", mockClient)
	if err == nil {
		t.Error("Esperado erro, mas não ocorreu")
	}

	if !strings.Contains(err.Error(), "erro ao decodificar resposta") {
		t.Errorf("Erro inesperado: %v", err)
	}
}

func TestConsultarCEP_EmptyData(t *testing.T) {
	resultado := Resultado{
		Source: "ViaCEP",
		Data:   map[string]interface{}{},
	}

	resultadoJSON, _ := json.Marshal(resultado)

	mockClient := mock.NewMockHTTPClient()
	mockClient.SetResponse("http://localhost:8080/consulta?cep=00000000", &http.Response{
		StatusCode: 200,
		Body:       &mockBody{data: string(resultadoJSON)},
	})

	output := captureOutput(func() {
		err := consultarCEP("00000000", mockClient)
		if err != nil {
			t.Errorf("Não esperado erro: %v", err)
		}
	})

	if !strings.Contains(output, "✅ Resposta recebida da: ViaCEP") {
		t.Errorf("Output não contém mensagem de sucesso esperada")
	}
}

func TestConsultarCEP_WithError(t *testing.T) {
	resultado := Resultado{
		Source: "BrasilAPI",
		Error:  "CEP não encontrado",
	}

	resultadoJSON, _ := json.Marshal(resultado)

	mockClient := mock.NewMockHTTPClient()
	mockClient.SetResponse("http://localhost:8080/consulta?cep=99999999", &http.Response{
		StatusCode: 200,
		Body:       &mockBody{data: string(resultadoJSON)},
	})

	output := captureOutput(func() {
		err := consultarCEP("99999999", mockClient)
		if err != nil {
			t.Errorf("Não esperado erro: %v", err)
		}
	})

	if !strings.Contains(output, "✅ Resposta recebida da: BrasilAPI") {
		t.Errorf("Output não contém mensagem de sucesso esperada")
	}
}

func TestDefaultHTTPClient_Get(t *testing.T) {
	// Teste com servidor real (opcional - pode ser comentado se não quiser fazer requisições reais)
	client := &DefaultHTTPClient{}

	// Criar um servidor de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "test"})
	}))
	defer server.Close()

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("Erro inesperado: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Esperado status 200, recebido %d", resp.StatusCode)
	}
}

// Helper functions para testes
func captureOutput(fn func()) string {
	// Salvar o stdout original
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Executar a função
	fn()

	// Fechar o writer e restaurar stdout
	w.Close()
	os.Stdout = old

	// Ler a saída capturada
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Mock body para testes
type mockBody struct {
	data string
	pos  int
}

func (m *mockBody) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.data) {
		return 0, fmt.Errorf("EOF")
	}

	n = copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockBody) Close() error {
	return nil
}
