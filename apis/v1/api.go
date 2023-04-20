package api_v1

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"webscrapper/constants"
	"webscrapper/utils"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string, string, string, string, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			log.Println("Error parsing the post data's form :/")
		}
		fmt.Println(m)

		userSelectorString := r.PostFormValue(constants.UserSelectorFormKey)
		userSelector, err := strconv.Atoi(userSelectorString)
		if err != nil {
			// ... handle error
			log.Println(err)
			log.Println(userSelectorString)
			log.Println("Error converting form value with key 'user' to integer^^")

		}

		fn(w, r, m[2], r.PostFormValue(constants.UsernameFormKey), r.PostFormValue(constants.PasswordFormKey), r.PostFormValue(constants.HighSchoolFormKey), userSelector)
	}
}

func LoginHandlerV1(w http.ResponseWriter, r *http.Request, title string, email string, password string, highSchool string, userSelector int) {
	c, e := utils.Init_colly()
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	e = utils.Login(c, email, password, highSchool)
	if e != nil {
		http.Error(w, e.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
	// data := map[string]string{
	// 	"name":  "John",
	// 	"email": "john@example.com",
	// }

	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(jsonData)
}
