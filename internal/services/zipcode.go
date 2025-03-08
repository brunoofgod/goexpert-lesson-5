package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

// GetCityByZip consulta o ViaCEP para obter a cidade do CEP
func GetCityByZipOnViaCEP(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("server-b")
	ctx, span := tracer.Start(ctx, "GetCityByZipOnViaCEP")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	// Criando um client HTTP com OpenTelemetry
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("can not find zipcode")
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Localidade == "" {
		return "", fmt.Errorf("invalid zipcode")
	}

	return data.Localidade, nil
}
