# Dota Commu. Tournament — Backend

A small read-only Go API that serves player ratings to the
[frontend](https://github.com/DosDeosos/Dota-Commu.-Tournament). It runs as
one Vercel Go function and needs no external database.

The data lives in this repo as SQL (`data/data.sql`). On startup the function
loads it into an in-memory SQLite database (pure Go, no cgo), validates every
row, and serves it as JSON. To change the data, edit the file and push.
Vercel redeploys with the new data.

The response has the same shape as `PLAYERS` and `RATINGS` in the frontend's
`src/data/players.ts`. Static config (axes, traits, labels) and the score
math stay in the frontend.

## Layout

```
data/schema.sql   SQLite tables
data/data.sql     the tournament data — edit this
api/index.go      Vercel entrypoint (vercel.json rewrites every path here)
app/              HTTP handlers, SQLite loading, validation
cmd/dev/          local dev server
```

## API

| Method | Path          | Returns                                      |
| ------ | ------------- | -------------------------------------------- |
| GET    | `/api/data`   | `{ players: Player[], ratings: Rating[] }`   |
| GET    | `/api/health` | `{"ok":true,"data":"loaded"}`, or 503 if the data is invalid |

`/api/data` is cached at Vercel's edge for 60 seconds, so a new deploy can
take up to a minute to show up.

## Editing the data

1. Edit `data/data.sql`, using the keys and ranges below.
2. Run `make test` (or `go test ./...`). It loads the file and fails on any
   mistake: a score out of range, a misspelled key or field such as `"socre"`,
   broken JSON, a rating for a player that doesn't exist, or a duplicate
   rating from the same rater.
3. Commit and push. CI runs the same checks, and Vercel redeploys.

Valid keys and ranges (checked on load by `app/model.go`):

| Field | Keys | Range |
| --- | --- | --- |
| `primary_roles`, `positions` | `carry`, `mid`, `offlane`, `hard`, `soft` | positions 1–20 |
| `scores` (Game Skill) | `skill`, `knowledge`, `consistency`, `teamplay`, `laning`, `lategame` | 1–10 |
| `scores` (Soft Skill, objective) | `captaincy`, `communication`, `roleFlexibility`, `heroPoolDepth` | 1–5 |
| `traits` | `summoner`, `auraMan`, `initiator`, `dogpiler`, `cardio`, `pest`, `firstMover`, `secondWind`, `realSupport`, `meCarry`, `madman`, `safetyFirst`, `speedrunner`, `slowpoke`, `yesMan`, `commander`, `farmGod`, `headfirst`, `chameleon` | — |

Each score is `{"score": n}` with an optional `"comment"`. Players appear in the
order they are inserted. Each `(player, rater)` pair can have only one rating.
If you add an axis, role or trait in the frontend's `players.ts`, also add it
to `app/model.go` and to this table.

## Local dev

```sh
make run    # http://localhost:8080/api/data
```

`go.mod` pins Go 1.27.2. Vercel builds with exactly that version, so bump it
when a Go patch release fixes something `make vuln-check` reports. A newer
local Go downloads the pinned toolchain automatically.

## Quality gate

`make quality-gate` runs the same checks as CI (`.github/workflows/ci.yml`).

| Target          | What it checks                                                                 |
| --------------- | ------------------------------------------------------------------------------ |
| `codeguard`     | no code comments (only `BREAKING:`, `LICENSE:`, `NOSONAR`, `//go:` directives), files ≤ 600 lines, no empty or duplicate functions |
| `lint`          | golangci-lint with `.golangci.yml`, which covers SonarQube for IDE's Sonar way rules for Go (cognitive complexity 15, shadowed builtins, duplicate branches, …) plus gosec, errcheck and others |
| `test`          | `go test -race`, including validation of `data/data.sql`. Race detection is skipped on Windows because it needs cgo |
| `coverage-gate` | at least 80% coverage of `app/`                                               |
| `vuln-check`    | govulncheck on dependencies and the Go standard library                       |
| `secret-scan`   | gitleaks over the git history                                                  |

`make fmt` applies gofmt and goimports.

## Deploy to Vercel

1. Import this repo into Vercel as a **new project**, separate from the
   frontend. Framework preset: **Other**.
2. Optional: under **Settings → Environment Variables**, set `CORS_ORIGIN` to
   the frontend URL, e.g. `https://<frontend>.vercel.app`. It defaults to `*`.
3. Deploy, then open `/api/health`. It should return
   `{"ok":true,"data":"loaded"}`.

## Using it from the frontend

Add `NEXT_PUBLIC_API_URL=https://<backend>.vercel.app` to the frontend's env,
then fetch `/api/data` and feed `players` and `ratings` into the existing
aggregation logic in place of the hardcoded `PLAYERS` and `RATINGS` constants.
