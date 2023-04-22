package pages

import (
	"log"
	"strings"
	"webscrapper/constants"
	"webscrapper/models"
	"webscrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func AssignmentsDataForACourse(c *colly.Collector, studentId string, mpToView string, courseCode string, courseSection string, courseName string, school string) []models.Assignment {
	assignments := make([]models.Assignment, 0)

	data := constants.ConstantLinks[school]["assignments"]
	data["studentid"] = studentId
	data["mpToView"] = mpToView
	data["courseCode"] = courseCode
	data["courseSection"] = courseSection
	assignemnts_url, err := utils.FormatDynamicUrl(data, school)
	//fmt.Printf("assignemnts_url: %v\n", assignemnts_url)
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
				tds := row.Children()

				stuff := utils.BasicDataExtractor(row, courseName)
				aAssignment.Category = stuff[constants.CourseSummaryNameCategory]
				aAssignment.Assignment = stuff[constants.CourseSummaryNameAssignment]
				aAssignment.Description = stuff[constants.CourseSummaryNameDescription]
				aAssignment.CourseName = courseName

				tds.Each(func(i int, s *goquery.Selection) {
					if i == constants.CourseSummaryMPIndex {
						aAssignment.MP = strings.TrimSpace(s.Text())
					} else if i == constants.CourseSummaryDueIndex {
						dayname, date := utils.ProcessDueCell(s)
						aAssignment.DayName = strings.TrimSpace(dayname)
						aAssignment.FullDayName = utils.DayClassifier(aAssignment.DayName)
						aAssignment.Date = strings.TrimSpace(date)
						aAssignment.FullDate = aAssignment.Date
					} else if i == constants.CourseSummaryGradeIndex {
						aAssignment.GradeNum, aAssignment.GradePercent = utils.ProcessGradeCellForAssignment(s)
					} else if i >= constants.CourseSummaryCommentIndex {
						aAssignment.Comment = utils.CleanAString(s.Text())
					} else if i == constants.CourseSummaryPrevIndex {
						aAssignment.Prev = utils.CleanAString(s.Text())
					} else if i == constants.CourseSummaryDocsIndex {
						aAssignment.Docs = utils.CleanAString(s.Text())
					}
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

	//fmt.Println(assignments)

	return assignments
}

func ScheduleDataForACourse(c *colly.Collector, studentId string, mpToView string, courseCode string, courseSection string, courseName string, school string) []models.ScheduleAssignment {
	assignments := make([]models.ScheduleAssignment, 0)

	data := constants.ConstantLinks[school]["assignments"]
	data["studentid"] = studentId
	data["mpToView"] = mpToView
	data["courseCode"] = courseCode
	data["courseSection"] = courseSection
	assignemnts_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		log.Println(err)
		return assignments
	}

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		rows := dom.Find(".list > tbody>tr")
		rows.Each(func(i int, row *goquery.Selection) {
			if row.Children().Length() == constants.CourseSummaryRowLength && i != 1 {
				notGraded := false
				tds := row.Children()
				tds.Each(func(i int, s *goquery.Selection) {

					if i == constants.CourseSummaryGradeIndex {
						_, notGraded = utils.ProcessGradeCellForSchedule(s)

					}
				})

				if notGraded {
					aAssignment := models.ScheduleAssignment{
						CourseName:  "",
						Category:    "",
						Assignment:  "",
						Description: "",
						Date:        "",
						Points:      "",
					}

					stuff := utils.BasicDataExtractor(row, courseName)
					aAssignment.Category = stuff[constants.CourseSummaryNameCategory]
					aAssignment.Assignment = stuff[constants.CourseSummaryNameAssignment]
					aAssignment.Description = stuff[constants.CourseSummaryNameDescription]
					aAssignment.CourseName = courseName

					tds.Each(func(i int, s *goquery.Selection) {
						if i == constants.CourseSummaryDueIndex {
							_, date := utils.ProcessDueCell(s)
							aAssignment.Date = strings.TrimSpace(date)
						} else if i == constants.CourseSummaryGradeIndex {
							aAssignment.Points, _ = utils.ProcessGradeCellForSchedule(s)
						}
					})

					assignments = append(assignments, aAssignment)
				}

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
