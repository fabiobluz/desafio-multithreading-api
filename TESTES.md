# Testes Unitários - Desafio Multithreading API

Este documento descreve os testes unitários criados para as aplicações Client e Server.

## Estrutura dos Testes

### Servidor (`server/main_test.go` + `server/mock/`)

Os testes do servidor cobrem:

1. **TestHandler_Success**: Testa o handler com sucesso, verificando se retorna dados válidos de uma das APIs
2. **TestHandler_MissingCEP**: Testa validação de parâmetro obrigatório (CEP)
3. **TestHandler_Timeout**: Testa comportamento de timeout quando APIs não respondem
4. **TestBuscarAPI_Success**: Testa a função `buscarAPI` com resposta válida
5. **TestBuscarAPI_Error**: Testa tratamento de erros de rede
6. **TestBuscarAPI_InvalidJSON**: Testa tratamento de JSON inválido

#### Melhorias de Testabilidade

- **Interface APIClient**: Criada para permitir injeção de dependência e mock de requisições HTTP
- **Pacote `server/mock`**: Contém `MockAPIClient` que permite simular diferentes cenários:
  - Respostas de sucesso
  - Erros de rede
  - Delays para simular timeouts
- **Função helper `handlerWithMock`**: Permite testar o handler com cliente HTTP mockado

### Cliente (`client/main_test.go` + `client/mock/`)

Os testes do cliente cobrem:

1. **TestConsultarCEP_Success**: Testa consulta bem-sucedida e exibição dos dados
2. **TestConsultarCEP_HTTPError**: Testa tratamento de erros de rede
3. **TestConsultarCEP_ServerError**: Testa tratamento de erros do servidor (status != 200)
4. **TestConsultarCEP_InvalidJSON**: Testa tratamento de JSON inválido
5. **TestConsultarCEP_EmptyData**: Testa resposta com dados vazios
6. **TestConsultarCEP_WithError**: Testa resposta com campo de erro
7. **TestDefaultHTTPClient_Get**: Testa o cliente HTTP padrão

#### Melhorias de Testabilidade

- **Interface HTTPClient**: Criada para permitir injeção de dependência
- **Pacote `client/mock`**: Contém `MockHTTPClient` para simular respostas HTTP
- **Função `consultarCEP`**: Extraída do main para permitir testes unitários
- **Captura de output**: Função `captureOutput` para verificar saída do console

## Como Executar os Testes

### Executar todos os testes do servidor:
```bash
cd server
go test -v
```

### Executar todos os testes do cliente:
```bash
cd client
go test -v
```

### Executar um teste específico:
```bash
# Servidor
go test -v -run TestHandler_Success

# Cliente
go test -v -run TestConsultarCEP_Success
```

### Executar com cobertura:
```bash
# Servidor
go test -v -cover

# Cliente
go test -v -cover
```

## Cobertura de Testes

Os testes cobrem os seguintes cenários:

### Servidor:
- ✅ Validação de parâmetros
- ✅ Requisições HTTP bem-sucedidas
- ✅ Tratamento de erros de rede
- ✅ Tratamento de JSON inválido
- ✅ Timeout de requisições
- ✅ Multithreading (goroutines)

### Cliente:
- ✅ Requisições HTTP bem-sucedidas
- ✅ Tratamento de erros de rede
- ✅ Tratamento de erros do servidor
- ✅ Parsing de JSON
- ✅ Exibição de dados no console
- ✅ Tratamento de dados vazios

## Estrutura de Arquivos

```
server/
├── main.go              # Código principal do servidor
├── main_test.go         # Testes unitários
├── mock/
│   └── api_client.go    # MockAPIClient para testes
└── go.mod

client/
├── main.go              # Código principal do cliente
├── main_test.go         # Testes unitários
├── mock/
│   └── http_client.go   # MockHTTPClient para testes
└── go.mod
```

## Padrões Utilizados

1. **Dependency Injection**: Uso de interfaces para permitir mocks
2. **Table-Driven Tests**: Estrutura organizada para múltiplos cenários
3. **Mock Objects**: Simulação de dependências externas em pacotes separados
4. **Test Helpers**: Funções auxiliares para reduzir duplicação
5. **Output Capture**: Verificação de saída do console
6. **Separation of Concerns**: Mocks organizados em pacotes dedicados

## Melhorias Futuras

1. **Testes de Integração**: Testes end-to-end entre cliente e servidor
2. **Benchmarks**: Testes de performance
3. **Testes de Concorrência**: Verificação de race conditions
4. **Mocks mais sofisticados**: Simulação de diferentes tipos de erro
5. **Cobertura de código**: Análise detalhada de cobertura
