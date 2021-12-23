package main

/*
 * IMPORTS
 */

import (
	"net/http"

	"github.com/balle/go-game-collection/controllers"
	"github.com/gorilla/mux"
)

/*
 * MAIN FUNCTION
 */
func main() {
	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.ListGames)
	r.HandleFunc("/game", controllers.AddGame).Methods(http.MethodPost)
	r.HandleFunc("/gamesystem", controllers.AddGameSystem).Methods(http.MethodPost)
	r.HandleFunc("/game/edit/{game:[0-9]+}", controllers.EditGame).Methods(http.MethodPost)
	r.HandleFunc("/game/edit/{game:[0-9]+}", controllers.ShowEditGame).Methods(http.MethodGet)
	r.HandleFunc("/game/delete/{game:[0-9]+}", controllers.DeleteGame)
	r.HandleFunc("/game/{game:[0-9]+}", controllers.ShowGame)

	// Register fileserver for images
	images := http.StripPrefix("/images/", http.FileServer(http.Dir("./templates/images/")))
	r.PathPrefix("/images/").Handler(images)

	// Run the webserver
	http.ListenAndServe(":8888", r)
}
