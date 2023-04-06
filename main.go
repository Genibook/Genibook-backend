package main

import (
	"log"
	"webscrapper/utils"

	"webscrapper/pages"
)

func main() {
	c, e := utils.Init_colly()
	if e != nil {
		log.Fatal("error initing colly scrapper")
		log.Fatal(e)
		return
	}
	pages.ProfileData(c)

}
