package main

import (
	"testing"
)

func TestGetTempUnit(t *testing.T) {
	tests := []struct {
		name       string
		unitSystem UnitSystem
		want       string
	}{
		{"metric", UnitMetric, "°C"},
		{"imperial", UnitImperial, "°F"},
		{"empty", "", "°C"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getTempUnit(tc.unitSystem)
			if got != tc.want {
				t.Errorf("getTempUnit(%s) = %s; want %s", tc.unitSystem, got, tc.want)
			}
		})
	}
}

func TestColorizeTemp(t *testing.T) {
	// Simple test just to check it doesn't crash
	result := colorizeTemp(20.0, UnitMetric)
	if result == "" {
		t.Error("colorizeTemp returned empty string")
	}
}

// Add more tests for other key functions
