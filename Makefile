MIN_COVERAGE  := 80
COVER_PROFILE := coverage.out
RACE          := $(if $(filter Windows_NT,$(OS)),,-race)

.PHONY: run test lint fmt codeguard coverage-gate vuln-check secret-scan quality-gate

run:
	go run ./cmd/dev

test:
	go test $(RACE) ./...

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

coverage-gate:
	go test -count=1 -coverprofile=$(COVER_PROFILE) ./app > /dev/null
	@go tool cover -func=$(COVER_PROFILE) | awk -v min=$(MIN_COVERAGE) '/^total:/ { sub("%", "", $$3); printf "coverage %s%% (min %s%%)\n", $$3, min; if ($$3 + 0 < min) exit 1 }'

vuln-check:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

secret-scan:
	gitleaks detect --source . --no-banner --redact

codeguard:
	go run ./cmd/codeguard

quality-gate:
	$(MAKE) codeguard
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) coverage-gate
	$(MAKE) vuln-check
	$(MAKE) secret-scan
