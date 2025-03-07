package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeather_InvalidCEP(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]string{"cep": "123"})
	req := httptest.NewRequest("POST", "/weather", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	GetWeather(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Esperado status 422, mas recebeu %d", res.StatusCode)
	}
}

func TestGetWeather_NotFoundCEP(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]string{"cep": "00000000"})
	req := httptest.NewRequest("POST", "/weather", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	GetWeather(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Esperado status 404, mas recebeu %d", res.StatusCode)
	}
}
