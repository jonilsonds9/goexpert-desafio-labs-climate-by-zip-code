package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                 string
	WeatherAPIKey        string
	OpenWeatherMapAPIKey string
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	weatherAPIKey := viper.GetString("WEATHERAPI_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHERAPI_KEY is required. Please set it in .env file or environment variables")
	}

	openWeatherMapAPIKey := viper.GetString("OPENWEATHERMAP_API_KEY")
	if openWeatherMapAPIKey == "" {
		log.Fatal("OPENWEATHERMAP_API_KEY is required. Please set it in .env file or environment variables")
	}

	return &Config{
		Port:                 port,
		WeatherAPIKey:        weatherAPIKey,
		OpenWeatherMapAPIKey: openWeatherMapAPIKey,
	}
}
