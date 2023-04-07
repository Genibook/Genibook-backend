package pages

import (
	"log"
	"webscrapper/constants"
	"webscrapper/utils"

	"github.com/gocolly/colly"
)

func GradebookData(c *colly.Collector, studentId int, mpToView string) {

	data := constants.ConstantLinks["gradebook"]
	gradebook_url, err := utils.FormatDynamicUrl(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	c.Visit(gradebook_url)
}
