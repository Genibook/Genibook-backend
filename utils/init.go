package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func Init_colly() (*colly.Collector, error) {

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return c, err
	}

	Login(c, os.Getenv("username"), os.Getenv("password"))

	c.OnRequest(func(r *colly.Request) {
		// Print the URL of the page being visited
		fmt.Println("~~~~~~~~~~")
		fmt.Println("Visiting:", r.URL.String())
		fmt.Println("~~~~~~~~~~")
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		title := e.Text
		fmt.Println("------------")
		fmt.Print("We are scrapping the website with title: ")
		fmt.Println(title)
		fmt.Println("------------")
	})

	fmt.Println("finished initialization!")
	return c, nil
}
