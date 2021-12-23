package utils

import (
	"log"
	"net/http"
)

// Got an error? Log it and return error code 500 to browser
func GotError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return true
	}

	return false
}
