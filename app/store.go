package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/DosDeosos/dota-commu-tournament-backend/data"

	_ "modernc.org/sqlite"
)

type snapshot struct {
	Players []Player `json:"players"`
	Ratings []Rating `json:"ratings"`
}

var loadData = sync.OnceValues(func() (snapshot, error) {
	return load(context.Background(), data.Schema+data.Rows)
})

func load(ctx context.Context, sqlText string) (snapshot, error) {
	db, err := sql.Open("sqlite", "file::memory:?_pragma=foreign_keys(1)")
	if err != nil {
		return snapshot{}, err
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	if _, err := db.ExecContext(ctx, sqlText); err != nil {
		return snapshot{}, fmt.Errorf("data.sql: %w", err)
	}
	players, err := loadPlayers(ctx, db)
	if err != nil {
		return snapshot{}, err
	}
	ratings, err := loadRatings(ctx, db)
	if err != nil {
		return snapshot{}, err
	}
	return snapshot{Players: players, Ratings: ratings}, nil
}

func loadPlayers(ctx context.Context, db *sql.DB) ([]Player, error) {
	rows, err := db.QueryContext(ctx, `select id, name, primary_roles from players order by rowid`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []Player{}
	for rows.Next() {
		var p Player
		var roles string
		if err := rows.Scan(&p.ID, &p.Name, &roles); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(roles), &p.PrimaryRoles); err != nil {
			return nil, fmt.Errorf("player %q primary_roles: %w", p.ID, err)
		}
		if err := p.validate(); err != nil {
			return nil, fmt.Errorf("player %q: %w", p.ID, err)
		}
		players = append(players, p)
	}
	return players, rows.Err()
}

func loadRatings(ctx context.Context, db *sql.DB) ([]Rating, error) {
	rows, err := db.QueryContext(ctx, `select subject_id, rater, scores, positions, traits from ratings order by subject_id, rowid`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ratings := []Rating{}
	for rows.Next() {
		var r Rating
		var scores, positions, traits string
		if err := rows.Scan(&r.SubjectID, &r.Rater, &scores, &positions, &traits); err != nil {
			return nil, err
		}
		if err := unmarshalRating(&r, scores, positions, traits); err != nil {
			return nil, fmt.Errorf("rating %q by %q: %w", r.SubjectID, r.Rater, err)
		}
		ratings = append(ratings, r)
	}
	return ratings, rows.Err()
}

func unmarshalRating(r *Rating, scores, positions, traits string) error {
	if err := strictUnmarshal(scores, &r.Scores); err != nil {
		return fmt.Errorf("scores: %w", err)
	}
	if err := strictUnmarshal(positions, &r.Positions); err != nil {
		return fmt.Errorf("positions: %w", err)
	}
	if err := json.Unmarshal([]byte(traits), &r.Traits); err != nil {
		return fmt.Errorf("traits: %w", err)
	}
	return r.validate()
}

func strictUnmarshal(text string, v any) error {
	dec := json.NewDecoder(strings.NewReader(text))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
