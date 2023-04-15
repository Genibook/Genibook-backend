package pages

import (
	"log"
	"strconv"
	"webscrapper/constants"
	"webscrapper/models"
	"webscrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func AssignmentsDataForACourse(c *colly.Collector, studentId int, mpToView string, courseCode string, courseSection string) []models.Assignment {
	assignments := make([]models.Assignment, 0)

	data := constants.ConstantLinks["gradebook"]
	data["studentid"] = strconv.Itoa(studentId)
	data["mpToView"] = mpToView
	assignemnts_url, err := utils.FormatDynamicUrl(data)
	if err != nil {
		log.Println(err)
		return assignments

	}

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		rows := dom.Find(".list > tbody>tr")
		rows.Each(func(i int, row *goquery.Selection) {
			if i != 0 && i != 1 && i != rows.Length()-1 && i != rows.Length()-2 && i != rows.Length()-2 {
				aAssignment := models.Assignment{
					CourseName:   "",
					MP:           "",
					DayName:      "",
					FullDayName:  "",
					Date:         "",
					FullDate:     "",
					Teacher:      "",
					Category:     "",
					Assignment:   "",
					Description:  "",
					GradePercent: "",
					GradeNum:     "",
					Comment:      "",
					Prev:         "",
					Docs:         "",
				}
				assignments = append(assignments, aAssignment)
			}
		})

	})

	err = c.Visit(assignemnts_url)
	if err != nil {
		log.Println("Couldn't visit assignment url: function AssignmentsDataForACourse, file assignments.go")
	}
	c.OnHTMLDetach("body")

	return assignments
}
