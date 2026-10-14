package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DosDeosos/dota-commu-tournament-backend/data"
)

func TestRatingValidateNormalizes(t *testing.T) {
	r := Rating{
		SubjectID: "ratlab",
		Rater:     " captain ",
		Scores:    map[string]Score{"skill": {Score: 10}, "captaincy": {Score: 5}},
		Traits:    []string{"pest", "pest"},
	}
	if err := r.validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Rater != "captain" || len(r.Traits) != 1 || r.Positions == nil {
		t.Fatalf("not normalized: %+v", r)
	}
}

func TestRatingValidateRejects(t *testing.T) {
	cases := []struct {
		name    string
		rating  Rating
		wantErr string
	}{
		{"game axis over 10", Rating{SubjectID: "ratlab", Rater: "a", Scores: map[string]Score{"skill": {Score: 11}}}, "scores.skill"},
		{"soft axis over 5", Rating{SubjectID: "ratlab", Rater: "a", Scores: map[string]Score{"communication": {Score: 6}}}, "scores.communication"},
		{"below 1", Rating{SubjectID: "ratlab", Rater: "a", Scores: map[string]Score{"skill": {Score: 0}}}, "between 1"},
		{"position over 20", Rating{SubjectID: "ratlab", Rater: "a", Positions: map[string]Score{"mid": {Score: 21}}}, "positions.mid"},
		{"unknown axis", Rating{SubjectID: "ratlab", Rater: "a", Scores: map[string]Score{"aim": {Score: 1}}}, "unknown scores key"},
		{"unknown trait", Rating{SubjectID: "ratlab", Rater: "a", Traits: []string{"tryhard"}}, "unknown traits"},
		{"blank rater", Rating{SubjectID: "ratlab", Rater: "  "}, "rater"},
		{"unknown position", Rating{SubjectID: "ratlab", Rater: "a", Positions: map[string]Score{"jungle": {Score: 1}}}, "unknown positions key"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.rating
			err := r.validate()
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestPlayerValidate(t *testing.T) {
	ok := Player{ID: "ratlab", Name: "หนูทดลอง", PrimaryRoles: []string{"soft", "hard"}}
	if err := ok.validate(); err != nil {
		t.Fatal(err)
	}
	for _, p := range []Player{
		{ID: "Bad ID", Name: "x"},
		{ID: "ok", Name: ""},
		{ID: "ok", Name: "x", PrimaryRoles: []string{"jungle"}},
	} {
		if err := p.validate(); err == nil {
			t.Errorf("expected error for %+v", p)
		}
	}
}

func TestAPI(t *testing.T) {
	h := New()
	do := func(method, path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
		return rec
	}

	rec := do("GET", "/api/data")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/data: got %d: %s", rec.Code, rec.Body)
	}
	var snap snapshot
	if err := json.Unmarshal(rec.Body.Bytes(), &snap); err != nil {
		t.Fatal(err)
	}
	if len(snap.Players) == 0 || len(snap.Ratings) == 0 {
		t.Fatalf("expected players and ratings, got %+v", snap)
	}

	if rec := do("GET", "/api/health"); rec.Code != http.StatusOK {
		t.Errorf("health: got %d", rec.Code)
	}
	for _, method := range []string{"POST", "PUT", "DELETE"} {
		if rec := do(method, "/api/data"); rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s /api/data: got %d, want 405", method, rec.Code)
		}
	}
	rec = do("OPTIONS", "/api/data")
	if rec.Code != http.StatusNoContent || rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("preflight: got %d %q", rec.Code, rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestDataFile(t *testing.T) {
	snap, err := loadData()
	if err != nil {
		t.Fatalf("data/data.sql is invalid: %v", err)
	}
	ids := map[string]bool{}
	for _, p := range snap.Players {
		ids[p.ID] = true
	}
	for _, r := range snap.Ratings {
		if !ids[r.SubjectID] {
			t.Errorf("rating by %q for unknown player %q", r.Rater, r.SubjectID)
		}
	}
}

func TestLoadRejectsBadData(t *testing.T) {
	const player = `insert into players (id, name, primary_roles) values ('p', 'P', '["mid"]');`
	rating := func(scores string) string {
		return player + `insert into ratings (subject_id, rater, scores) values ('p', 'a', '` + scores + `');`
	}
	cases := []struct {
		name, rows, wantErr string
	}{
		{"score too high", rating(`{"skill": {"score": 11}}`), "scores.skill"},
		{"typo in field", rating(`{"skill": {"socre": 5}}`), "unknown field"},
		{"broken JSON", rating(`{"skill": `), "CHECK constraint"},
		{"unknown player", `insert into ratings (subject_id, rater) values ('ghost', 'a');`, "FOREIGN KEY"},
		{"unknown role", `insert into players (id, name, primary_roles) values ('p', 'P', '["jungle"]');`, "primaryRoles"},
		{"duplicate rating", rating(`{}`) + `insert into ratings (subject_id, rater) values ('p', 'a');`, "UNIQUE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := load(t.Context(), data.Schema+tc.rows)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, tc.wantErr)
			}
		})
	}
}
