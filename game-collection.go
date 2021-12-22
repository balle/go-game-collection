package main

/*
 * IMPORTS
 */

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/*
 * TYPES AND VARIABLES
 */
type GameSystem struct {
	gorm.Model
	Name   string
	GameID uint
}

type Game struct {
	gorm.Model
	Name         string
	Played       bool
	Finished     bool
	GameSystemID int
	GameSystems  []GameSystem
}

var db *gorm.DB

/*
 * FUNCTIONS
 */

// Got an error? Log it and return error code 500 to browser
func gotError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return true
	}

	return false
}

// Receive all games from the db and call the game list template
func listGames(w http.ResponseWriter, r *http.Request) {
	var games []Game
	var gameSystems []GameSystem
	t, err := template.ParseFiles("templates/game/list.html")

	if gotError(w, err) {
		return
	}

	// Fetch all games and game systems from db
	db.Find(&games)
	db.Find(&gameSystems)

	// Sort them by name
	sort.Slice(games, func(i, j int) bool {
		return sort.StringsAreSorted([]string{games[i].Name, games[j].Name})
	})

	sort.Slice(gameSystems, func(i, j int) bool {
		return sort.StringsAreSorted([]string{gameSystems[i].Name, gameSystems[j].Name})
	})

	// Compile template with our data
	data := struct {
		Games       []Game
		GameSystems []GameSystem
	}{
		games,
		gameSystems}

	err = t.Execute(w, data)

	if gotError(w, err) {
		return
	}
}

// Add a new game to the db
func addGame(w http.ResponseWriter, r *http.Request) {
	var gameSystem GameSystem
	played := false
	finished := false

	// Parse Form data
	err := r.ParseForm()

	if gotError(w, err) {
		return
	}

	params := r.Form

	// Make sure gamesystem id is really a number
	gameSystemId, err := strconv.Atoi(params.Get("system"))

	if gotError(w, err) {
		return
	}

	// Fetch gamesystem with given id
	db.First(&gameSystem, gameSystemId)

	// Create new game struct and save it in the database
	if params.Get("played") == "on" {
		played = true
	}

	if params.Get("finished") == "on" {
		finished = true
	}

	db.Create(&Game{
		Name:        params.Get("name"),
		GameSystems: []GameSystem{gameSystem},
		Played:      played,
		Finished:    finished,
	})

	listGames(w, r)
}

// Fetch a single game from the db and call the show game template
func showGame(w http.ResponseWriter, r *http.Request) {
	var game Game
	gameId, err := strconv.Atoi(mux.Vars(r)["game"])

	if gotError(w, err) {
		return
	}

	db.Preload("GameSystems").First(&game, gameId)

	t, err := template.ParseFiles("templates/game/show.html")

	if gotError(w, err) {
		return
	}

	err = t.Execute(w, struct{ Game Game }{game})

	if gotError(w, err) {
		return
	}
}

// Fetch a single game from the db and call the edit game template
func showEditGame(w http.ResponseWriter, r *http.Request) {
	var game Game
	var gameSystems []GameSystem
	gameId, err := strconv.Atoi(mux.Vars(r)["game"])

	if gotError(w, err) {
		return
	}

	db.Preload("GameSystems").First(&game, gameId)
	db.Find(&gameSystems)

	sort.Slice(gameSystems, func(i, j int) bool {
		return sort.StringsAreSorted([]string{gameSystems[i].Name, gameSystems[j].Name})
	})

	funcMap := template.FuncMap{
		"onGameSystem": func(gameOnSystems []GameSystem, checkSystemId uint) bool {
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

	if gotError(w, err) {
		return
	}

	data := struct {
		Game        Game
		GameSystems []GameSystem
	}{
		game,
		gameSystems}

	err = t.Execute(w, data)

	if gotError(w, err) {
		return
	}
}

// Update a single game in the db
// TODO: Implement me
func editGame(w http.ResponseWriter, r *http.Request) {
	var game Game
	gameId, err := strconv.Atoi(mux.Vars(r)["game"])

	if gotError(w, err) {
		return
	}

	// Parse Form data
	err = r.ParseForm()

	if gotError(w, err) {
		return
	}

	params := r.Form

	// Fetch game from db and update its properties
	db.Preload("GameSystems").First(&game, gameId)

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
func deleteGame(w http.ResponseWriter, r *http.Request) {
	var game Game
	gameId, err := strconv.Atoi(mux.Vars(r)["game"])

	if gotError(w, err) {
		return
	}

	db.First(&game, gameId)
	db.Delete(&game)

	listGames(w, r)
}

// Add a new game system to the db
func addGameSystem(w http.ResponseWriter, r *http.Request) {
	// Parse Form data
	err := r.ParseForm()

	if gotError(w, err) {
		return
	}

	params := r.Form

	// Create new game system struct and save it in the database
	db.Create(&GameSystem{
		Name: params.Get("name"),
	})

	listGames(w, r)
}

/*
 * MAIN FUNCTION
 */
func main() {
	// Connect to the database
	var err error
	db, err = gorm.Open(sqlite.Open("games.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&GameSystem{}, &Game{})

	// Setup routes, register fileserver for images and run the webserver
	r := mux.NewRouter()
	r.HandleFunc("/", listGames)
	r.HandleFunc("/game", addGame).Methods(http.MethodPost)
	r.HandleFunc("/gamesystem", addGameSystem).Methods(http.MethodPost)
	r.HandleFunc("/game/edit/{game:[0-9]+}", editGame).Methods(http.MethodPost)
	r.HandleFunc("/game/edit/{game:[0-9]+}", showEditGame).Methods(http.MethodGet)
	r.HandleFunc("/game/delete/{game:[0-9]+}", deleteGame)
	r.HandleFunc("/game/{game:[0-9]+}", showGame)

	images := http.StripPrefix("/images/", http.FileServer(http.Dir("./templates/images/")))
	r.PathPrefix("/images/").Handler(images)

	http.ListenAndServe(":8888", r)
}
