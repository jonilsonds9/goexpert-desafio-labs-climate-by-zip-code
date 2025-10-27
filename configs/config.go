package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	WeatherAPIKey        string
	OpenWeatherMapAPIKey string
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	weatherAPIKey := viper.GetString("WEATHERAPI_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHERAPI_KEY is required. Please set it in .env file or environment variables")
	}

	openWeatherMapAPIKey := viper.GetString("OPENWEATHERMAP_API_KEY")
	if openWeatherMapAPIKey == "" {
		log.Fatal("OPENWEATHERMAP_API_KEY is required. Please set it in .env file or environment variables")
	}

	return &Config{
		WeatherAPIKey:        weatherAPIKey,
		OpenWeatherMapAPIKey: openWeatherMapAPIKey,
	}
}
