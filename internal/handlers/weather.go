package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/clients/openweathermap"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/clients/viacep"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/clients/weatherapi"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/domain"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/utils"
)

type WeatherHandler struct {
	viaCEP     viacep.Client
	geoClient  openweathermap.Client
	weatherAPI weatherapi.Client
}

func NewWeatherHandler(viaCEP viacep.Client, geoClient openweathermap.Client, weather weatherapi.Client) http.Handler {
	return &WeatherHandler{viaCEP: viaCEP, geoClient: geoClient, weatherAPI: weather}
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cep := strings.TrimSpace(r.URL.Query().Get("cep"))
	if !utils.IsValidCEP(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	ctx, cancel := contextWithTimeout(r, 5*time.Second)
	defer cancel()

	addr, status, err := h.viaCEP.ConsultCEP(ctx, cep)
	if err != nil || status >= 400 || addr == nil || addr.Erro {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// Get coordinates from OpenWeatherMap Geocoding API
	geoLocation, statusGeo, err := h.geoClient.GetCoordinates(ctx, addr.Localidade, "BR")
	if err != nil || statusGeo >= 400 || geoLocation == nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// Get temperature using coordinates
	tempC, statusW, err := h.weatherAPI.CurrentTempCByCoords(ctx, geoLocation.Lat, geoLocation.Lon)
	if err != nil || statusW >= 400 {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	resp := domain.WeatherResponse{
		TempC: tempC,
		TempF: utils.CelsiusToFahrenheit(tempC),
		TempK: utils.CelsiusToKelvin(tempC),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func contextWithTimeout(r *http.Request, d time.Duration) (context.Context, context.CancelFunc) {
	if r.Context() != nil {
		return context.WithTimeout(r.Context(), d)
	}
	return context.WithTimeout(context.Background(), d)
}
