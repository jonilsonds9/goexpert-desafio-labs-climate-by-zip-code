# Desafio Labs - Clima por CEP

Sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin).

## Endereço do projeto no Google Cloud Run

```
https://SERVICE_APP/api/weather?cep=01153000
```

## Como rodar o projeto localmente

### Pré-requisitos

- Go 1.21 ou superior
- Docker e Docker Compose (opcional)
- Chaves de API:
  - [WeatherAPI](https://www.weatherapi.com/) - Para consultar temperatura
  - [OpenWeatherMap](https://openweathermap.org/api) - Para geocoding (obter lat/lon)

### Configuração

1. Clone o repositório:
```bash
git clone <seu-repositorio>
cd labs-climate-by-zip-code
```

2. Crie um arquivo `.env` na raiz do projeto com suas chaves de API:
```env
WEATHERAPI_KEY=sua_chave_weatherapi_aqui
OPENWEATHERMAP_API_KEY=sua_chave_openweathermap_aqui
```

> **Nota**: Você pode copiar o arquivo `.env.example` e preencher com suas chaves:
> ```bash
> cp .env.example .env
> ```

### Executando com Docker Compose (RECOMENDADO)

O Docker Compose automaticamente carrega as variáveis do arquivo `.env`:

```bash
docker-compose up --build
```

O servidor estará disponível em `http://localhost:8080`

### Executando com Go

```bash
go run cmd/server/main.go
```

O servidor estará disponível em `http://localhost:8080`

### Testando a aplicação

```bash
# Teste com CEP válido
curl "http://localhost:8080/api/weather?cep=01153000"

# Resposta esperada:
# {"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}
```

### Executando os testes

```bash
go test ./...
```

## API Endpoints

### GET /api/weather

Consulta a temperatura atual para um CEP brasileiro.

**Parâmetros:**
- `cep` (query string, obrigatório): CEP brasileiro com 8 dígitos (pode conter ou não hífen)

**Exemplos:**
```bash
# Com hífen
curl "http://localhost:8080/api/weather?cep=01153-000"

# Sem hífen
curl "http://localhost:8080/api/weather?cep=01153000"
```

**Respostas:**

✅ **Sucesso (200)**
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

❌ **CEP inválido (422)**
```json
invalid zipcode
```

❌ **CEP não encontrado (404)**
```json
can not find zipcode
```

## Descrição do desafio

**Objetivo**: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

**Requisitos**:

- O sistema deve receber um CEP válido de 8 digitos
- O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin.
- O sistema deve responder adequadamente nos seguintes cenários:
  - Em caso de sucesso:
    - Código HTTP: 200
    - Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
  - Em caso de falha, caso o CEP não seja válido (com formato correto):
    - Código HTTP: 422
    - Mensagem: invalid zipcode
  - Em caso de falha, caso o CEP não seja encontrado:
    - Código HTTP: 404
    - Mensagem: can not find zipcode
- Deverá ser realizado o deploy no Google Cloud Run.

**Dicas**:

- Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a temperatura: https://viacep.com.br/
- Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas: https://www.weatherapi.com/
- Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32
- Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
  - Sendo F = Fahrenheit
  - Sendo C = Celsius
  - Sendo K = Kelvin

**Entrega**:

- O código-fonte completo da implementação.
- Testes automatizados demonstrando o funcionamento.
- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- Deploy realizado no Google Cloud Run (free tier) e endereço ativo para ser acessado.
