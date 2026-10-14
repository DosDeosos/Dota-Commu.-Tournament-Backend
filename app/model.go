package app

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var roles = set("carry", "mid", "offlane", "hard", "soft")

var axisMax = map[string]float64{
	"skill":           10,
	"knowledge":       10,
	"consistency":     10,
	"teamplay":        10,
	"laning":          10,
	"lategame":        10,
	"captaincy":       5,
	"communication":   5,
	"roleFlexibility": 5,
	"heroPoolDepth":   5,
}

const positionMax = 20

var traits = set(
	"summoner", "auraMan", "initiator", "dogpiler", "cardio", "pest",
	"firstMover", "secondWind", "realSupport", "meCarry", "madman",
	"safetyFirst", "speedrunner", "slowpoke", "yesMan", "commander",
	"farmGod", "headfirst", "chameleon",
)

const (
	maxNameLen    = 100
	maxCommentLen = 2000
)

var idPattern = regexp.MustCompile(`^[a-z0-9_-]{1,64}$`)

type Player struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	PrimaryRoles []string `json:"primaryRoles"`
}

type Score struct {
	Score   float64 `json:"score"`
	Comment string  `json:"comment,omitempty"`
}

type Rating struct {
	SubjectID string           `json:"subjectId"`
	Rater     string           `json:"rater"`
	Scores    map[string]Score `json:"scores"`
	Positions map[string]Score `json:"positions"`
	Traits    []string         `json:"traits"`
}

func (p *Player) validate() error {
	p.Name = strings.TrimSpace(p.Name)
	if !idPattern.MatchString(p.ID) {
		return fmt.Errorf("id must be 1-64 chars of a-z, 0-9, - or _")
	}
	if p.Name == "" || utf8.RuneCountInString(p.Name) > maxNameLen {
		return fmt.Errorf("name must be 1-%d characters", maxNameLen)
	}
	var err error
	if p.PrimaryRoles, err = checkKeys("primaryRoles", p.PrimaryRoles, roles); err != nil {
		return err
	}
	return nil
}

func (r *Rating) validate() error {
	r.Rater = strings.TrimSpace(r.Rater)
	if !idPattern.MatchString(r.SubjectID) {
		return fmt.Errorf("invalid subjectId")
	}
	if r.Rater == "" || utf8.RuneCountInString(r.Rater) > maxNameLen {
		return fmt.Errorf("rater must be 1-%d characters", maxNameLen)
	}
	var err error
	if r.Scores, err = checkScores("scores", r.Scores, func(key string) (float64, bool) {
		limit, ok := axisMax[key]
		return limit, ok
	}); err != nil {
		return err
	}
	if r.Positions, err = checkScores("positions", r.Positions, func(key string) (float64, bool) {
		return positionMax, roles[key]
	}); err != nil {
		return err
	}
	if r.Traits, err = checkKeys("traits", r.Traits, traits); err != nil {
		return err
	}
	return nil
}

func checkScores(field string, scores map[string]Score, limitFor func(string) (float64, bool)) (map[string]Score, error) {
	if scores == nil {
		return map[string]Score{}, nil
	}
	for key, s := range scores {
		limit, ok := limitFor(key)
		if !ok {
			return nil, fmt.Errorf("unknown %s key %q", field, key)
		}
		if err := s.validate(field+"."+key, limit); err != nil {
			return nil, err
		}
	}
	return scores, nil
}

func (s Score) validate(field string, limit float64) error {
	if s.Score < 1 || s.Score > limit {
		return fmt.Errorf("%s must be between 1 and %g", field, limit)
	}
	if utf8.RuneCountInString(s.Comment) > maxCommentLen {
		return fmt.Errorf("%s comment is longer than %d characters", field, maxCommentLen)
	}
	return nil
}

func checkKeys(field string, values []string, allowed map[string]bool) ([]string, error) {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, v := range values {
		if !allowed[v] {
			return nil, fmt.Errorf("unknown %s value %q", field, v)
		}
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out, nil
}

func set(values ...string) map[string]bool {
	m := make(map[string]bool, len(values))
	for _, v := range values {
		m[v] = true
	}
	return m
}
