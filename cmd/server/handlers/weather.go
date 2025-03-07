package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/brunoofgod/goexpert-lesson-5/internal/services"
)

type WeatherRequest struct {
	CEP string `json:"cep"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// GetWeather processa a requisição do usuário
// @Summary Obtém a temperatura de uma cidade a partir do CEP
// @Description Retorna a temperatura em Celsius, Fahrenheit e Kelvin
// @Tags Clima
// @Accept json
// @Produce json
// @Param request body WeatherRequest true "CEP para consulta"
// @Success 200 {object} WeatherResponse
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /weather [post]
func GetWeather(w http.ResponseWriter, r *http.Request) {
	var req WeatherRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(req.CEP) != 8 {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	city, err := getCityByZipCode(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	tempC, err := services.GetWeatherByCity(http.DefaultClient, city)
	if err != nil {
		http.Error(w, `{"message": "error fetching weather"}`, http.StatusInternalServerError)
		return
	}

	response := WeatherResponse{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getCityByZipCode(req WeatherRequest, w http.ResponseWriter) (*string, error) {

	url := fmt.Sprintf("%s/get-city-by-zip?zipcode=%s", os.Getenv("SERVER_B_HOST"), req.CEP)
	resp, err := http.DefaultClient.Get(url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`{"message": "%s"}`, bodyString)
	}

	if err != nil {
		return nil, err
	}

	return &bodyString, nil
}
