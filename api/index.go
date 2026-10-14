package handler

import (
	"net/http"

	"github.com/DosDeosos/dota-commu-tournament-backend/app"
)

var h = app.New()

func Handler(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
}
