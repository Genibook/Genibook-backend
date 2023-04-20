package utils

import (
	"fmt"
	"log"
	"os"
	"time"
	"webscrapper/constants"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func Init_colly() *colly.Collector {

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	now := time.Now()
	fmt.Println("Connection Established With Backend at Time: " + now.Format("2006-01-02 15:04:05"))
	return c
}

func Login(c *colly.Collector, username string, password string, school string) error {
	// authenticate
	data := map[string]string{
		constants.ConstantLinks[school]["login"]["username"]: username,
		constants.ConstantLinks[school]["login"]["password"]: password,
	}
	err := c.Post(constants.ConstantLinks[school]["login"]["url"], data)
	if err != nil {
		log.Println(err)
		log.Println("Failed login")
		// maybe an api call to login or smth and then deal with it like that return the error etc
		return err
	}
	//log.Println("Successful login!")
	return nil

}

func InitAndLogin(email string, password string, highSchool string) (*colly.Collector, error) {
	c := Init_colly()
	e := Login(c, email, password, highSchool)

	if e != nil {
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
