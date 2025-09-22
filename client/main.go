package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Resultado struct {
	Source string                 `json:"source"`
	Data   map[string]interface{} `json:"data"`
	Error  string                 `json:"error,omitempty"`
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <CEP>")
		return
	}
	cep := os.Args[1]

	client := &DefaultHTTPClient{}
	if err := consultarCEP(cep, client); err != nil {
		fmt.Println("Erro:", err)
		return
	}
}

func consultarCEP(cep string, client HTTPClient) error {
	url := "http://localhost:8080/consulta?cep=" + cep
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("erro ao chamar o servidor: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro do servidor: %s", resp.Status)
	}

	var resultado Resultado
	if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
		return fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	fmt.Println("✅ Resposta recebida da:", resultado.Source)
	for k, v := range resultado.Data {
		fmt.Printf("  %s: %v\n", k, v)
	}

	return nil
}
