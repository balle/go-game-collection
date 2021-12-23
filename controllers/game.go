package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"text/template"

	"github.com/balle/go-game-collection/models"
	"github.com/balle/go-game-collection/store"
	"github.com/balle/go-game-collection/utils"
	"github.com/gorilla/mux"
)

// Receive all games from the db and call the game list template
func ListGames(w http.ResponseWriter, r *http.Request) {
	var games []models.Game
	var gameSystems []models.GameSystem
	t, err := template.ParseFiles("templates/game/list.html")

	if utils.GotError(w, err) {
		return
	}

	// Fetch all games and game systems from db
	store.Db.Find(&games)
	store.Db.Find(&gameSystems)

	// Sort them by name
	sort.Slice(games, func(i, j int) bool {
		return sort.StringsAreSorted([]string{games[i].Name, games[j].Name})
	})

	sort.Slice(gameSystems, func(i, j int) bool {
		return sort.StringsAreSorted([]string{gameSystems[i].Name, gameSystems[j].Name})
	})

	// Compile template with our data
	data := struct {
		Games       []models.Game
		GameSystems []models.GameSystem
	}{
		games,
		gameSystems}

	err = t.Execute(w, data)

	if utils.GotError(w, err) {
		return
	}
}

// Add a new game to the db
func AddGame(w http.ResponseWriter, r *http.Request) {
	var gameSystem models.GameSystem
	played := false
	finished := false

	// Parse Form data
	err := r.ParseForm()

	if utils.GotError(w, err) {
		return
	}

	params := r.Form

	// Make sure gamesystem id is really a number
	gameSystemId, err := strconv.Atoi(params.Get("system"))

	if utils.GotError(w, err) {
		return
	}

	// Fetch gamesystem with given id
	store.Db.First(&gameSystem, gameSystemId)

	// Create new game struct and save it in the database
	if params.Get("played") == "on" {
		played = true
	}

	if params.Get("finished") == "on" {
		finished = true
	}

	store.Db.Create(&models.Game{
		Name:        params.Get("name"),
		GameSystems: []models.GameSystem{gameSystem},
		Played:      played,
		Finished:    finished,
	})

	ListGames(w, r)
}

// Fetch a single game from the db and call the show game template
func ShowGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	gameId, err := strconv.Atoi(mux.Vars(r)["game"])

	if utils.GotError(w, err) {
		return
	}

	store.Db.Preload("GameSystems").First(&game, gameId)

	t, err := template.ParseFiles("templates/game/show.html")

	if utils.GotError(w, err) {
		return
	}

	err = t.Execute(w, struct{ Game models.Game }{game})

	if utils.GotError(w, err) {
		return
	}
}

// Fetch a single game from the db and call the edit game template
func ShowEditGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	var gameSystems []models.GameSystem
	gameId, err := strconv.Atoi(mux.Vars(r)["game"])

	if utils.GotError(w, err) {
		return
	}

	store.Db.Preload("GameSystems").First(&game, gameId)
	store.Db.Find(&gameSystems)

	// Sort by Name
	sort.Slice(gameSystems, func(i, j int) bool {
		return sort.StringsAreSorted([]string{gameSystems[i].Name, gameSystems[j].Name})
	})

	// Define template function onGameSystem to check if a GameSystemId is in the Systems of a Game
	funcMap := template.FuncMap{
		"onGameSystem": func(gameOnSystems []models.GameSystem, checkSystemId uint) bool {
			result := false

			for _, x := range gameOnSystems {
				if x.ID == checkSystemId {
					result = true
					break
				}
			}

			return result
		},
	}

	t, err := template.New("edit.html").Funcs(funcMap).ParseFiles("templates/game/edit.html")

	if utils.GotError(w, err) {
		return
	}

	data := struct {
		Game        models.Game
		GameSystems []models.GameSystem
	}{
		game,
		gameSystems}

	err = t.Execute(w, data)

	if utils.GotError(w, err) {
		return
	}
}

// Update a single game in the db
// TODO: Implement me
func EditGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	gameId, err := strconv.Atoi(mux.Vars(r)["game"])

	if utils.GotError(w, err) {
		return
	}

	// Parse Form data
	err = r.ParseForm()

	if utils.GotError(w, err) {
		return
	}

	params := r.Form

	// Fetch game from db and update its properties
	store.Db.Preload("GameSystems").First(&game, gameId)

	game.Name = params.Get("name")

	if params.Get("finished") == "on" {
		game.Finished = true
	}

	if params.Get("played") == "on" {
		game.Played = true
	}

	// Save the updated game

	// Redirect to edit games view
}

// Fetch a single game from the db and call the edit game template
func DeleteGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	gameId, err := strconv.Atoi(mux.Vars(r)["game"])

	if utils.GotError(w, err) {
		return
	}

	store.Db.First(&game, gameId)
	store.Db.Delete(&game)

	ListGames(w, r)
}
