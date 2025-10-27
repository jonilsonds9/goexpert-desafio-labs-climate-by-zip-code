package utils

import (
	"regexp"
)

var cepRegex = regexp.MustCompile(`^\d{8}$`)

func IsValidCEP(cep string) bool {
	return cepRegex.MatchString(cep)
}

func CelsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func CelsiusToKelvin(c float64) float64 {
	return c + 273
}
