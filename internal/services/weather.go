package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// GetWeatherByCity consulta a WeatherAPI para obter a temperatura atual
func GetWeatherByCity(client *http.Client, city *string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("API Key não configurada")
	}
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&key=%s", url.QueryEscape(*city), apiKey)
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("erro ao ler resposta de erro: %v", err)
		}
		return 0, fmt.Errorf("erro ao buscar clima body: %s KEY: %s", string(bodyBytes), apiKey)
	}

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
