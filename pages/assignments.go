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

func AssignmentsDataForACourse(c *colly.Collector, studentId int, mpToView string, courseCode string, courseSection string, courseName string) []models.Assignment {
	assignments := make([]models.Assignment, 0)

	data := constants.ConstantLinks["assignments"]
	data["studentid"] = strconv.Itoa(studentId)
	data["mpToView"] = mpToView
	data["courseCode"] = courseCode
	data["courseSection"] = courseSection
	assignemnts_url, err := utils.FormatDynamicUrl(data)
	if err != nil {
		log.Println(err)
		return assignments

	}

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		rows := dom.Find(".list > tbody>tr")
		rows.Each(func(i int, row *goquery.Selection) {
			if row.Children().Length() == constants.CourseSummaryRowLength && i != 1 {
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

				stuff := utils.BasicDataExtractor(row, courseName)
				aAssignment.Category = stuff[constants.CourseSummaryNameCategory]
				aAssignment.Assignment = stuff[constants.CourseSummaryNameAssignment]
				aAssignment.Description = stuff[constants.CourseSummaryNameDescription]
				aAssignment.CourseName = courseName

				tds := row.Children()
				tds.Each(func(i int, s *goquery.Selection) {
					// if i ==
				})

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
