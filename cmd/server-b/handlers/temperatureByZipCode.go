package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brunoofgod/goexpert-lesson-5/internal/services"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

// GetTemperatureByZipCode processa a requisição do usuário
// @Summary Obtém as temperaturas de uma cidade a partir do CEP
// @Description Retorna os graus de temperatura em Celsius, Fahrenheit e Kelvin
// @Tags Temperature
// @Accept json
// @Produce json
// @Param zipcode query string true "CEP para consulta"
// @Success 200 {object} services.WeatherResponse
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /get-temperature-by-zipcode [get]
func GetTemperatureByZipCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("server-b")
	ctx, span := tracer.Start(ctx, "GetTemperatureByZipCode-handler")
	defer span.End()

	zipCode := r.URL.Query().Get("zipcode")
	if zipCode == "" {
		http.Error(w, "ZipCode is required", http.StatusBadRequest)
		return
	}

	cityName, err := services.GetCityByZipOnViaCEP(ctx, zipCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	response, err := services.GetWeatherByCity(ctx, client, &cityName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
