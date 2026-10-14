# Dota-Commu.-Tournament-Backend

The read-only Go API behind the Dota 2 community player-rating site (sibling repo `Dota-Commu.-Tournament`). It
serves every player and rating at `GET /api/data`, in the same shape as `PLAYERS` and `RATINGS` in the frontend's
`src/data/players.ts`. It runs as one Vercel Go function and needs no external database.

Read [CONSTRAINTS.md](./CONSTRAINTS.md) before writing code. Do not weaken it to make a change pass. API, data
editing and deploy steps for people: [README.md](./README.md).

---

## Stack — do not substitute without discussion

| Concern  | Choice                                                                  | Where                   |
| -------- | ----------------------------------------------------------------------- | ----------------------- |
| Language | Go 1.27.2, pinned in `go.mod` (Vercel builds with exactly that version) | `go.mod`                |
| HTTP     | standard library `net/http` `ServeMux` with method patterns             | `app/app.go`            |
| Storage  | SQLite in memory via `modernc.org/sqlite` (pure Go, no cgo)             | `app/store.go`          |
| Data     | SQL text embedded with `go:embed`                                       | `data/`                 |
| Hosting  | Vercel Go runtime, one function; `vercel.json` rewrites every path      | `api/index.go`          |
| Lint     | golangci-lint v2, Sonar way rules for Go                                | `.golangci.yml`         |

## Structure

```
api/index.go        Vercel entrypoint (package handler)
app/                app.go (routes, CORS), store.go (SQLite load), model.go (types, validation)
data/               schema.sql, data.sql (the tournament data), embed.go
cmd/dev/            local dev server on $PORT (default 8080)
cmd/codeguard/      comment, file-length, empty-function and duplicate-function guard
```

Keep the dev server out of `main.go`, `cmd/api/` and `cmd/server/`: Vercel treats those paths as a whole-server Go
app instead of the `api/` function.

## How it works

- On the first request, `loadData` runs `data.Schema + data.Rows` in a fresh in-memory SQLite database, reads every
  row back, validates it and keeps the result in memory. The data only changes on redeploy, so later requests never
  touch SQLite.
- The database is opened with one connection: every connection to `:memory:` is a separate, empty database.
- Foreign keys are on (`_pragma=foreign_keys(1)`), and JSON columns have `check (json_valid(...))`, so a broken edit
  fails on load.
- Scores and positions are decoded with `DisallowUnknownFields`, so a typo such as `"socre"` fails instead of
  silently becoming a score of 0, which validation then rejects anyway.
- Role, axis and trait keys and score limits in `app/model.go` mirror the frontend's `players.ts`. Change both
  together, and update the key table in the README.
- `/api/data` is cached at Vercel's edge (`s-maxage=60`). `CORS_ORIGIN` limits browser origins (default `*`).

## Code quality

- **No comments in code.** `make codeguard` fails on any comment. Exceptions: `// BREAKING: …`, `// LICENSE: …`,
  `// NOSONAR`, and `//go:` directives. Explanations belong in this file, the README, commit messages and PRs.
- **Cognitive complexity of 15 or less per function.** golangci-lint runs `gocognit` and `gocyclo` at 15, `funlen` at
  80 lines / 50 statements, `nestif` at 4. Split into named helpers and return early.
- **Files of 600 lines or less.** `make codeguard` fails above 600.
- **No lint suppression.** `//nolint` is a comment and fails `codeguard`. If a linter is wrong for a category of
  files, change `.golangci.yml` with a path rule and a reason in the PR.
- **Standard library first.** The only dependency is the SQLite driver.
- **No magic numbers, no dead code, no speculative abstraction.** No interface without a second implementation.
- **Every validation rule has a test**, including a bad-data case in `TestLoadRejectsBadData`.
- Files are `snake_case.go`; tests sit next to their source as `<file>_test.go`.

## Agent skills

Team skills live in `.claude/skills/`, vendored and pinned (sources in `.claude/skills/README.md`).

- **`ponytail` applies to every change.** Choose the laziest solution that is actually correct: reuse before writing,
  stdlib before dependencies, delete before adding.
- **Re-check every diff with `code-review-and-quality` and `code-simplification` before committing.**
- Use `debugging-and-error-recovery` for failures and `test-driven-development` for new logic.
- The rules in this file and `CONSTRAINTS.md` win when a skill disagrees (no `ponytail:` comments; put the note in the
  PR).
- Never edit the vendored skill files.

## Quality gate

```bash
make test          # all tests, including validation of data/data.sql
make quality-gate  # codeguard + lint + test + coverage + govulncheck + gitleaks (CI runs this)
make fmt           # gofmt + goimports
```

GitHub Actions (`.github/workflows/ci.yml`) runs `make quality-gate` on every pull request and every push to `main`.

## Branch flow and commits

- `main` is production: every push deploys to Vercel. Run `make quality-gate` before pushing.
- Data edits can go straight to `main` once `make test` passes. Code changes go on a `feature/*` branch and a PR.
- Commit messages: `type: Summary.` (`feat`, `fix`, `refactor`, `docs`, `chore`, `ci`, `data`). One concern per
  commit.
- No AI attribution in commits or PRs.

### PR checklist

- [ ] `make quality-gate` green locally
- [ ] Keys in `app/model.go` still match the frontend's `players.ts` and the README table
- [ ] No comment, `//nolint` or secret
- [ ] README updated when the API or the data format changed

## What to avoid

- Writing to disk or expecting writes to persist: Vercel functions have no lasting filesystem. Data changes go
  through `data/data.sql` and a redeploy.
- Adding a dependency for what the standard library already does.
- A cgo SQLite driver: Vercel and Windows builds have no C compiler.
- Growing a file past 600 lines or a function past complexity 15 "just this once".
