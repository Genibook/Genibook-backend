package utils

import (
	"fmt"
	"net/http"
)

// prints the error if the error is not nil
// returns http ___ code
// also logs the string and error
func APIPrintSpecificError(prompt string, w http.ResponseWriter, theError error, code int) {
	if theError != nil {
		fmt.Printf("Error prompt: %v\n", prompt)
		fmt.Printf("Error: %v\n", theError.Error())
		http.Error(w, theError.Error(), code)
		return
	}
}
