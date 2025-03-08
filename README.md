# Clima API com serviços separados utilizando Open telemetry

## Descrição
Este projeto é composto por 2 API's em Go que recebe um CEP como entrada, identifica a cidade correspondente e retorna a temperatura atual em três unidades: Celsius, Fahrenheit e Kelvin.

### **Tecnologias e Pacotes Utilizados**
- **Linguagem**: Go
- **Framework**: `go-chi/chi` (roteamento)
- **Documentação**: `swaggo/swag`
- **Open Telemetry**: `opentelemetry-go-contrib` 
- **Serviços Externos**:
  - [ViaCEP](https://viacep.com.br/) para busca de informações do CEP
  - [WeatherAPI](https://www.weatherapi.com/) para obter dados meteorológicos
  - [Zipkin](https://zipkin.io/)
- **Docker** para execução e deploy

## **Funcionamento**
A API denominada server aceita requisições HTTP POST com um corpo JSON contendo um CEP de 8 dígitos. A API nomeada server-b é chamada pela API server para fazer o processamento dos dados e obter a temperatura atual. 

### **Exemplo de Request:**
```json
POST /weather
Content-Type: application/json
{
  "cep": "88310630"
}
```

### **Exemplo de Response:**
#### **Em caso de sucesso (200 OK):**
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### **CEP inválido (422 Unprocessable Entity):**
```json
{
  "message": "invalid zipcode"
}
```

#### **CEP não encontrado (404 Not Found):**
```json
{
  "message": "can not find zipcode"
}
```

## **Pré-requisitos para funcionamento**
- Docker e Docker Compose instalados
- Chave de API do [WeatherAPI](https://www.weatherapi.com/) (definir no ambiente)

## **Instruções de Execução**
### **Rodando com Docker Compose**
1. Clone o repositório:
   ```sh
   git clone https://github.com/brunoofgod/goexpert-lesson-5.git
   cd goexpert-lesson-5
   ```

2. Configure seus enviroments dentro do arquivo docker-compose.yml:
   ```sh
      environment:
      - WEATHER_API_KEY=[SUA-API-KEY]
      - HOSTNAME=localhost:8085
   ```

3. Inicie o projeto com Docker Compose:
   ```sh
    docker-compose up -d
   ```

4. Acesse a API em:
   ```
   API server: http://localhost:8080
   API server-b: http://localhost:8085
   Zipkin: http://localhost:9411/zipkin/
   ```
   Para ver a documentação Swagger:
   ```
   http://localhost:8080/swagger/index.html
   http://localhost:8085/swagger/index.html
   ```