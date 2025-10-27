package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/clients/openweathermap"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/domain"
)

type stubViaCEP struct {
	addr   *domain.ViaCEPAddress
	status int
	err    error
}

func (s *stubViaCEP) ConsultCEP(ctx context.Context, cep string) (*domain.ViaCEPAddress, int, error) {
	return s.addr, s.status, s.err
}

type stubWeather struct {
	tempC  float64
	status int
	err    error
}

func (s *stubWeather) CurrentTempC(ctx context.Context, location string) (float64, int, error) {
	return s.tempC, s.status, s.err
}

func (s *stubWeather) CurrentTempCByCoords(ctx context.Context, lat, lon float64) (float64, int, error) {
	return s.tempC, s.status, s.err
}

type stubGeoClient struct {
	location *openweathermap.GeoLocation
	status   int
	err      error
}

func (s *stubGeoClient) GetCoordinates(ctx context.Context, cityName, countryCode string) (*openweathermap.GeoLocation, int, error) {
	return s.location, s.status, s.err
}

func TestWeatherHandler_Success(t *testing.T) {
	geoLoc := &openweathermap.GeoLocation{
		Name:    "Sao Paulo",
		Lat:     -23.5505,
		Lon:     -46.6333,
		Country: "BR",
		State:   "São Paulo",
	}
	h := NewWeatherHandler(
		&stubViaCEP{addr: &domain.ViaCEPAddress{Localidade: "Sao Paulo", Uf: "SP"}, status: 200},
		&stubGeoClient{location: geoLoc, status: 200},
		&stubWeather{tempC: 25.0, status: 200},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01153000", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp domain.WeatherResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.TempC != 25.0 {
		t.Fatalf("unexpected tempC: %v", resp.TempC)
	}
}

func TestWeatherHandler_InvalidCEP(t *testing.T) {
	h := NewWeatherHandler(&stubViaCEP{}, &stubGeoClient{}, &stubWeather{})
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=abc", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}

func TestWeatherHandler_NotFound(t *testing.T) {
	h := NewWeatherHandler(&stubViaCEP{addr: &domain.ViaCEPAddress{Erro: true}, status: 200}, &stubGeoClient{}, &stubWeather{})
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01153000", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestWeatherHandler_ViaCEPError(t *testing.T) {
	h := NewWeatherHandler(
		&stubViaCEP{addr: nil, status: 500, err: errors.New("service unavailable")},
		&stubGeoClient{},
		&stubWeather{tempC: 25.0, status: 200},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01153000", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestWeatherHandler_GeoClientError(t *testing.T) {
	h := NewWeatherHandler(
		&stubViaCEP{addr: &domain.ViaCEPAddress{Localidade: "Sao Paulo", Uf: "SP"}, status: 200},
		&stubGeoClient{location: nil, status: 500, err: errors.New("geocoding service unavailable")},
		&stubWeather{tempC: 25.0, status: 200},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01153000", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestWeatherHandler_WeatherAPIError(t *testing.T) {
	geoLoc := &openweathermap.GeoLocation{
		Name:    "Sao Paulo",
		Lat:     -23.5505,
		Lon:     -46.6333,
		Country: "BR",
		State:   "São Paulo",
	}
	h := NewWeatherHandler(
		&stubViaCEP{addr: &domain.ViaCEPAddress{Localidade: "Sao Paulo", Uf: "SP"}, status: 200},
		&stubGeoClient{location: geoLoc, status: 200},
		&stubWeather{tempC: 0, status: 500, err: errors.New("weather service unavailable")},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01153000", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestWeatherHandler_MethodNotAllowed(t *testing.T) {
	h := NewWeatherHandler(&stubViaCEP{}, &stubGeoClient{}, &stubWeather{})
	req := httptest.NewRequest(http.MethodPost, "/weather?cep=01153000", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestWeatherHandler_EmptyCEP(t *testing.T) {
	h := NewWeatherHandler(&stubViaCEP{}, &stubGeoClient{}, &stubWeather{})
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}

func TestWeatherHandler_CEPWithHyphen(t *testing.T) {
	h := NewWeatherHandler(&stubViaCEP{}, &stubGeoClient{}, &stubWeather{})
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01153-000", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}

func TestWeatherHandler_TemperatureConversion(t *testing.T) {
	geoLoc := &openweathermap.GeoLocation{
		Name:    "Sao Paulo",
		Lat:     -23.5505,
		Lon:     -46.6333,
		Country: "BR",
		State:   "São Paulo",
	}
	h := NewWeatherHandler(
		&stubViaCEP{addr: &domain.ViaCEPAddress{Localidade: "Sao Paulo", Uf: "SP"}, status: 200},
		&stubGeoClient{location: geoLoc, status: 200},
		&stubWeather{tempC: 28.5, status: 200},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01153000", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp domain.WeatherResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.TempC != 28.5 {
		t.Errorf("expected tempC 28.5, got %v", resp.TempC)
	}
	if resp.TempF != 83.3 {
		t.Errorf("expected tempF 83.3, got %v", resp.TempF)
	}
	if resp.TempK != 301.5 {
		t.Errorf("expected tempK 301.5, got %v", resp.TempK)
	}
}

func TestWeatherHandler_NegativeTemperature(t *testing.T) {
	geoLoc := &openweathermap.GeoLocation{
		Name:    "Moscow",
		Lat:     55.7558,
		Lon:     37.6173,
		Country: "RU",
		State:   "Moscow",
	}
	h := NewWeatherHandler(
		&stubViaCEP{addr: &domain.ViaCEPAddress{Localidade: "Moscow", Uf: "MC"}, status: 200},
		&stubGeoClient{location: geoLoc, status: 200},
		&stubWeather{tempC: -10.0, status: 200},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=12345678", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp domain.WeatherResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.TempC != -10.0 {
		t.Errorf("expected tempC -10.0, got %v", resp.TempC)
	}
	if resp.TempF != 14.0 {
		t.Errorf("expected tempF 14.0, got %v", resp.TempF)
	}
	if resp.TempK != 263.0 {
		t.Errorf("expected tempK 263.0, got %v", resp.TempK)
	}
}
