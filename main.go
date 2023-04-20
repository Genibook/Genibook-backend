package main

import (
	"fmt"
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
	assignments := pages.AssignmentsDataForACourse(c, 107604, "MP2", "22000", "2", "US History I")
	for _, assignment := range assignments {
		fmt.Println(assignment)

	}
	//pages.GradebookData(c, 107604, "MP1")

}
