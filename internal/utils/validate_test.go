package utils

import (
	"math"
	"testing"
)

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{
			name:     "Valid CEP with 8 digits",
			cep:      "01153000",
			expected: true,
		},
		{
			name:     "Valid CEP another example",
			cep:      "12345678",
			expected: true,
		},
		{
			name:     "Invalid CEP with hyphen",
			cep:      "01153-000",
			expected: false,
		},
		{
			name:     "Invalid CEP with less than 8 digits",
			cep:      "0115300",
			expected: false,
		},
		{
			name:     "Invalid CEP with more than 8 digits",
			cep:      "011530000",
			expected: false,
		},
		{
			name:     "Invalid CEP with letters",
			cep:      "0115300a",
			expected: false,
		},
		{
			name:     "Invalid empty CEP",
			cep:      "",
			expected: false,
		},
		{
			name:     "Invalid CEP with spaces",
			cep:      "01153 000",
			expected: false,
		},
		{
			name:     "Invalid CEP with special characters",
			cep:      "01153@00",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCEP(tt.cep)
			if result != tt.expected {
				t.Errorf("IsValidCEP(%s) = %v, expected %v", tt.cep, result, tt.expected)
			}
		})
	}
}

func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name      string
		celsius   float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "Freezing point",
			celsius:   0,
			expected:  32,
			tolerance: 0.01,
		},
		{
			name:      "Boiling point",
			celsius:   100,
			expected:  212,
			tolerance: 0.01,
		},
		{
			name:      "Negative temperature",
			celsius:   -40,
			expected:  -40,
			tolerance: 0.01,
		},
		{
			name:      "Room temperature",
			celsius:   25,
			expected:  77,
			tolerance: 0.01,
		},
		{
			name:      "Body temperature",
			celsius:   37,
			expected:  98.6,
			tolerance: 0.01,
		},
		{
			name:      "Decimal temperature",
			celsius:   28.5,
			expected:  83.3,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CelsiusToFahrenheit(tt.celsius)
			if !floatEquals(result, tt.expected, tt.tolerance) {
				t.Errorf("CelsiusToFahrenheit(%f) = %f, expected %f", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "Freezing point",
			celsius:  0,
			expected: 273,
		},
		{
			name:     "Boiling point",
			celsius:  100,
			expected: 373,
		},
		{
			name:     "Absolute zero",
			celsius:  -273,
			expected: 0,
		},
		{
			name:     "Room temperature",
			celsius:  25,
			expected: 298,
		},
		{
			name:     "Decimal temperature",
			celsius:  28.5,
			expected: 301.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CelsiusToKelvin(tt.celsius)
			if result != tt.expected {
				t.Errorf("CelsiusToKelvin(%f) = %f, expected %f", tt.celsius, result, tt.expected)
			}
		})
	}
}
