package api_v1

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"webscrapper/constants"
	"webscrapper/pages"

	"github.com/gocolly/colly"
)

func GetMP(w http.ResponseWriter, r *http.Request) (string, error) {
	mp := r.URL.Query().Get(constants.MPFormKey)
	if !strings.Contains(mp, "MP") {
		log.Println("Marking Period Not Valid: " + mp)
		http.Error(w, "Marking Period Not Valid: "+mp, http.StatusNotAcceptable)
		return "", http.ErrBodyNotAllowed
	}
	return mp, nil
}

func GetIDs(userSelector int, c *colly.Collector, highSchool string, w http.ResponseWriter) ([]string, error) {
	IDS := pages.StudentIdAndCurrMP(c, highSchool)
	if userSelector > len(IDS) {
		log.Printf("User selector index > len(available IDS) Length: %d\n", len(IDS))
		http.Error(w, fmt.Sprintf("User selector index > len(available IDS) Length: %d", len(IDS)), http.StatusNotAcceptable)
		return make([]string, 0), http.ErrBodyNotAllowed
	}
	return IDS, nil
}
