package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/balle/go-game-collection/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", controllers.ListGames)
	return router
}

func TestListGames(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	if response.Code != 200 {
		t.Errorf("Got response code %d expected 200", response.Code)
	}
}
