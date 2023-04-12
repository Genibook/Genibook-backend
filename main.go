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
	pages.ProfileData(c, 1)

}
