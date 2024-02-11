package utils

import (
	"log"
	"os"
	constants "webscrapper/constants/v1"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func Init_colly() *colly.Collector {

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"
	log.Println("Connection Established With Backend")
	return c
}

func Login(c *colly.Collector, username string, password string, school string) (*colly.Collector, error) {
	// authenticate
	data := map[string]string{
		constants.ConstantLinks[school]["login"]["username"]: username,
		constants.ConstantLinks[school]["login"]["password"]: password,
	}
	// c.OnResponse(func(r *colly.Response) {
	// 	cookies := c.Cookies(r.Request.URL.String())
	// 	for _, cookie := range cookies {
	// 		log.Println("Cookie:", cookie.Name, "Value:", cookie.Value)
	// 	}
	// })
	err := c.Post(constants.ConstantLinks[school]["login"]["url"], data)
	// c.OnResponse(func(r *colly.Response) {})
	if err != nil {
		log.Println(err)
		log.Println("[ERROR Login() init.go]: Failed login")
		// maybe an api call to login or smth and then deal with it like that return the error etc
		return c, err
	}
	return c, nil

}

func InitAndLogin(email string, password string, highSchool string) (*colly.Collector, error) {
	c := Init_colly()
	c, e := Login(c, email, password, highSchool)

	if e != nil {
		log.Println("[ERROR InitAndLogin() init.go]")
		return c, e
	}
	return c, e
}

func DebugLogin(c *colly.Collector) (*colly.Collector, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return c, err
	}

	Login(c, os.Getenv("username"), os.Getenv("password"), constants.MontgomeryHighSchoolKeyName)

	return c, nil
}
