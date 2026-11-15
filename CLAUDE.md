# crusher - Claude guidance

## What this project is

`crusher` is a Google Cloud BigTable maintenance CLI. `crusher clip` compacts a table by
deleting rows older than `--days` (default 180) whose row key matches `--keyfilter`
(default `.*`). Sibling of ghat, holden and picard, same conventions (urfave/cli v2,
zerolog, `src/`, `src/version`).

Module path: `github.com/JamesWoolfenden/crusher`

## Layout

- `main.go` - the CLI: `clip` (alias `c`) and `version` (alias `v`); `validateInput`
  checks project/instance/table, `confirmDeletion` is the prompt skipped by `--yes`
- `src/crusher.go` - the `Crusher` struct (table, project, instance, key filter, days,
  dry-run, bigtable filter)
- `src/query.go` - `ReadWithFilter` finds the rows to delete, `DeleteRows` removes them,
  `Clip` ties the two together and returns the number of rows
- `src/version` - `Version`, set at build time by goreleaser via `-ldflags`
- `scripts/run-staticcheck.sh` - staticcheck wrapper used by `make staticcheck`

## Safety

This tool **deletes data**. Never run `crusher clip` against a real instance while testing:
use `--dry-run`, or the BigTable emulator. Tests must not need GCP credentials.

## Build and check

- `make build` / `make install` / `make test`
- `make check-all` - vet, staticcheck, gosec, govulncheck
- `make release-check` - `goreleaser check` plus a snapshot build
- `make bump` - tag the next patch version and push it; the tag triggers `release.yml`
  (goreleaser, Docker Hub, Homebrew cask, scoop dispatch)
- `pre-commit` runs gofmt, go vet, go test and go mod tidy; never bypass it

## Conventions

- Dependency vulnerabilities: use `picard` (`picard scan -d .`, `picard dependabot -d . --prs`).
  CI runs picard on every push to main; Dependabot only opens security PRs
  (`open-pull-requests-limit: 0`).
- Pin GitHub Actions to a commit SHA with a version comment (`ghat all` does this).
- The Go version comes from `go.mod` everywhere (CI, release, Docker); raise it there only.
- Files in this repo mix line endings: edit byte-for-byte, and check `git diff --stat`
  before committing.
