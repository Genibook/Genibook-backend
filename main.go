package main

import (
	"log"
	"webscrapper/pages"
	"webscrapper/utils"
)

func main() {
	c, e := utils.Init_colly()
	if e != nil {
		log.Fatal("error initing colly scrapper")
		log.Fatal(e)
		return
	}
	// student := pages.ProfileData(c, 1)
	// fmt.Printf("student: %v\n", student)
	//pages.ProfileData(c, 1)
	//pages.GimmeCourseCodes(c, 107604, "MP1")
	pages.AssignmentsDataForACourse(c, 107604, "MP2", "22000", "2", "US History I")
	//pages.GradebookData(c, 107604, "MP1")

}
