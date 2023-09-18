package pages

import (
	"fmt"
	"log"
	"strings"
	constants "webscrapper/constants/v1"
	"webscrapper/models"
	"webscrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

func AssignmentsDataForACourse(c *colly.Collector, studentId string, mpToView string, courseCode string, courseSection string, courseName string, school string) ([]models.Assignment, error) {
	assignments := make([]models.Assignment, 0)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		rows := dom.Find(".list > tbody>tr")

		rows.Each(func(i int, row *goquery.Selection) {
			if row.Children().Length() == constants.CourseSummaryRowLength && i != 0 {
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

				stuff := basicDataExtractor(row, courseName)
				aAssignment.Category = stuff[constants.CourseSummaryNameCategory]
				aAssignment.Assignment = stuff[constants.CourseSummaryNameAssignment]
				aAssignment.Description = stuff[constants.CourseSummaryNameDescription]
				aAssignment.CourseName = courseName

				tds.Each(func(i int, s *goquery.Selection) {
					if i == constants.CourseSummaryMPIndex {
						aAssignment.MP = strings.TrimSpace(s.Text())
					} else if i == constants.CourseSummaryDueIndex {
						dayname, date := processDueCell(s)
						aAssignment.DayName = strings.TrimSpace(dayname)
						aAssignment.FullDayName = dayClassifier(aAssignment.DayName)
						aAssignment.Date = strings.TrimSpace(date)
						aAssignment.FullDate = aAssignment.Date
					} else if i == constants.CourseSummaryGradeIndex {
						aAssignment.GradeNum, aAssignment.GradePercent = processGradeCellForAssignment(s)
					} else if i == constants.CourseSummaryCommentIndex {
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

	if mpToView == "FG" {
		mpToView = "allMP"
	}

	data := constants.ConstantLinks[school]["assignments"]
	data["studentid"] = studentId
	data["dateRange"] = mpToView
	data["courseAndSection"] = fmt.Sprintf("%v:%v", courseCode, courseSection)
	assignemnts_url, err := utils.FormatDynamicUrl(data, school)

	if err != nil {
		log.Println(err)
		return assignments, err

	}

	err = c.Visit(assignemnts_url)
	if err != nil {
		log.Println("[AssignmentsDataForACourse] : Couldn't visit assignment url: file assignments.go")
		return assignments, err
	}

	c.OnHTMLDetach("body")

	//fmt.Println(assignments)

	return assignments, nil
}

func ScheduleDataForACourse(c *colly.Collector, studentId string, mpToView string, courseCode string, courseSection string, courseName string, school string) ([]models.ScheduleAssignment, error) {
	assignments := make([]models.ScheduleAssignment, 0)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		rows := dom.Find(".list > tbody>tr")
		rows.Each(func(i int, row *goquery.Selection) {
			if row.Children().Length() == constants.CourseSummaryRowLength && i != 0 {
				notGraded := false
				tds := row.Children()
				tds.Each(func(i int, s *goquery.Selection) {

					if i == constants.CourseSummaryGradeIndex {
						_, notGraded = processGradeCellForSchedule(s)

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

					stuff := basicDataExtractor(row, courseName)
					aAssignment.Category = stuff[constants.CourseSummaryNameCategory]
					aAssignment.Assignment = stuff[constants.CourseSummaryNameAssignment]
					aAssignment.Description = stuff[constants.CourseSummaryNameDescription]
					aAssignment.CourseName = courseName

					tds.Each(func(i int, s *goquery.Selection) {
						if i == constants.CourseSummaryDueIndex {
							_, date := processDueCell(s)
							aAssignment.Date = strings.TrimSpace(date)
						} else if i == constants.CourseSummaryGradeIndex {
							aAssignment.Points, _ = processGradeCellForSchedule(s)
						}
					})

					assignments = append(assignments, aAssignment)
				}

			}
		})

	})
	if mpToView == "FG" {
		mpToView = "allMP"
	}
	data := constants.ConstantLinks[school]["assignments"]
	data["studentid"] = studentId
	data["dateRange"] = mpToView
	data["courseAndSection"] = fmt.Sprintf("%v:%v", courseCode, courseSection)
	assignemnts_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		return assignments, err
	}

	//fmt.Printf("assignemnts_url: %v\n", assignemnts_url)
	err = c.Visit(assignemnts_url)
	if err != nil {
		log.Println("[ScheduleDataForACourse]: Couldn't visit assignment url file assignments.go")
		return assignments, err
	}

	c.OnHTMLDetach("body")

	return assignments, err
}

func dayClassifier(day string) string {
	switch day {
	case "Mon":
		return "Monday"
	case "Tue":
		return "Tuesday"
	case "Wed":
		return "Wednesday"
	case "Thu":
		return "Thursday"
	case "Fri":
		return "Friday"
	case "Sat":
		return "Saturday"
	case "Sun":
		return "Sunday"
	default:
		return ""
	}
}

func basicDataExtractor(row *goquery.Selection, courseName string) map[string]string {
	data := map[string]string{}

	data["course_name"] = courseName
	data[constants.CourseSummaryNameAssignment] = ""
	data[constants.CourseSummaryNameDescription] = ""
	data[constants.CourseSummaryNameCategory] = ""

	tds := row.Children()
	tds.Each(func(i int, s *goquery.Selection) {
		if i == constants.CourseSummaryAssignmentNameIndex {
			data[constants.CourseSummaryNameAssignment] = s.Find("b").Text()

			data[constants.CourseSummaryNameDescription] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s.Find("div[style=\"font-style:italic;padding-left:5px;\"]").Text(), data[constants.CourseSummaryNameAssignment], ""), "\n", ""))
			//fmt.Println(data[constants.CourseSummaryNameDescription])

		} else if i == constants.CourseSummaryCategoryIndex {
			nodes := s.Contents()
			for _, node := range nodes.Nodes {
				if node.Type == html.TextNode {
					text := node.Data
					text = strings.ReplaceAll(strings.TrimSpace(text), "\n", "")
					data[constants.CourseSummaryNameCategory] = text
				}
			}

		}
	})

	return data
}

func processGradeCellForAssignment(s *goquery.Selection) (string, string) {

	gradeNum := ""
	gradePercent := ""

	nodes := s.Contents()
	for i, node := range nodes.Nodes {
		if node.Type == html.TextNode {
			// this one is ALWAYS the x / z
			if i%2 == 0 {
				if utils.CleanAString(node.Data) != "" {
					gradeNum = strings.ReplaceAll(utils.CleanAString(node.Data), " ", "")
					//fmt.Printf("gradeNum: %v\n", gradeNum)
				}

			}

		}

	}

	divs := s.Find("div")
	lenDivs := divs.Length()

	if lenDivs == constants.ANotSussyGradeCellDivCount {
		gradePercent = strings.ReplaceAll(utils.CleanAString(divs.Text()), "%", "")

	} else if lenDivs == 2 {
		//sussy eetash cell

		//also like these r the stuff for like excempts and stuff

		/*
			e.g.
			<div>
			<span>
			exempt
			</span>
			</div>
			<div> assign points = 5 </div>
		*/
		divs.Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				x_modifier := utils.CleanAString(s.Find("span").Text())
				gradeNum = x_modifier + " " + gradeNum
			} else if i == 1 {
				// we replace the % and then add it in the front end when neccessary
				gradePercent = strings.ReplaceAll(utils.CleanAString(s.Text()), "%", "")
			}
		})
		if strings.Contains(gradeNum, "Exempt") || strings.Contains(gradeNum, "Absent") || strings.Contains(gradeNum, "Missing") || strings.Contains(gradeNum, "Incomplete") {
			gradeNum = strings.TrimSpace(gradeNum)
			gradePercent = utils.CleanAString(strings.ReplaceAll(strings.ReplaceAll(utils.CleanAString(s.Text()), constants.AssignmentPtsString, ""), gradeNum, ""))

		}

	} else if lenDivs == constants.UngradedCellDivLength {
		subDivs := divs.Find("div")
		// ungraded cell
		subDivs.Each(func(i int, s *goquery.Selection) {
			if i == 1 {

				gradeNum = constants.NotGradedString
				gradePoints := utils.CleanAString(s.Text())

				gradePercent = strings.ReplaceAll(gradePoints, constants.AssignmentPtsString, "")
			}
		})

	}

	return gradeNum, gradePercent
}

func processGradeCellForSchedule(s *goquery.Selection) (string, bool) {
	gradePoints := ""
	notGraded := false

	divs := s.Find("div")
	lenDivs := divs.Length()

	if lenDivs == constants.UngradedCellDivLength {
		subDivs := divs.Find("div")
		if subDivs.Length() == constants.GradeCellThatHasNotGradedSubDivCount {
			// ungraded cell
			//fmt.Println("found un graded cell")
			subDivs.Each(func(i int, s *goquery.Selection) {
				if i == 1 {
					gradePoints = strings.ReplaceAll(utils.CleanAString(s.Text()), constants.AssignmentPtsString, "")
					notGraded = true
				}
			})

		}
	}

	return gradePoints, notGraded
}

func processDueCell(s *goquery.Selection) (dayname string, date string) {
	dayname = ""
	date = ""

	s.Find("div").Each(func(k int, l *goquery.Selection) {
		if k == 0 {
			dayname = l.Text()

		} else if k == 1 {
			date = l.Text()
		}
	})
	return dayname, date
}
