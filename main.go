package main

import (
	"log"
	"net/http"
	api_v1 "webscrapper/apis/v1"
)

func main() {
	// c, e := utils.Init_colly()
	// if e != nil {
	// 	log.Fatal("error initing colly scrapper")
	// 	log.Fatal(e)
	// 	return
	// }
	// utils.DebugLogin(c)
	// // student := pages.ProfileData(c, 1)
	// // fmt.Printf("student: %v\n", student)
	// //pages.ProfileData(c, 1)
	// //pages.GimmeCourseCodes(c, 107604, "MP1")
	// // assignments := pages.AssignmentsDataForACourse(c, 107604, "MP2", "22000", "2", "US History I")
	// // for _, assignment := range assignments {
	// // 	fmt.Println(assignment)

	// // }
	// assignments_more := pages.ScheduleDataForACourse(c, 107604, "MP2", "22000", "2", "US History I")
	// for _, assignment := range assignments_more {
	// 	fmt.Println(assignment)

	// }
	// http.HandleFunc("/view/", makeHandler(viewHandler))
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/login/", api_v1.MakeHandler(api_v1.LoginHandlerV1))
	log.Fatal(http.ListenAndServe(":8080", nil))
	//pages.GradebookData(c, 107604, "MP1")

}
