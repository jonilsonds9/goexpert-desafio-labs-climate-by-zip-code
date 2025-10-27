package configs

import (
	"os"
	"testing"
)

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("WEATHERAPI_KEY", "test-weather-key")
	os.Setenv("OPENWEATHERMAP_API_KEY", "test-openweather-key")
	defer os.Unsetenv("WEATHERAPI_KEY")
	defer os.Unsetenv("OPENWEATHERMAP_API_KEY")

	config := LoadConfig()

	if config.WeatherAPIKey != "test-weather-key" {
		t.Errorf("Expected WeatherAPIKey to be 'test-weather-key', got '%s'", config.WeatherAPIKey)
	}

	if config.OpenWeatherMapAPIKey != "test-openweather-key" {
		t.Errorf("Expected OpenWeatherMapAPIKey to be 'test-openweather-key', got '%s'", config.OpenWeatherMapAPIKey)
	}
}
