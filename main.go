package main

import (
	"log"
	"net/http"
	"webscrapper/utils"
)

func main() {
	c, e := utils.Init_colly()
	utils.DebugLogin(c)
	if e != nil {
		log.Fatal("error initing colly scrapper")
		log.Fatal(e)
		return
	}
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
	// http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
	//pages.GradebookData(c, 107604, "MP1")

}
