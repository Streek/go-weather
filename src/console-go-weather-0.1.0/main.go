package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Application constants
const (
	configFileName = "weather_config.json"
	cacheDuration  = 1 * time.Hour
	appName        = "Weather Console"
	appVersion     = "1.0.0"
)

// DisplayMode represents how weather data should be presented
type DisplayMode string

// Available display modes
const (
	DisplayText  DisplayMode = "text"
	DisplayTable DisplayMode = "table"
)

// UnitSystem represents measurement units to use
type UnitSystem string

// Available unit systems
const (
	UnitMetric   UnitSystem = "metric"
	UnitImperial UnitSystem = "imperial"
)

// Config stores user preferences
type Config struct {
	ZipCode     string      `json:"zip_code"`
	DisplayMode DisplayMode `json:"display_mode"`
	Units       UnitSystem  `json:"units"`
	UseColors   bool        `json:"use_colors"`
}

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
)

// Main function - entry point for the application
func main() {
	// Parse command line flags and handle commands
	cmd := parseFlags()
	if err := cmd.execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Command represents a user operation
type Command struct {
	showHelp       bool
	showDaily      bool
	showHourly     bool
	zipOverride    string
	displayMode    DisplayMode
	forceTextMode  bool
	forceTableMode bool
	unitSystem     UnitSystem
	useColors      *bool
	noColors       bool
	saveAll        bool // New flag to save all settings
}

// parseFlags processes command-line arguments and returns a Command
func parseFlags() *Command {
	cmd := &Command{}

	// Define command line flags
	flag.BoolVar(&cmd.showHelp, "help", false, "Show help information")
	flag.BoolVar(&cmd.showDaily, "daily", false, "Show 7-day forecast")
	flag.BoolVar(&cmd.showHourly, "hourly", false, "Show hourly forecast")
	flag.StringVar(&cmd.zipOverride, "zip", "", "Override default ZIP/postal code")
	flag.BoolVar(&cmd.forceTableMode, "table", false, "Show output in table format")
	flag.BoolVar(&cmd.forceTextMode, "text", false, "Show output in text format")
	flag.StringVar((*string)(&cmd.unitSystem), "units", "", "Use specific units (metric or imperial)")

	// Add save flag
	flag.BoolVar(&cmd.saveAll, "save", false, "Save current settings as defaults")

	var useColors bool
	flag.BoolVar(&useColors, "color", false, "Enable colored output")
	flag.BoolVar(&cmd.noColors, "no-color", false, "Disable colored output")

	// Add short flag aliases
	flag.BoolVar(&cmd.showHelp, "?", false, "Short for -help")
	flag.BoolVar(&cmd.showDaily, "d", false, "Short for -daily")
	flag.BoolVar(&cmd.showHourly, "h", false, "Short for -hourly")
	flag.StringVar(&cmd.zipOverride, "z", "", "Short for -zip")
	flag.BoolVar(&cmd.forceTableMode, "t", false, "Short for -table")
	flag.BoolVar(&cmd.forceTextMode, "T", false, "Short for -text")
	flag.StringVar((*string)(&cmd.unitSystem), "u", "", "Short for -units")
	flag.BoolVar(&useColors, "c", false, "Short for -color")
	flag.BoolVar(&cmd.noColors, "nc", false, "Short for -no-color")
	flag.BoolVar(&cmd.saveAll, "s", false, "Short for -save")

	// Override default usage output
	flag.Usage = printHelp

	// Parse flags and capture original args before parsing
	originalArgs := os.Args[1:]
	flag.Parse()

	// Process color flags manually by checking if they appeared in the args
	for _, arg := range originalArgs {
		if arg == "-color" || arg == "-c" {
			cmd.useColors = &useColors
			break
		}
	}

	return cmd
}

// execute runs the command based on flags
func (cmd *Command) execute() error {
	if cmd.showHelp {
		printHelp()
		return nil
	}

	// Load config (or create default)
	config := loadConfig()

	// Determine display mode
	displayMode := config.DisplayMode
	if cmd.forceTableMode && !cmd.forceTextMode {
		displayMode = DisplayTable
	} else if cmd.forceTextMode && !cmd.forceTableMode {
		displayMode = DisplayText
	}

	// If no display mode is set, default to text
	if displayMode == "" {
		displayMode = DisplayText
	}

	cmd.displayMode = displayMode

	// Determine unit system
	unitSystem := config.Units
	if cmd.unitSystem != "" {
		unitSystem = cmd.unitSystem
	}

	// If no unit system is set, default to metric
	if unitSystem == "" {
		unitSystem = UnitMetric
	}

	// Handle color settings
	useColors := config.UseColors
	if cmd.useColors != nil {
		useColors = *cmd.useColors
	} else if cmd.noColors {
		useColors = false
	}

	// Get location coordinates
	zipCode := cmd.zipOverride
	if zipCode == "" {
		zipCode = config.ZipCode
	}

	// If no zip code is set, ask the user
	if zipCode == "" {
		fmt.Print("Enter your location (ZIP/postal code or city name): ")
		fmt.Scanln(&zipCode)
	}

	// Handle saving settings if --save flag is provided
	if cmd.saveAll {
		// Only save valid values and not flags starting with -
		if zipCode != "" && !strings.HasPrefix(zipCode, "-") {
			config.ZipCode = zipCode
		}

		// Save display mode if explicitly set
		if cmd.forceTableMode || cmd.forceTextMode {
			config.DisplayMode = displayMode
		}

		// Save unit system if explicitly set
		if cmd.unitSystem != "" {
			config.Units = unitSystem
		}

		// Save color preference if explicitly set
		if cmd.useColors != nil {
			config.UseColors = *cmd.useColors
		} else if cmd.noColors {
			config.UseColors = false
		}

		// Save config to file
		if err := saveConfig(config); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Println("All settings saved:")
		fmt.Printf("- Location: %s\n", config.ZipCode)
		fmt.Printf("- Display mode: %s\n", config.DisplayMode)
		fmt.Printf("- Unit system: %s\n", getUnitSystemName(config.Units))
		fmt.Printf("- Colors: %v\n", config.UseColors)
	}

	// Get geographical coordinates
	latitude, longitude, err := getCoordinates(zipCode)
	if err != nil {
		return fmt.Errorf("could not get coordinates: %w", err)
	}

	// Fetch and display weather information
	return fetchWeather(latitude, longitude, cmd.showDaily, cmd.showHourly, cmd.displayMode, unitSystem, useColors)
}

// Print detailed help information
func printHelp() {
	fmt.Printf("%s v%s - Command Line Weather Information\n\n", appName, appVersion)
	fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	fmt.Printf("  -help, -?           Show this help message\n")
	fmt.Printf("  -daily, -d          Show 7-day forecast\n")
	fmt.Printf("  -hourly, -h         Show hourly forecast for the next 24 hours\n")
	fmt.Printf("  -zip, -z [location] Override default location (ZIP code or city name)\n")
	fmt.Printf("  -table, -t          Display output in table format\n")
	fmt.Printf("  -text, -T           Display output in text format\n")
	fmt.Printf("  -units, -u [system] Use specific units (metric or imperial)\n")
	fmt.Printf("  -color, -c          Enable colored output\n")
	fmt.Printf("  -no-color, -nc      Disable colored output\n")
	fmt.Printf("  -save, -s           Save current settings as defaults\n\n")

	fmt.Printf("Examples:\n")
	fmt.Printf("  Basic usage (shows only current weather for default location):\n")
	fmt.Printf("    %s\n\n", os.Args[0])

	fmt.Printf("  Show 7-day forecast for a different location in imperial units:\n")
	fmt.Printf("    %s -daily -zip 10001 -units imperial\n\n", os.Args[0])

	fmt.Printf("  Show hourly forecast in table format with colors and save settings:\n")
	fmt.Printf("    %s -hourly -table -color -save\n\n", os.Args[0])

	fmt.Printf("  Show current weather in metric units without colors:\n")
	fmt.Printf("    %s -units metric -no-color\n\n", os.Args[0])

	fmt.Printf("  Save imperial as default unit system:\n")
	fmt.Printf("    %s -units imperial -save\n\n", os.Args[0])

	fmt.Printf("Configuration:\n")
	fmt.Printf("  Your preferences are stored in: %s\n", getConfigPath())
	fmt.Printf("  Weather data is cached for one hour in: %s\n", getCacheDir())
}

// Load configuration from file
func loadConfig() Config {
	configPath := getConfigPath()
	config := Config{
		DisplayMode: DisplayText, // Default to text mode
		Units:       UnitMetric,  // Default to metric
		UseColors:   true,        // Default to colors enabled
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config
	}

	// Try to unmarshal, ignore errors (will use default values)
	_ = json.Unmarshal(data, &config)
	return config
}

// Save configuration to file
func saveConfig(config Config) error {
	configPath := getConfigPath()

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// getZipCode returns the location to use for weather lookup
func getZipCode(override string, config *Config) string {
	if override != "" {
		// Don't save the new zip code to config
		return override
	}

	// Use config zip code if available
	if config.ZipCode != "" {
		return config.ZipCode
	}

	// Otherwise ask user
	fmt.Print("Enter your location (ZIP/postal code or city name): ")
	var zip string
	fmt.Scanln(&zip)

	// Don't save to config here

	return zip
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.weather_config"
	}
	return filepath.Join(home, ".weather_config", configFileName)
}

// GeoLocation represents a geographical point
type GeoLocation struct {
	Latitude  float64
	Longitude float64
	Name      string
	Country   string
}

// Use Open-Meteo's geocoding endpoint to convert location to lat/lon
func getCoordinates(location string) (float64, float64, error) {
	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", location)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	type GeoResponse struct {
		Results []struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Name      string  `json:"name"`
			Country   string  `json:"country"`
		} `json:"results"`
	}

	var geoResp GeoResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &geoResp)
	if err != nil || len(geoResp.Results) == 0 {
		return 0, 0, fmt.Errorf("location not found")
	}

	result := geoResp.Results[0]
	fmt.Printf("Location detected: %s, %s\n", result.Name, result.Country)
	return result.Latitude, result.Longitude, nil
}

