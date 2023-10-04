package pages_v2

import (
	"fmt"
	"log"
	"strings"
	constants "webscrapper/constants/v2"
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
		rows := dom.Find("table.notecard > tbody > tr:nth-child(2) > td > div.desktop>table.list > tbody>tr")

		rows.Each(func(i int, row *goquery.Selection) {
			if row.Children().Length() == constants.RowLength && i != 0 {
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
				//fmt.Println(tds.Length())

				tds.Each(func(i int, td *goquery.Selection) {
					switch i {
					case constants.DueIndex:
						divs := td.Find("div")
						utils.Assert((divs.Length()) == 2, "AssignmentsData constants.DueIndex pages/v2/assignments.go")
						divs.Each(func(i int, div *goquery.Selection) {
							if i == 0 {
								aAssignment.DayName = utils.CleanAString(div.Text())
								aAssignment.FullDayName = dayClassifier(aAssignment.DayName)
							} else if i == 1 {
								aAssignment.Date = utils.CleanAString(div.Text())
								//TODO: um like actually do this 9/17/2023
								aAssignment.FullDate = aAssignment.Date
							}
						})
					case constants.CourseIndex:
						divs := td.Find("div")
						utils.Assert((divs.Length()) == 2, "AssignmentsData constants.CourseIndex pages/v2/assignments.go")
						divs.Each(func(i int, div *goquery.Selection) {
							if i == 0 {
								aAssignment.CourseName = utils.CleanAString(div.Text())
							} else if i == 1 {
								aAssignment.Teacher = utils.CleanAString(div.Text())
							}

						})
					case constants.AssignmentIndex:
						//fmt.Println("hi")
						divs := td.Find("div")
						aAssignment.Assignment = utils.CleanAString(td.Find("b").Text())
						divs.Each(func(i int, div *goquery.Selection) {
							if i == 0 {
								aAssignment.Category = utils.CleanAString(div.Text())
							} else if i == 1 {

								if !strings.Contains(div.Text(), "Close") {
									aAssignment.Description = utils.CleanAString(div.Text())
								}
							}
						})
					case constants.GradeIndex:
						aAssignment.GradeNum, aAssignment.GradePercent = processGradeCellForAssignment(td)
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
	//FIXME: 10/3/2023 I JUST DISCOVERED BIGGEST ERROR - I AM NOT COPYING ANOTHER DICT IN MEMORY IM ACTUALLY JUST USING THE DICT ITSELF AND CHANGING VALUES BIG ERROR
	data["courseAndSection"] = fmt.Sprintf("%v:%v", courseCode, courseSection)
	data["status"] = ""
	assignemnts_url, err := utils.FormatDynamicUrl(data, school)
	//fmt.Printf("assignemnts_url: %v\n", assignemnts_url)

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

func ScheduleDataForACourse(c *colly.Collector, studentId string, mpToView string, school string) (map[string][]models.ScheduleAssignment, error) {
	assignments := map[string][]models.ScheduleAssignment{}
	//assignments := make([]models.ScheduleAssignment, 0)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		rows := dom.Find("table.notecard > tbody > tr:nth-child(2) > td > div.desktop>table.list > tbody>tr")
		rows.Each(func(i int, row *goquery.Selection) {
			if row.Children().Length() == constants.RowLength && i != 0 {
				notGraded := false
				tds := row.Children()
				tds.Each(func(i int, td *goquery.Selection) {

					if i == constants.GradeIndex {
						_, notGraded = processGradeCellForSchedule(td)
						//fmt.Println(notGraded)
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

					tds := row.Children()

					tds.Each(func(i int, td *goquery.Selection) {
						switch i {
						case constants.DueIndex:
							divs := td.Find("div")
							utils.Assert((divs.Length()) == 2, "ScheduleDataForACourse constants.DueIndex pages/v2/assignments.go")
							divs.Each(func(i int, div *goquery.Selection) {
								if i == 1 {
									aAssignment.Date = utils.CleanAString(div.Text())

								}
							})
						case constants.CourseIndex:
							divs := td.Find("div")
							utils.Assert((divs.Length()) == 2, "ScheduleDataForACourse constants.CourseIndex pages/v2/assignments.go")
							divs.Each(func(i int, div *goquery.Selection) {
								if i == 0 {
									aAssignment.CourseName = utils.CleanAString(div.Text())
								}

							})
						case constants.AssignmentIndex:
							divs := td.Find("div")
							aAssignment.Assignment = utils.CleanAString(td.Find("b").Text())
							divs.Each(func(i int, div *goquery.Selection) {
								if i == 0 {
									aAssignment.Category = utils.CleanAString(div.Text())
								} else if i == 1 {
									if !strings.Contains(div.Text(), "Close") {
										aAssignment.Description = utils.CleanAString(div.Text())
									}

								}
							})
						case constants.GradeIndex:
							aAssignment.Points, _ = processGradeCellForSchedule(td)

						}
					})
					//fmt.Printf("aAssignment.CourseName: %v\n", aAssignment.CourseName)
					courseassigns, ok := assignments[aAssignment.CourseName]
					if ok {
						assignments[aAssignment.CourseName] = append(courseassigns, aAssignment)
					} else {
						assignments[aAssignment.CourseName] = make([]models.ScheduleAssignment, 0)
						assignments[aAssignment.CourseName] = append(courseassigns, aAssignment)
					}

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
	data["status"] = "UNGRADED"
	data["courseAndSection"] = ""
	assignemnts_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		return assignments, err
	}
	//fmt.Printf("SCHEDULE URL: %v\n", assignemnts_url)

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

func processGradeCellForAssignment(s *goquery.Selection) (string, string) {

	gradeNum := ""
	gradePercent := ""

	nodes := s.Contents()
	//fmt.Println(nodes)
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

	if lenDivs == constants.NormalCell {
		gradePercent = strings.ReplaceAll(utils.CleanAString(divs.Text()), "%", "")

	} else if lenDivs == constants.EttashCell {
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

	} else if lenDivs == constants.UngradedCell {
		//len divs is recursive which means it finds all divs in the subtree in our case it is going tobe 3 divs.
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

func processGradeCellForSchedule(td *goquery.Selection) (string, bool) {
	gradePoints := ""
	notGraded := false

	divs := td.Find("div")
	lenDivs := divs.Length()
	//fmt.Println(divs.Text())
	//fmt.Printf("lenDivs: %v\n", lenDivs)

	if lenDivs == constants.UngradedCell {
		subDivs := divs.Find("div")
		if subDivs.Length() == constants.NotGradedSubDivCount {
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
