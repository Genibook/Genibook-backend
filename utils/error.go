package utils

import (
	"log"
	"net/http"
)

func APIPrintSpecificError(prompt string, w http.ResponseWriter, theError error, code int) {
	if theError != nil {
		log.Println(prompt)
		http.Error(w, theError.Error(), code)
		return
	}
}
