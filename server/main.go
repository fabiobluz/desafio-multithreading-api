package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Resultado struct {
	Source string                 `json:"source"`
	Data   map[string]interface{} `json:"data"`
	Error  string                 `json:"error,omitempty"`
}

type APIClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DefaultAPIClient struct{}

func (c *DefaultAPIClient) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func main() {
	http.HandleFunc("/consulta", handler)
	fmt.Println("🚀 Servidor rodando em http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "cep é obrigatório", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resultChan := make(chan Resultado, 2)

	brasilAPI := "https://brasilapi.com.br/api/cep/v1/" + cep
	viaCepAPI := "http://viacep.com.br/ws/" + cep + "/json/"

	apiClient := &DefaultAPIClient{}
	go buscarAPI(ctx, "BrasilAPI", brasilAPI, resultChan, apiClient)
	go buscarAPI(ctx, "ViaCEP", viaCepAPI, resultChan, apiClient)

	var res Resultado
	select {
	case res = <-resultChan:
	case <-ctx.Done():
		http.Error(w, "timeout: nenhuma API respondeu em 1s", http.StatusGatewayTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func buscarAPI(ctx context.Context, source, url string, ch chan<- Resultado, apiClient APIClient) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- Resultado{Source: source, Error: err.Error()}
		return
	}

	resp, err := apiClient.Do(req)
	if err != nil {
		ch <- Resultado{Source: source, Error: err.Error()}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- Resultado{Source: source, Error: err.Error()}
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		ch <- Resultado{Source: source, Error: err.Error()}
		return
	}

	select {
	case ch <- Resultado{Source: source, Data: data}:
	case <-ctx.Done():
	}
}
