package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/configs"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/clients/openweathermap"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/clients/viacep"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/clients/weatherapi"
	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/handlers"
)

func main() {
	cfg := configs.LoadConfig()

	httpClient := &http.Client{Timeout: 10 * time.Second}

	viaCEPClient := viacep.NewClient(httpClient)
	geoClient := openweathermap.NewClient(httpClient, cfg.OpenWeatherMapAPIKey)
	weatherClient := weatherapi.NewClient(httpClient, cfg.WeatherAPIKey)

	mux := http.NewServeMux()
	mux.Handle("/api/weather", handlers.NewWeatherHandler(viaCEPClient, geoClient, weatherClient))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, logRequests(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
