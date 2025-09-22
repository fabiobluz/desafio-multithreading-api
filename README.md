# 🚀 Desafio Multithreading API

Uma aplicação em Go que demonstra o uso de **multithreading** e **concorrência** para buscar informações de endereço a partir de um CEP, utilizando duas APIs públicas simultaneamente e retornando a resposta mais rápida.

## 📋 Sobre o Projeto

Esta aplicação implementa um padrão de **"race condition"** controlado, onde duas APIs são consultadas simultaneamente para o mesmo CEP, e apenas a resposta mais rápida é retornada ao usuário. Isso melhora significativamente a performance e a experiência do usuário.

### 🎯 Objetivos
- Demonstrar o uso de **goroutines** e **channels** em Go
- Implementar **timeout** para evitar requisições muito lentas
- Mostrar **injeção de dependência** para facilitar testes
- Aplicar boas práticas de **tratamento de erros**

## 🏗️ Arquitetura

A aplicação é composta por dois módulos principais:

```
desafio-multithreading-api/
├── server/          # Servidor HTTP que consulta as APIs
├── client/          # Cliente que consome o servidor
└── TESTES.md        # Documentação dos testes unitários
```

### 🔄 Fluxo de Funcionamento

1. **Cliente** envia CEP para o **Servidor**
2. **Servidor** dispara 2 goroutines simultâneas:
   - Uma para **BrasilAPI** (`https://brasilapi.com.br/api/cep/v1/`)
   - Outra para **ViaCEP** (`http://viacep.com.br/ws/`)
3. O servidor aguarda a **primeira resposta** (race condition)
4. Retorna os dados da API mais rápida
5. **Cliente** exibe os dados formatados

## 🚀 Como Executar

### Pré-requisitos
- Go 1.19+ instalado
- Conexão com a internet (para acessar as APIs)

### 1. Executar o Servidor

```bash
cd server
go run main.go
```

O servidor estará disponível em: `http://localhost:8080`

### 2. Executar o Cliente

Em outro terminal:

```bash
cd client
go run main.go 01310100
```

**Exemplo de saída:**
```
✅ Resposta recebida da: BrasilAPI
  cep: 01310-100
  state: SP
  city: São Paulo
  neighborhood: Bela Vista
  street: Avenida Paulista
```

## 🔧 API Endpoints

### GET `/consulta?cep={cep}`

Consulta informações de endereço para um CEP específico.

**Parâmetros:**
- `cep` (obrigatório): CEP no formato `01310100` ou `01310-100`

**Resposta de Sucesso (200):**
```json
{
  "source": "BrasilAPI",
  "data": {
    "cep": "01310-100",
    "state": "SP",
    "city": "São Paulo",
    "neighborhood": "Bela Vista",
    "street": "Avenida Paulista"
  }
}
```

**Resposta de Erro (400):**
```json
{
  "error": "cep é obrigatório"
}
```

**Resposta de Timeout (504):**
```json
{
  "error": "timeout: nenhuma API respondeu em 1s"
}
```

## ⚡ Características Técnicas

### 🧵 Multithreading
- **Goroutines**: Duas requisições simultâneas para APIs diferentes
- **Channels**: Comunicação entre goroutines de forma segura
- **Select**: Aguarda a primeira resposta disponível

### ⏱️ Timeout
- **Timeout de 1 segundo**: Evita requisições muito lentas
- **Context cancellation**: Cancela requisições em andamento

### 🧪 Testabilidade
- **Interfaces**: `APIClient` e `HTTPClient` para injeção de dependência
- **Mocks**: Pacotes dedicados para simulação em testes
- **Cobertura**: Testes unitários para todas as funcionalidades

## 🧪 Testes

A aplicação possui uma suíte completa de testes unitários:

### Executar Testes do Servidor
```bash
cd server
go test -v
```

### Executar Testes do Cliente
```bash
cd client
go test -v
```

### Cobertura de Testes
```bash
# Servidor
cd server && go test -v -cover

# Cliente
cd client && go test -v -cover
```

**Cenários testados:**
- ✅ Requisições bem-sucedidas
- ✅ Tratamento de erros de rede
- ✅ Timeout de requisições
- ✅ Validação de parâmetros
- ✅ Parsing de JSON
- ✅ Comportamento de multithreading

## 📊 APIs Utilizadas

### 1. BrasilAPI
- **URL**: `https://brasilapi.com.br/api/cep/v1/{cep}`
- **Formato**: JSON
- **Campos**: `cep`, `state`, `city`, `neighborhood`, `street`

### 2. ViaCEP
- **URL**: `http://viacep.com.br/ws/{cep}/json/`
- **Formato**: JSON
- **Campos**: `cep`, `logradouro`, `bairro`, `localidade`, `uf`

## 🛠️ Tecnologias Utilizadas

- **Go 1.24.2**: Linguagem principal
- **net/http**: Cliente e servidor HTTP
- **context**: Controle de timeout e cancelamento
- **encoding/json**: Serialização/deserialização JSON
- **testing**: Framework de testes

## 📁 Estrutura de Arquivos

```
server/
├── main.go              # Servidor HTTP principal
├── main_test.go         # Testes unitários
├── mock/
│   └── api_client.go    # Mock para testes
└── go.mod

client/
├── main.go              # Cliente HTTP
├── main_test.go         # Testes unitários
├── mock/
│   └── http_client.go   # Mock para testes
└── go.mod
```

## 🎯 Padrões de Design Aplicados

1. **Dependency Injection**: Interfaces para facilitar testes
2. **Race Condition Controlado**: Primeira resposta vence
3. **Timeout Pattern**: Evita requisições lentas
4. **Error Handling**: Tratamento robusto de erros
5. **Separation of Concerns**: Mocks em pacotes separados

## 🚀 Melhorias Futuras

- [ ] **Cache**: Implementar cache de respostas
- [ ] **Logging**: Adicionar logs estruturados
- [ ] **Métricas**: Monitoramento de performance
- [ ] **Rate Limiting**: Controle de taxa de requisições
- [ ] **Health Check**: Endpoint de saúde da aplicação
- [ ] **Docker**: Containerização da aplicação

## 📝 Exemplo de Uso Completo

```bash
# Terminal 1 - Iniciar servidor
cd server
go run main.go
# 🚀 Servidor rodando em http://localhost:8080 ...

# Terminal 2 - Consultar CEP
cd client
go run main.go 01310100
# ✅ Resposta recebida da: BrasilAPI
#   cep: 01310-100
#   state: SP
#   city: São Paulo
#   neighborhood: Bela Vista
#   street: Avenida Paulista

# Terminal 3 - Teste via curl
curl "http://localhost:8080/consulta?cep=01310100"
# {"source":"BrasilAPI","data":{"cep":"01310-100","state":"SP",...}}
```

## 📄 Licença

Este projeto é um exemplo educacional demonstrando conceitos de concorrência em Go.

---

**Desenvolvido com ❤️ em Go**