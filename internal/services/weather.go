package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// GetWeatherByCity consulta a WeatherAPI para obter a temperatura atual
func GetWeatherByCity(ctx context.Context, client *http.Client, city *string) (*WeatherResponse, error) {
	tracer := otel.Tracer("server-b")
	ctx, span := tracer.Start(ctx, "GetWeatherByCity")
	defer span.End()

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API Key não configurada")
	}
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&key=%s", url.QueryEscape(*city), apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler resposta de erro: %v", err)
		}
		return nil, fmt.Errorf("erro ao buscar clima body: %s KEY: %s", string(bodyBytes), apiKey)
	}

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	response := WeatherResponse{
		TempC: data.Current.TempC,
		TempF: data.Current.TempC*1.8 + 32,
		TempK: data.Current.TempC + 273,
	}

	return &response, nil
}
