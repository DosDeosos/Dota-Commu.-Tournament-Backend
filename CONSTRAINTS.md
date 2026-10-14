# Constraints

Last reviewed: 2026-10-14.

## Floor (always enforced)

- No comments in code except `BREAKING:`, `LICENSE:`, `NOSONAR` and `//go:` directives
- No `//nolint` or other suppression
- No skipped or deleted tests without a reason in the commit message
- No secrets in source
- This file does not get weakened to make a change pass

## Enforced with numbers

| Dimension       | Rule                                                                                         | Checked by           | Runs at   |
| --------------- | -------------------------------------------------------------------------------------------- | -------------------- | --------- |
| Build           | Compiles, `go vet` clean                                                                     | `make test`          | local, CI |
| Data            | `data/data.sql` loads and every row passes validation                                        | `make test`          | local, CI |
| Unit tests      | All pass                                                                                     | `make test`          | local, CI |
| Race            | `go test -race` clean                                                                        | `make test` (Linux)  | CI        |
| Complexity      | `gocognit` ≤ 15, `gocyclo` ≤ 15, `funlen` ≤ 80 lines / 50 statements, `nestif` ≤ 4           | `make lint`          | local, CI |
| Lint            | Zero golangci-lint issues (gosec, staticcheck, errcheck, revive, …), incl. Sonar way for Go  | `make lint`          | local, CI |
| File size       | Every Go file ≤ 600 lines                                                                    | `make codeguard`     | local, CI |
| Comments        | Only the allowed prefixes                                                                    | `make codeguard`     | local, CI |
| Functions       | No empty function outside tests (Sonar S1186); no two identical functions in a file (S4144)  | `make codeguard`     | local, CI |
| Coverage        | `app/` ≥ 80%                                                                                 | `make coverage-gate` | local, CI |
| Vulnerabilities | No reachable known vulnerability                                                             | `make vuln-check`    | local, CI |
| Secrets         | No secrets in the tree                                                                       | `make secret-scan`   | local, CI |

## Exceptions

| ID  | Rule                                           | Path        | Reason                                                                                                   | Expires |
| --- | ---------------------------------------------- | ----------- | -------------------------------------------------------------------------------------------------------- | ------- |
| E1  | funlen, dupl, goconst, gosec, noctx, bodyclose | `_test\.go` | Table-driven tests repeat literals and setup and use plain `httptest` requests; the code under test is linted in full | never   |
| E2  | gosec G104                                     | all         | Unchecked errors are already reported by `errcheck`; G104 would duplicate it                             | never   |
