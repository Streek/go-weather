# Repository Guidance

## Scope

This repository contains the Go command-line weather client. Preserve its CLI
interface and text/table output; do not add web or footer UI to this project.

## Safe Changes

- Keep credentials, tokens, private keys, account IDs, internal hostnames, and
  infrastructure identifiers out of source, documentation, commits, and logs.
- Do not commit or push from an update workflow unless explicitly requested.
- Preserve existing untracked source snapshots, archives, and package artifacts.
- Use the public Open-Meteo endpoints already defined in `main.go`.

## Validation

- `make check` runs formatting, `go vet`, and tests.
- `make build` creates the local `go-weather` binary.
- `make package` requires Arch Linux `makepkg` and is not part of `make check`.

The application requires Go 1.21 or newer, as declared by `go.mod` and the
AUR package metadata.

## Release

Releases are tag-driven. The tag must be `v<VERSION>` and match `VERSION`.
See [RELEASE.md](RELEASE.md) for the release procedure and [SECURITY.md](SECURITY.md)
for reporting and secret-handling rules.
