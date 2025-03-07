package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brunoofgod/goexpert-lesson-5/internal/services"
)

// GetCityByZip processa a requisição do usuário
// @Summary Obtém o nome da cidade através do CEP
// @Description Retorna o nome da cidade pelo CEP
// @Tags City
// @Accept json
// @Produce json
// @Param zipcode query string true "CEP para consulta"
// @Success 200 {object} string
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /get-city-by-zip [get]
func GetCityByZip(w http.ResponseWriter, r *http.Request) {

	zipCode := r.URL.Query().Get("zipcode")

	if zipCode == "" {
		http.Error(w, "ZipCode is required", http.StatusBadRequest)
		return
	}

	cityName, err := services.GetCityByZipOnViaCEP(zipCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cityName)
}
