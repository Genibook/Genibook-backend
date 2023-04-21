package api_v1

import (
	"webscrapper/models"
)

// func ReturnJsonData(object interface{}, w http.ResponseWriter, name string) error {
// 	jsonData, e := json.Marshal(object)

// 	// this handles the erro already
// 	utils.APIPrintSpecificError(name, w, e, http.StatusInternalServerError)

// 	if e != nil {
// 		//technically i don't need to handle the error in api.go cuz its the last line and if i return here
// 		// it means that i "break" here so like the code won't run and error messages would display :D
// 		return e
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(jsonData))
// 	return nil
// }

func CombineGradeAssiandProfile(assignments map[string][]models.Assignment, grades map[string]map[string]string, profile models.Student) map[string]interface{} {
	student := profile.ToDict()
	student["assignments"] = assignments
	student["grades"] = grades

	return student

}
