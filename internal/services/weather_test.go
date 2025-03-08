package services

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func mockWeatherAPI() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current": {"temp_c": 25.0}}`))
	})
	return httptest.NewServer(handler)
}

type mockRoundTripper struct {
	mockServerURL string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	mockURL, err := url.Parse(m.mockServerURL)
	if err != nil {
		return nil, err
	}

	req.URL.Scheme = mockURL.Scheme
	req.URL.Host = mockURL.Host
	req.RequestURI = ""

	return http.DefaultTransport.RoundTrip(req)
}

func TestGetWeatherByCity(t *testing.T) {
	mockServer := mockWeatherAPI()
	defer mockServer.Close()

	os.Setenv("WEATHER_API_KEY", "fake-key")
	defer os.Unsetenv("WEATHER_API_KEY")

	mockClient := &http.Client{
		Transport: &mockRoundTripper{mockServerURL: mockServer.URL},
	}

	city := "São Paulo"
	temp, err := GetWeatherByCity(mockClient, &city)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	if temp.TempC != 25.0 {
		t.Errorf("Esperava 25.0°C, mas recebeu %f", temp)
	}
}
