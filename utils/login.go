package utils

import (
	"log"
	"webscrapper/constants"

	"github.com/gocolly/colly"
)

func Login(c *colly.Collector, username string, password string) {
	// authenticate
	data := map[string]string{
		constants.ConstantLinks["login"]["username"]: username,
		constants.ConstantLinks["login"]["password"]: password,
	}
	err := c.Post(constants.ConstantLinks["login"]["url"], data)
	if err != nil {
		log.Println(err)
		log.Println("Failed login")
		// maybe an api call to login or smth and then deal with it like that return the error etc
	}
	log.Println("successful login!")

}
