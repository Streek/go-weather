# Engineering

The CLI is a Go 1.21 module. `make check` is the documented validation command and runs formatting checks, `go vet`, and tests. `make build` creates the local binary, while `make package` is an Arch Linux packaging path outside the normal check.

HTTP requests use a 15-second timeout and cap response bodies at 1 MiB. Configuration and cache writes use restrictive permissions and atomic replacement.
