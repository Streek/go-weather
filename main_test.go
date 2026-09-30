package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReadResponseBodyRejectsOversizedResponse(t *testing.T) {
	_, err := readResponseBody(strings.NewReader(strings.Repeat("x", maxAPIResponse+1)))
	if err == nil {
		t.Fatal("readResponseBody accepted an oversized response")
	}
}

func TestCheckCacheRejectsFutureTimestamp(t *testing.T) {
	cacheKey := "future-timestamp-test"
	cachePath := filepath.Join(getCacheDir(), cacheKey+".json")
	t.Cleanup(func() { _ = os.Remove(cachePath) })
	data := CacheFile{Timestamp: time.Now().Add(time.Hour), Data: WeatherData{}}
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cachePath, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	if _, ok := checkCache(cacheKey); ok {
		t.Fatal("checkCache accepted a future-dated cache")
	}
}

func TestValidateWeatherDataRejectsMismatchedArrays(t *testing.T) {
	var weather WeatherData
	weather.CurrentWeather.Time = "2026-01-01T00:00"
	weather.Daily.Time = []string{"2026-01-01"}
	if err := validateWeatherData(weather, false, false); err == nil {
		t.Fatal("validateWeatherData accepted mismatched daily arrays")
	}
}

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
