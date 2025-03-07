package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

// GetCityByZip consulta o ViaCEP para obter a cidade do CEP
func GetCityByZipOnViaCEP(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CEP não encontrado")
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Localidade == "" {
		return "", fmt.Errorf("CEP inválido")
	}

	return data.Localidade, nil
}
