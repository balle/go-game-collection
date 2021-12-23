package controllers

import (
	"net/http"

	"github.com/balle/go-game-collection/models"
	"github.com/balle/go-game-collection/store"
	"github.com/balle/go-game-collection/utils"
)

// Add a new game system to the db
func AddGameSystem(w http.ResponseWriter, r *http.Request) {
	// Parse Form data
	err := r.ParseForm()

	if utils.GotError(w, err) {
		return
	}

	params := r.Form

	// Create new game system struct and save it in the database
	store.Db.Create(&models.GameSystem{
		Name: params.Get("name"),
	})

	ListGames(w, r)
}
