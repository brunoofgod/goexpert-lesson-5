# Etapa de build
FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

# Garante que as dependências estejam atualizadas
RUN go mod tidy

# Compila o binário sem dependências de GLIBC
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server-b ./cmd/server-b/main.go

# Etapa final: Usa Debian e instala certificados SSL
FROM debian:bookworm

WORKDIR /root/

# Instala pacotes
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/server-b .

# Expõe as portas dos serviços
EXPOSE 8080
EXPOSE 8085

CMD ["/root/server"]
