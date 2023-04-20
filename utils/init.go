package utils

import (
	"log"
	"os"
	"webscrapper/constants"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func Init_colly() (*colly.Collector, error) {

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// c.OnRequest(func(r *colly.Request) {
	// 	// Print the URL of the page being visited
	// 	fmt.Println("~~~~~~~~~~")
	// 	fmt.Println("Visiting:", r.URL.String())
	// 	fmt.Println("~~~~~~~~~~")
	// })

	// c.OnHTML("title", func(e *colly.HTMLElement) {
	// 	title := e.Text
	// 	fmt.Println("------------")
	// 	fmt.Print("We are scrapping the website with title: ")
	// 	fmt.Println(title)
	// 	fmt.Println("------------")
	// })

	log.Println("finished initialization!")
	return c, nil
}

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

func DebugLogin(c *colly.Collector) (*colly.Collector, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return c, err
	}

	Login(c, os.Getenv("username"), os.Getenv("password"))

	return c, nil
}
