[![Release and AUR Publish](https://github.com/Streek/go-weather/actions/workflows/release.yml/badge.svg)](https://github.com/Streek/go-weather/actions/workflows/release.yml)

# Console Go Weather

A powerful command-line weather application written in Go that provides current, daily, and hourly weather forecasts.

![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)

## Features

- Get current weather conditions for any location
- View 7-day forecasts with daily highs and lows
- Check hourly forecasts for the next 24 hours
- Support for multiple display formats (text and table)
- Location caching and persistent user preferences
- Work with any location worldwide (postal code or city name)

## Installation

### Prerequisites

- Go 1.16 or higher

### Building from source

```bash
# Clone the repository
git clone https://github.com/streek/go-weather.git
cd go-weather

# Build the application
make build

# Install system-wide (requires root privileges)
sudo make install
```

### Using package managers

#### Arch Linux

```bash
# Install from AUR
yay -S go-weather
```

## Usage

```bash
# Basic usage (shows only current weather for default location)
go-weather

# Show 7-day forecast for a different location
go-weather -daily -zip 10001

# Show hourly forecast in table format
go-weather -hourly -table

# Show both daily and hourly forecasts for another location
go-weather -d -h -z "Paris, France"

# Save a new default location without running a query
go-weather --save-zip "Tokyo, Japan"

# Save table display format as default
go-weather --save-display table
```

### Command-line Options

- `-help`, `-?`: Show help information
- `-daily`, `-d`: Show 7-day forecast
- `-hourly`, `-h`: Show hourly forecast for the next 24 hours
- `-zip`, `-z` [location]: Override default location (ZIP code or city name)
- `-table`, `-t`: Display output in table format (save preference)
- `-text`, `-T`: Display output in text format (save preference)
- `-save-zip` [location]: Save a default location without querying weather
- `-save-display` [format]: Save display format preference (text or table)

## Configuration

The application stores your preferences in `~/.weather_config/weather_config.json`.
Weather data is cached for one hour in your system's temporary directory.

## Weather Data Source

This application uses the [Open-Meteo API](https://open-meteo.com/) for weather data, which is a free and open-source weather API.

## Contributing

Contributions to go-weather are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please note that any contributions you make will be under the GNU GPL v3 license.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Open-Meteo](https://open-meteo.com/) for their free weather API
- All the awesome contributors to this project
