# Release Guide

## Preconditions

1. Update `VERSION`, `CHANGELOG.md`, `PKGBUILD`, and `.SRCINFO` together.
2. Run `make check` and verify the archive checksum for the exact public tag.
3. Confirm the tag name is `v<VERSION>` and contains the intended source.

## GitHub Actions

The release workflow runs for `v*` tag pushes or manual dispatch. It verifies
that the tag matches `VERSION`, builds the binary, creates the GitHub release,
and updates the AUR package checksum from that exact tag archive.

Do not add credentials to workflow files. AUR authentication is supplied only
through the protected GitHub Actions secret required by the workflow.

## Recovery

If a release is wrong, stop publication, identify the affected tag and release,
and correct them through the provider interfaces. Do not rewrite history or
force-push as part of routine recovery.
