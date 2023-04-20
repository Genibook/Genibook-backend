package api_v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"webscrapper/constants"
	"webscrapper/pages"
	"webscrapper/utils"
)

var validPath = regexp.MustCompile("^/(edit|login|profile|grades|)/")

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string, string, string, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.Find([]byte(r.URL.Path))
		if m == nil {
			http.NotFound(w, r)
			return
		}
		err := r.ParseForm()
		if err != nil {
			utils.APIPrintSpecificError("Error parsing the post data's form :/", w, err, http.StatusInternalServerError)
			return
		}

		userSelectorString := r.URL.Query().Get(constants.UserSelectorFormKey)
		userSelector, err := strconv.Atoi(userSelectorString)
		if err != nil {
			utils.APIPrintSpecificError("Error converting form value with key 'user' to integer: "+userSelectorString, w, err, http.StatusInternalServerError)
			return
		}
		if userSelector <= 0 {
			log.Println("Someone tried to use a userselector of <= 0")
			http.Error(w, "user key is <=0", http.StatusNotAcceptable)
			return
		}

		fn(w, r, r.URL.Query().Get(constants.UsernameFormKey), r.URL.Query().Get(constants.PasswordFormKey), r.URL.Query().Get(constants.HighSchoolFormKey), userSelector)
	}
}

func LoginHandlerV1(w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {
	c := utils.Init_colly()
	e := utils.Login(c, email, password, highSchool)

	if e != nil {
		log.Println("Func Login Hanlder - Incorrect Password and Username <Note: It is OK if this happens>")
		http.Error(w, e.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)

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

func ProfileHandlerV1(w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {
	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError("Func Profile Handler V1: Couldn't init/login", w, e, http.StatusInternalServerError)
	student := pages.ProfileData(c, userSelector, highSchool)
	jsonData, e := student.ToJson()
	utils.APIPrintSpecificError("Func Profile Handler V1: Json Parsing Error", w, e, http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

// <note>: userSelector is 1st indexed meaning the first user is 1, second is 2.
// Backend processes it like that
func GradesHandlerV1(w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {
	mp := r.URL.Query().Get(constants.MPFormKey)

	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError("Func GradesHandlerV1: Couldn't init/login", w, e, http.StatusInternalServerError)

	IDS := pages.StudentIdAndCurrMP(c, highSchool)

	if userSelector > len(IDS) {
		log.Printf("User selector index >= len(available IDS) Length: %d\n", len(IDS))
		http.Error(w, fmt.Sprintf("User selector index >= len(available IDS) Length: %d", len(IDS)), http.StatusNotAcceptable)
		return
	}

	weeklySumData := pages.GradebookData(c, IDS[userSelector-1], mp, highSchool)

	grades := map[string]map[string]string{}
	for key := range weeklySumData {
		oneGrade := weeklySumData[key]
		grades[key] = oneGrade.ToDict()
	}
	jsonData, e := json.Marshal(grades)
	utils.APIPrintSpecificError("Func GradesHandlerV1: Json Parsing Error", w, e, http.StatusInternalServerError)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}
