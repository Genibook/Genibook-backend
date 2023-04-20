package main

import (
	"log"
	"net/http"
	api_v1 "webscrapper/apis/v1"
)

func main() {
	// c := utils.Init_colly()
	// utils.DebugLogin(c)
	// student := pages.ProfileData(c, 1, constants.MontgomeryHighSchoolKeyName)
	// fmt.Printf("student: %v\n", student)
	// //pages.GimmeCourseCodes(c, 107604, "MP1")
	// // assignments := pages.AssignmentsDataForACourse(c, 107604, "MP2", "22000", "2", "US History I")
	// // for _, assignment := range assignments {
	// // 	fmt.Println(assignment)
	// fmt.Printf("pages.StudentIdAndCurrMP(c): %v\n", pages.StudentIdAndCurrMP(c, constants.MontgomeryHighSchoolKeyName))
	// // }
	//pages.GradebookData(c, 107604, "MP1", constants.MontgomeryHighSchoolKeyName)
	// assignments_more := pages.ScheduleDataForACourse(c, 107604, "MP2", "22000", "2", "US History I")
	// for _, assignment := range assignments_more {
	// 	fmt.Println(assignment)

	// }
	// http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/student/", api_v1.MakeHandler(api_v1.StudentHandlerV1))
	http.HandleFunc("/schedule/", api_v1.MakeHandler(api_v1.ScheduleAssignmentHandlerV1))
	http.HandleFunc("/assignments/", api_v1.MakeHandler(api_v1.AssignmentHandlerV1))
	http.HandleFunc("/grades/", api_v1.MakeHandler(api_v1.GradesHandlerV1))
	http.HandleFunc("/profile/", api_v1.MakeHandler(api_v1.ProfileHandlerV1))
	http.HandleFunc("/hi/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hiiii")) })
	http.HandleFunc("/login/", api_v1.MakeHandler(api_v1.LoginHandlerV1))
	log.Fatal(http.ListenAndServe(":6969", nil))

}
