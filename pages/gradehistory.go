package pages

import (
	"log"
	"strconv"
	"strings"
	"webscrapper/constants"
	"webscrapper/models"
	"webscrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func GradeHistoryData(c *colly.Collector, studentId string, school string) (map[string][]map[string]models.Course, error) {
	courses := map[string][]map[string]models.Course{}

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		table := dom.Find(".list")
		trs := table.Find("tr.listroweven, tr.listrowodd")
		trs.Each(func(i int, s *goquery.Selection) {
			if i < trs.Length()-2 {
				course := models.Course{}
				s.Find("td").Each(func(i int, td *goquery.Selection) {
					switch i {
					case 0:

						course.SchoolYear = utils.CleanAString(td.Text())
					case 1:
						grade, err := strconv.Atoi(utils.CleanAString(td.Text()))
						if err != nil {
							return
						}
						course.Grade = grade
					case 2:
						course.Name = strings.TrimSpace(utils.CleanAString(td.Text()))
					case 3:
						course.School = utils.CleanAString(td.Text())
					case 4:
						course.FG = utils.CleanAString(td.Text())
					case 5:
						attempted, err := strconv.ParseFloat(utils.CleanAString(td.Text()), 64)
						if err != nil {
							return
						}
						course.Attempted = float32(attempted)
					case 6:
						earned, err := strconv.ParseFloat(utils.CleanAString(td.Text()), 64)
						if err != nil {
							return
						}
						course.Earned = float32(earned)
					}
				})

				if course.SchoolYear != "" {
					courses[course.SchoolYear] = append(courses[course.SchoolYear], map[string]models.Course{course.Name: course})
				}
			}

		})
	})

	data := constants.ConstantLinks[school]["history"]
	data["studentid"] = studentId
	history_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		log.Println(err)
		return courses, err

	}

	err = c.Visit(history_url)

	if err != nil {
		log.Println("Couldn't visit GradeHistoryData url: function GradebookData, file gradehistory.go.go")
		return courses, err
	}

	c.OnHTMLDetach("body")

	return courses, nil
}
func CurrentGradeHistoryData(c *colly.Collector, studentId string, school string) (map[string]string, error) {
	currCourses := map[string]string{}
	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		table := dom.Find(".list")
		trs := table.Find("tbody>tr.listroweven, tbody>tr.listrowodd")
		trs.Each(func(i int, s *goquery.Selection) {
			name := ""
			att := ""

			s.Find("td").Each(func(k int, td *goquery.Selection) {
				switch k {
				case 0:
					name = utils.CleanAString(td.Text())
				case 6:
					att = utils.CleanAString(td.Text())
				}
			})

			if name != "" {
				currCourses[name] = att
			}

		})
	})

	data := constants.ConstantLinks[school]["currHistory"]
	data["studentid"] = studentId
	history_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		log.Println(err)
		return currCourses, err

	}

	err = c.Visit(history_url)

	if err != nil {
		log.Println("Couldn't visit history_url: function CurrentGradeHistoryData, file gradehistory.go")
		return currCourses, err
	}

	c.OnHTMLDetach("body")

	return currCourses, nil
}
