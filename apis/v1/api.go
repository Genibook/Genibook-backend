package api_v1

import (
	"log"
	"net/http"
	"strconv"
	"webscrapper/constants"
	"webscrapper/utils"

	"github.com/gin-gonic/gin"
)

// var validPath = regexp.MustCompile("^/(edit|login|profile|grades|assignments|schedule|student)/")

func MakeHandler(fn func(*gin.Context, http.ResponseWriter, *http.Request, string, string, string, int)) func(c *gin.Context) {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		//, w http.ResponseWriter, r *http.Request
		err := r.ParseForm()
		if err != nil {
			utils.APIPrintSpecificError("Error parsing the post data's form :/", w, err, http.StatusInternalServerError)
			return
		}

		userSelectorString := c.Query(constants.UserSelectorFormKey)
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
		key := c.Query(constants.HighSchoolFormKey)
		kValid := false
		for k := range constants.ConstantLinks {
			if k == key {
				kValid = true
			}
		}
		if !kValid {
			log.Println("Someone tried to use a sussy highschool")
			http.Error(w, "High School Not Available", http.StatusNoContent)
			return
		}

		fn(c, w, r, c.Query(constants.UsernameFormKey), c.Query(constants.PasswordFormKey), key, userSelector)
	}
}

func LoginHandlerV1(context *gin.Context, w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {
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

func ProfileHandlerV1(c *gin.Context, w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {
	functionName := "Func ProfileHandlerV1"

	student, err := GetProfile(w, functionName, email, password, highSchool, userSelector)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, student)
}

// <note>: userSelector is 1st indexed meaning the first user is 1, second is 2.
// Backend processes it like that
func GradesHandlerV1(c *gin.Context, w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {

	functionName := "Func GradesHandlerV1"

	grades, err := GetGrades(w, r, functionName, email, password, highSchool, userSelector)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, grades)

}

func AssignmentHandlerV1(c *gin.Context, w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {

	functionName := "Func AssignmentHandlerV1"

	assignments, err := GetAssignments(w, r, functionName, email, password, highSchool, userSelector)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, assignments)

}

func ScheduleAssignmentHandlerV1(c *gin.Context, w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {

	functionName := "Func ScheduleAssignmentHandlerV1"

	scheduleAssignments, err := GetSchedule(w, r, functionName, email, password, highSchool, userSelector)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, scheduleAssignments)
}

func StudentHandlerV1(c *gin.Context, w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int) {

	functionName := "Func StudentHandlerV1"

	student, err := GetProfile(w, functionName, email, password, highSchool, userSelector)
	if err != nil {
		return
	}

	grades, err := GetGrades(w, r, functionName, email, password, highSchool, userSelector)
	if err != nil {
		return
	}

	assignments, err := GetAssignments(w, r, functionName, email, password, highSchool, userSelector)
	if err != nil {
		return
	}

	ret := CombineGradeAssiandProfile(assignments, grades, student)

	c.JSON(http.StatusOK, ret)

}
