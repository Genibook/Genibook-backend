package api_v1

import (
	"encoding/json"
	"net/http"
	"webscrapper/utils"
)

func ReturnJsonData(object interface{}, w http.ResponseWriter, name string) {
	jsonData, e := json.Marshal(object)

	// this handles the erro already
	utils.APIPrintSpecificError(name, w, e, http.StatusInternalServerError)

	if e != nil {
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
	return
}
