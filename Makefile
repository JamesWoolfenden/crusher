.PHONY: build install test lint fmt vet staticcheck gosec govulncheck complexity check-all release-check sbom bump psbump update
BINARY=crusher
EXE=$(if $(filter Windows_NT,$(OS)),.exe,)

default: build

build:
	go build -o ${BINARY}${EXE}

install:
	go install .

test:
	go test ./...

lint:
	golangci-lint run --fix

fmt:
	gofumpt -l -w .

vet:
	go vet ./...

staticcheck: ## Run staticcheck
	@./scripts/run-staticcheck.sh

gosec: ## Run security scanner
	gosec -quiet -exclude-generated ./...

govulncheck: ## Check for known vulnerabilities
	govulncheck ./...

complexity: ## Check cyclomatic complexity
	gocyclo -over 15 -avg .

check-all: vet staticcheck gosec govulncheck ## Run all checks
	@echo "All checks passed!"

sbom:
	syft . -o cyclonedx-json=crusher.cdx.json

release-check:
	goreleaser check
	GITHUB_TOKEN= GITLAB_TOKEN= GITEA_TOKEN= goreleaser build --snapshot --clean

bump: release-check
	git push
	$(eval VERSION=$(shell git describe --tags --abbrev=0 | awk -F. '{OFS="."; $$NF+=1; print $0}'))
	git tag -a $(VERSION) -m "new release"
	git push origin $(VERSION)

psbump:
	git push
	powershell -command "./bump.ps1"

update:
	go get -u
	go mod tidy
	pre-commit autoupdate
