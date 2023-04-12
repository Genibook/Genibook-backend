package utils

import (
	"fmt"
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
		log.Fatal(err)
		log.Fatal("Failed login")
	}
	fmt.Print("successful login!\n")

}
