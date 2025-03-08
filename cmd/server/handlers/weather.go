package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/brunoofgod/goexpert-lesson-5/internal/services"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type WeatherRequest struct {
	CEP string `json:"cep"`
}

// GetWeather processa a requisição do usuário
// @Summary Obtém a temperatura de uma cidade a partir do CEP
// @Description Retorna a temperatura em Celsius, Fahrenheit e Kelvin
// @Tags Clima
// @Accept json
// @Produce json
// @Param request body WeatherRequest true "CEP para consulta"
// @Success 200 {object} services.WeatherResponse
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /weather [post]
func GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("server")

	// Criando um span para a requisição do usuário
	ctx, span := tracer.Start(ctx, "GetWeather-handler")
	defer span.End()

	var req WeatherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(req.CEP) != 8 {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	bodyBytes, shouldReturn := getTemperatures(ctx, req, w)
	if shouldReturn {
		return
	}

	var response services.WeatherResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getTemperatures(ctx context.Context, req WeatherRequest, w http.ResponseWriter) ([]byte, bool) {
	url := fmt.Sprintf("%s/get-temperature-by-zipcode?zipcode=%s", os.Getenv("SERVER_B_HOST"), req.CEP)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(httpReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, true
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, true
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, string(bodyBytes), resp.StatusCode)
		return nil, true
	}
	return bodyBytes, false
}
