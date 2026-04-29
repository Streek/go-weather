# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Location input now properly URL-encodes query parameters, so multi-word city names like "New York" or "San Francisco" actually work instead of producing malformed API requests
- Interactive location prompt now reads the full line of input; previously `fmt.Scanln` would silently truncate anything after the first space
- All HTTP requests now use a client with a 15-second timeout instead of hanging indefinitely on slow or unresponsive networks
- Geocoding and weather API calls now check HTTP status codes and surface meaningful errors instead of trying to parse error responses as valid JSON
- `io.ReadAll` error in the geocoding path is no longer silently discarded
- Table display mode now correctly aligns columns when ANSI color codes are present; previously the invisible escape sequences threw off all column widths
- Cache directory creation (`MkdirAll`) now logs a warning on failure instead of silently discarding the error

### Changed

- `-h` now maps to help (matching standard CLI convention); hourly forecast short flag is now `-H`
- PKGBUILD `pkgver` is now a static value instead of a dynamic `$(cat VERSION)` subshell that some AUR tooling couldn't evaluate
- PKGBUILD `makedepends` now specifies `go>=1.21` so users with older Go get a clear dependency error instead of a cryptic build failure
- `go.mod` minimum Go version updated from 1.16 to 1.21 to match the CI build environment
- Release workflow SSH config for AUR no longer uses a heredoc inside indented YAML (the leading whitespace was being baked into the SSH config, breaking host matching and causing auth failures on `git push`)

## [1.0.1] - YYYY-MM-DD

### Added

- Initial release
- Current weather conditions display
- 7-day forecast option
- Hourly forecast option
- Table and text display modes
- Metric and imperial units support
- Color-coded temperature output
- Location saving and preferences
- Weather data caching
