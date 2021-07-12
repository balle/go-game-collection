package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GameSystem struct {
	gorm.Model
	Name   string
	GameID uint
}

type Game struct {
	gorm.Model
	Name         string
	Played       bool
	GameSystemID int
	GameSystems  []GameSystem
}

var db *gorm.DB

func gotError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return true
	}

	return false
}

func listGames(w http.ResponseWriter, r *http.Request) {
	var games []Game
	var gameSystems []GameSystem
	t, err := template.ParseFiles("game.html")

	if gotError(w, err) {
		return
	}

	db.Find(&games)
	db.Find(&gameSystems)

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

func addGame(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if gotError(w, err) {
		return
	}

	var gameSystem GameSystem
	played := false
	params := r.Form
	gameSystemId, err := strconv.Atoi(params.Get("system"))

	if gotError(w, err) {
		return
	}

	db.First(&gameSystem, gameSystemId)

	if params.Get("played") == "on" {
		played = true
	}

	db.Create(&Game{
		Name:        params.Get("name"),
		GameSystems: []GameSystem{gameSystem},
		Played:      played,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func showGame(w http.ResponseWriter, r *http.Request) {
	game := mux.Vars(r)["game"]
	fmt.Fprintf(w, "Showing game %s", game)
}

func insertTestData() {
	nintendoSwitch := GameSystem{Name: "Nintendo Switch"}
	db.Create(&nintendoSwitch)

	playstation1 := GameSystem{Name: "Playstation 1"}
	db.Create(&playstation1)

	playstation3 := GameSystem{Name: "Playstation 3"}
	db.Create(&playstation3)

	gameboy := GameSystem{Name: "Gameboy"}
	db.Create(&gameboy)

	db.Create(&Game{Name: "Need for Speed", GameSystems: []GameSystem{playstation1, playstation3}})
}

func main() {
	// Connect to the database
	var err error
	db, err = gorm.Open(sqlite.Open("games.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&GameSystem{}, &Game{})

	insertTestData()

	// Setup webserver
	r := mux.NewRouter()
	r.HandleFunc("/", listGames)
	r.HandleFunc("/game", addGame).Methods(http.MethodPost)
	r.HandleFunc("/game/{game:[0-9]+}", showGame)
	http.ListenAndServe(":8888", r)
}