// WeatherData structure to hold all weather information
type WeatherData struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		WindSpeed   float64 `json:"windspeed"`
		WeatherCode int     `json:"weathercode"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
	Daily struct {
		Time             []string  `json:"time"`
		WeatherCode      []int     `json:"weathercode"`
		TemperatureMax   []float64 `json:"temperature_2m_max"`
		TemperatureMin   []float64 `json:"temperature_2m_min"`
		PrecipitationSum []float64 `json:"precipitation_sum"`
	} `json:"daily"`
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature   []float64 `json:"temperature_2m"`
		Precipitation []float64 `json:"precipitation"`
		WeatherCode   []int     `json:"weathercode"`
	} `json:"hourly"`
}

// Cache file structure with timestamp and data
type CacheFile struct {
	Timestamp time.Time   `json:"timestamp"`
	Data      WeatherData `json:"data"`
}

// Fetch weather data from API or cache
func fetchWeather(lat, lon float64, showDaily, showHourly bool, displayMode DisplayMode, unitSystem UnitSystem, useColors bool) error {
	// Check cache first
	cacheKey := generateCacheKey(lat, lon, showDaily, showHourly, string(unitSystem))
	cachedData, cacheExists := checkCache(cacheKey)
	if cacheExists {
		fmt.Println("Using cached weather data")
		displayWeatherData(cachedData, showDaily, showHourly, displayMode, unitSystem, useColors)
		return nil
	}

	// Build URL with parameters for requested forecast types
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true", lat, lon)

	// Add unit-specific parameters
	if unitSystem == UnitImperial {
		url += "&temperature_unit=fahrenheit&windspeed_unit=mph&precipitation_unit=inch"
	}

	if showDaily {
		url += "&daily=weathercode,temperature_2m_max,temperature_2m_min,precipitation_sum"
	}

	if showHourly {
		url += "&hourly=temperature_2m,precipitation,weathercode&forecast_hours=24"
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	// Parse and save to cache
	var weather WeatherData
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return fmt.Errorf("could not parse weather data: %w", err)
	}

	// Save to cache
	if err := saveToCache(cacheKey, body); err != nil {
		// Non-critical error, just log it
		fmt.Fprintf(os.Stderr, "Warning: Failed to cache weather data: %v\n", err)
	}

	// Display the weather data
	displayWeatherData(weather, showDaily, showHourly, displayMode, unitSystem, useColors)
	return nil
}

// Generate a cache key from request parameters
func generateCacheKey(lat, lon float64, daily, hourly bool, unitSystem string) string {
	key := fmt.Sprintf("%.4f-%.4f-d%v-h%v-u%s", lat, lon, daily, hourly, unitSystem)
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// Get cache directory
func getCacheDir() string {
	cacheDir := filepath.Join(os.TempDir(), "weather-cache")
	os.MkdirAll(cacheDir, 0755)
	return cacheDir
}

// Check if a valid cache exists
func checkCache(cacheKey string) (WeatherData, bool) {
	cacheFile := filepath.Join(getCacheDir(), cacheKey+".json")

	// Check if file exists and is not too old
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return WeatherData{}, false
	}

	var cache CacheFile
	if err := json.Unmarshal(data, &cache); err != nil {
		return WeatherData{}, false
	}

	// Check if cache is still valid
	if time.Since(cache.Timestamp) > cacheDuration {
		return WeatherData{}, false
	}

	return cache.Data, true
}

// Save data to cache
func saveToCache(cacheKey string, data []byte) error {
	// First parse the data to verify it's valid
	var weather WeatherData
	if err := json.Unmarshal(data, &weather); err != nil {
		return err
	}

	cache := CacheFile{
		Timestamp: time.Now(),
		Data:      weather,
	}

	cacheData, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(getCacheDir(), cacheKey+".json")
	return os.WriteFile(cacheFile, cacheData, 0644)
}

// Display weather data in appropriate format
func displayWeatherData(weather WeatherData, showDaily, showHourly bool, mode DisplayMode, unitSystem UnitSystem, useColors bool) {
	switch mode {
	case DisplayTable:
		displayWeatherAsTable(weather, showDaily, showHourly, unitSystem, useColors)
	default:
		displayWeatherAsText(weather, showDaily, showHourly, unitSystem, useColors)
	}
}

// Get the appropriate temperature units based on unit system
func getTempUnit(unitSystem UnitSystem) string {
	if unitSystem == UnitImperial {
		return "°F"
	}
	return "°C"
}

// Get the appropriate wind speed units based on unit system
func getWindUnit(unitSystem UnitSystem) string {
	if unitSystem == UnitImperial {
		return "mph"
	}
	return "km/h"
}

// Get the appropriate precipitation units based on unit system
func getPrecipUnit(unitSystem UnitSystem) string {
	if unitSystem == UnitImperial {
		return "in"
	}
	return "mm"
}

// Helper function to get unit system name
func getUnitSystemName(unit UnitSystem) string {
	switch unit {
	case UnitMetric:
		return "Metric (°C, km/h, mm)"
	case UnitImperial:
		return "Imperial (°F, mph, in)"
	default:
		return string(unit)
	}
}

// colorizeTemp applies color to temperature based on its value
func colorizeTemp(temp float64, unitSystem UnitSystem) string {
	// Convert to Celsius for standard comparison if needed
	tempC := temp
	if unitSystem == UnitImperial {
		tempC = (temp - 32) * 5 / 9
	}

	// Color based on temperature ranges (in Celsius)
	var colorCode string
	switch {
	case tempC < -10:
		colorCode = colorBlue // Very cold
	case tempC < 0:
		colorCode = colorCyan // Cold
	case tempC < 15:
		colorCode = colorWhite // Cool
	case tempC < 25:
		colorCode = colorGreen // Pleasant
	case tempC < 30:
		colorCode = colorYellow // Warm
	case tempC < 35:
		colorCode = colorMagenta // Hot
	default:
		colorCode = colorRed // Very hot
	}

	// Format with units
	unit := getTempUnit(unitSystem)
	return fmt.Sprintf("%s%.1f%s%s", colorCode, temp, unit, colorReset)
}

// Text-based display format
func displayWeatherAsText(weather WeatherData, showDaily, showHourly bool, unitSystem UnitSystem, useColors bool) {
	tempUnit := getTempUnit(unitSystem)
	windUnit := getWindUnit(unitSystem)
	precipUnit := getPrecipUnit(unitSystem)

	fmt.Println("Current Weather:")
	if useColors {
		fmt.Printf("  Temperature: %s\n", colorizeTemp(weather.CurrentWeather.Temperature, unitSystem))
	} else {
		fmt.Printf("  Temperature: %.1f%s\n", weather.CurrentWeather.Temperature, tempUnit)
	}

	// Add high/low temperatures for today if daily data is available
	if len(weather.Daily.Time) > 0 {
		today := time.Now().Format("2006-01-02")
		for i, day := range weather.Daily.Time {
			if day == today {
				if useColors {
					fmt.Printf("  High/Low: %s/%s\n",
						colorizeTemp(weather.Daily.TemperatureMax[i], unitSystem),
						colorizeTemp(weather.Daily.TemperatureMin[i], unitSystem))
				} else {
					fmt.Printf("  High/Low: %.1f%s/%.1f%s\n",
						weather.Daily.TemperatureMax[i], tempUnit,
						weather.Daily.TemperatureMin[i], tempUnit)
				}
				break
			}
		}
	}

	fmt.Printf("  Wind Speed: %.1f %s\n", weather.CurrentWeather.WindSpeed, windUnit)
	fmt.Printf("  Time: %s\n", formatTime(weather.CurrentWeather.Time))
	fmt.Printf("  Weather: %s\n", getWeatherDescription(weather.CurrentWeather.WeatherCode))

	// Display daily forecast if requested
	if showDaily && len(weather.Daily.Time) > 0 {
		fmt.Println("\n7-Day Forecast:")
		for i, day := range weather.Daily.Time {
			t, _ := time.Parse("2006-01-02", day)

			if useColors {
				fmt.Printf("  %s: %s, %s to %s, Precipitation: %.1f%s\n",
					t.Format("Mon Jan 2"),
					getWeatherDescription(weather.Daily.WeatherCode[i]),
					colorizeTemp(weather.Daily.TemperatureMin[i], unitSystem),
					colorizeTemp(weather.Daily.TemperatureMax[i], unitSystem),
					weather.Daily.PrecipitationSum[i],
					precipUnit)
			} else {
				fmt.Printf("  %s: %s, %.1f%s to %.1f%s, Precipitation: %.1f%s\n",
					t.Format("Mon Jan 2"),
					getWeatherDescription(weather.Daily.WeatherCode[i]),
					weather.Daily.TemperatureMin[i], tempUnit,
					weather.Daily.TemperatureMax[i], tempUnit,
					weather.Daily.PrecipitationSum[i], precipUnit)
			}
		}
	}

	// Display hourly forecast if requested
	if showHourly && len(weather.Hourly.Time) > 0 {
		fmt.Println("\nHourly Forecast (next 24h):")
		for i := 0; i < 24 && i < len(weather.Hourly.Time); i++ {
			t, _ := time.Parse("2006-01-02T15:04", weather.Hourly.Time[i])

			if useColors {
				fmt.Printf("  %s: %s, %s, Precipitation: %.1f%s\n",
					t.Format("15:04"),
					getWeatherDescription(weather.Hourly.WeatherCode[i]),
					colorizeTemp(weather.Hourly.Temperature[i], unitSystem),
					weather.Hourly.Precipitation[i], precipUnit)
			} else {
				fmt.Printf("  %s: %s, %.1f%s, Precipitation: %.1f%s\n",
					t.Format("15:04"),
					getWeatherDescription(weather.Hourly.WeatherCode[i]),
					weather.Hourly.Temperature[i], tempUnit,
					weather.Hourly.Precipitation[i], precipUnit)
			}
		}
	}
}

// Table-based display format
func displayWeatherAsTable(weather WeatherData, showDaily, showHourly bool, unitSystem UnitSystem, useColors bool) {
	tempUnit := getTempUnit(unitSystem)
	windUnit := getWindUnit(unitSystem)
	precipUnit := getPrecipUnit(unitSystem)

	// Current weather display
	fmt.Println("Current Weather:")
	printLine(60) // Increased width to accommodate high/low
	fmt.Printf("| %-10s | %-12s | %-10s | %-12s | %-15s |\n", "Temperature", "High/Low", "Wind", "Time", "Condition")
	printLine(60)

	// Find today's high/low if available
	highTemp, lowTemp := weather.CurrentWeather.Temperature, weather.CurrentWeather.Temperature
	if len(weather.Daily.Time) > 0 {
		today := time.Now().Format("2006-01-02")
		for i, day := range weather.Daily.Time {
			if day == today {
				highTemp = weather.Daily.TemperatureMax[i]
				lowTemp = weather.Daily.TemperatureMin[i]
				break
			}
		}
	}

	if useColors {
		fmt.Printf("| %-10s | %-12s | %-10.1f %s | %-12s | %-15s |\n",
			colorizeTemp(weather.CurrentWeather.Temperature, unitSystem),
			fmt.Sprintf("%s/%s",
				colorizeTemp(highTemp, unitSystem),
				colorizeTemp(lowTemp, unitSystem)),
			weather.CurrentWeather.WindSpeed, windUnit,
			formatTime(weather.CurrentWeather.Time),
			truncateString(getWeatherDescription(weather.CurrentWeather.WeatherCode), 15))
	} else {
		fmt.Printf("| %-10.1f%s | %-12s | %-10.1f %s | %-12s | %-15s |\n",
			weather.CurrentWeather.Temperature, tempUnit,
			fmt.Sprintf("%.1f/%.1f%s", highTemp, lowTemp, tempUnit),
			weather.CurrentWeather.WindSpeed, windUnit,
			formatTime(weather.CurrentWeather.Time),
			truncateString(getWeatherDescription(weather.CurrentWeather.WeatherCode), 15))
	}
	printLine(60)

	// Display daily forecast if requested
	if showDaily && len(weather.Daily.Time) > 0 {
		fmt.Println("\n7-Day Forecast:")
		printLine(80)
		fmt.Printf("| %-10s | %-15s | %-12s | %-12s | %-15s |\n",
			"Date", "Condition", "Min Temp", "Max Temp", "Precipitation")
		printLine(80)

		for i, day := range weather.Daily.Time {
			t, _ := time.Parse("2006-01-02", day)

			if useColors {
				fmt.Printf("| %-10s | %-15s | %-12s | %-12s | %-15.1f%s |\n",
					t.Format("Mon Jan 2"),
					truncateString(getWeatherDescription(weather.Daily.WeatherCode[i]), 15),
					colorizeTemp(weather.Daily.TemperatureMin[i], unitSystem),
					colorizeTemp(weather.Daily.TemperatureMax[i], unitSystem),
					weather.Daily.PrecipitationSum[i], precipUnit)
			} else {
				fmt.Printf("| %-10s | %-15s | %-12.1f%s | %-12.1f%s | %-15.1f%s |\n",
					t.Format("Mon Jan 2"),
					truncateString(getWeatherDescription(weather.Daily.WeatherCode[i]), 15),
					weather.Daily.TemperatureMin[i], tempUnit,
					weather.Daily.TemperatureMax[i], tempUnit,
					weather.Daily.PrecipitationSum[i], precipUnit)
			}
		}
		printLine(80)
	}

	// Display hourly forecast if requested
	if showHourly && len(weather.Hourly.Time) > 0 {
		fmt.Println("\nHourly Forecast (next 24h):")
		printLine(60)
		fmt.Printf("| %-5s | %-15s | %-12s | %-15s |\n",
			"Time", "Condition", "Temperature", "Precipitation")
		printLine(60)

		for i := 0; i < 24 && i < len(weather.Hourly.Time); i++ {
			t, _ := time.Parse("2006-01-02T15:04", weather.Hourly.Time[i])

			if useColors {
				fmt.Printf("| %-5s | %-15s | %-12s | %-15.1f%s |\n",
					t.Format("15:04"),
					truncateString(getWeatherDescription(weather.Hourly.WeatherCode[i]), 15),
					colorizeTemp(weather.Hourly.Temperature[i], unitSystem),
					weather.Hourly.Precipitation[i], precipUnit)
			} else {
				fmt.Printf("| %-5s | %-15s | %-12.1f%s | %-15.1f%s |\n",
					t.Format("15:04"),
					truncateString(getWeatherDescription(weather.Hourly.WeatherCode[i]), 15),
					weather.Hourly.Temperature[i], tempUnit,
					weather.Hourly.Precipitation[i], precipUnit)
			}
		}
		printLine(60)
	}
}

// Helper function to print a horizontal line for tables
func printLine(width int) {
	fmt.Print("+")
	for i := 0; i < width-2; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

// Helper function to truncate strings to fit in table cells
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Helper function to format ISO time string
func formatTime(timeStr string) string {
	// Parse time and format it to a more readable form
	t, err := time.Parse("2006-01-02T15:04", timeStr)
	if err != nil {
		return timeStr
	}
	return t.Format("15:04")
}

// getWeatherDescription converts weather code to human-readable description
func getWeatherDescription(code int) string {
	descriptions := map[int]string{
		0:  "Clear sky",
		1:  "Mainly clear",
		2:  "Partly cloudy",
		3:  "Overcast",
		45: "Fog",
		48: "Depositing rime fog",
		51: "Light drizzle",
		53: "Moderate drizzle",
		55: "Dense drizzle",
		56: "Light freezing drizzle",
		57: "Dense freezing drizzle",
		61: "Slight rain",
		63: "Moderate rain",
		65: "Heavy rain",
		66: "Light freezing rain",
		67: "Heavy freezing rain",
		71: "Slight snow fall",
		73: "Moderate snow fall",
		75: "Heavy snow fall",
		77: "Snow grains",
		80: "Slight rain showers",
		81: "Moderate rain showers",
		82: "Violent rain showers",
		85: "Slight snow showers",
		86: "Heavy snow showers",
		95: "Thunderstorm",
		96: "Thunderstorm with slight hail",
		99: "Thunderstorm with heavy hail",
	}

	if desc, ok := descriptions[code]; ok {
		return desc
	}
	return "Unknown"
}
