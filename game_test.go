package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/balle/go-game-collection/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", controllers.ListGames)
	return router
}

func get(url string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("GET", url, nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	return response
}

func defaultTests(t *testing.T, url string, expected string) {
	response := get(url)

	if response.Code != 200 {
		t.Errorf("Got response code %d expected 200", response.Code)
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Errorf("Failed to read body")
	} else if strings.Contains(string(body), "Error") {
		t.Errorf("Request %s returned error: %s", url, body)
	} else if !strings.Contains(string(body), expected) {
		t.Errorf("Expected response %s \nGot %s", expected, body)
	}
}

func TestListGames(t *testing.T) {
	defaultTests(t, "/", "My games")
}
