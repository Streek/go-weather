# Security Policy

## Reporting

Report suspected vulnerabilities privately to the repository owner through a
trusted GitHub security channel. Do not include credentials or private keys in
issues, pull requests, commits, or chat.

## Secret Handling

This CLI uses public Open-Meteo endpoints and has no repository-stored service
credentials. Keep local configuration and environment files outside version
control. If a secret is exposed, stop publication, rotate it at the provider,
remove it from reachable history through the appropriate repository process,
and verify the replacement before resuming.

## Runtime Notes

Network responses are bounded and structurally validated before display. User
configuration and cached weather data are stored with owner-only permissions
when the operating system permits it.
