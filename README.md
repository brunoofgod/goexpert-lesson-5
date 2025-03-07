# Clima API com serviços separados utilizando Open telemetry

## Descrição
Este projeto é uma API em Go que recebe um CEP como entrada, identifica a cidade correspondente e retorna a temperatura atual em três unidades: Celsius, Fahrenheit e Kelvin.

### **Tecnologias e Pacotes Utilizados**
- **Linguagem**: Go
- **Framework**: `go-chi/chi` (roteamento)
- **Documentação**: `swaggo/swag`
- **Serviços Externos**:
  - [ViaCEP](https://viacep.com.br/) para busca de informações do CEP
  - [WeatherAPI](https://www.weatherapi.com/) para obter dados meteorológicos
- **Docker** para execução e deploy

## **URL Pública**
O projeto está publicado na seguinte URL:
[https://goexpert-lesson-5-480219334057.us-central1.run.app/](https://goexpert-lesson-5-480219334057.us-central1.run.app/)

## **Funcionamento**
A API aceita requisições HTTP POST com um corpo JSON contendo um CEP de 8 dígitos.

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
      - HOSTNAME=localhost:8080
   ```

3. Inicie o projeto com Docker Compose:
   ```sh
    docker-compose up -d
   ```

4. Acesse a API em:
   ```
   API Clima: http://localhost:8080
   API Cidade: http://localhost:8085
   ```
   Para ver a documentação Swagger:
   ```
   http://localhost:8080/swagger/index.html
   http://localhost:8085/swagger/index.html
   ```