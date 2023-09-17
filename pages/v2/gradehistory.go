package pages_v2

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
					// 4 tds in a row now...
					// case 0 now has grade and school year.
					// 1 has course name
					// 2 has FG
					// 3 has attempted and earned

					case 0:
						course.SchoolYear = utils.CleanAString(td.Find("div[style=\"padding-left: 5px;\"]").Text())
						gradeString := utils.CleanAString(strings.ReplaceAll(td.Find("div").First().Text(), "Grade ", ""))
						grade, err := strconv.Atoi(gradeString)
						//fmt.Println(gradeString)
						if err != nil {
							return
						}
						course.Grade = grade
					case 1:
						course.Name = strings.TrimSpace(utils.CleanAString(td.Find("div[style=\"padding-left: 5px;font-weight: bold;\"]").Text()))
						course.School = utils.CleanAString(td.Find("div").First().Text())
						//fmt.Println(course.School)
					case 2:
						//FG
						course.FG = utils.CleanAString(td.Text())
					case 3:

						divs := td.Find("div")
						divs.Each(func(idx int, div *goquery.Selection) {
							switch idx {
							case 0:
								atmString := utils.CleanAString(strings.ReplaceAll(div.Text(), "Attempted", ""))
								//fmt.Printf("atmString: %v\n", atmString)
								attempted, err := strconv.ParseFloat(atmString, 64)
								if err != nil {
									return
								}
								course.Attempted = float32(attempted)
							case 1:
								eString := utils.CleanAString(strings.ReplaceAll(div.Text(), "Earned", ""))
								//fmt.Println(eString)
								earned, err := strconv.ParseFloat(eString, 64)
								if err != nil {
									return
								}
								course.Earned = float32(earned)
							}
						})

						// case 0:

						// 	course.SchoolYear = utils.CleanAString(td.Text())
						// case 1:
						// 	grade, err := strconv.Atoi(utils.CleanAString(td.Text()))
						// 	if err != nil {
						// 		return
						// 	}
						// 	course.Grade = grade
						// case 2:
						// 	course.Name = strings.TrimSpace(utils.CleanAString(td.Text()))
						// case 3:
						// 	course.School = utils.CleanAString(td.Text())
						// case 4:
						// 	course.FG = utils.CleanAString(td.Text())
						// case 5:
						// 	attempted, err := strconv.ParseFloat(utils.CleanAString(td.Text()), 64)
						// 	if err != nil {
						// 		return
						// 	}
						// 	course.Attempted = float32(attempted)
						// case 6:
						// 	earned, err := strconv.ParseFloat(utils.CleanAString(td.Text()), 64)
						// 	if err != nil {
						// 		return
						// 	}
						// 	course.Earned = float32(earned)
					}
				})

				if course.SchoolYear != "" {
					courses[course.SchoolYear] = append(courses[course.SchoolYear], map[string]models.Course{course.Name: course})
					//fmt.Println(course.SchoolYear)
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

	//TODO: change this shit what the tf
	currCourses := map[string]string{}
	c.OnHTML("body", func(h *colly.HTMLElement) {
		ATTCREDITSINDEX := 6
		dom := h.DOM
		table := dom.Find(".list")
		th := table.Find("tbody>tr.listheading")
		th.Each(func(i int, s *goquery.Selection) {
			s.Find("td").Each(func(k int, td *goquery.Selection) {
				text := utils.CleanAString(td.Text())
				if strings.ToLower(text) == "att." {
					ATTCREDITSINDEX = k
				}
			})
		})
		trs := table.Find("tbody>tr.listroweven, tbody>tr.listrowodd")
		trs.Each(func(i int, s *goquery.Selection) {
			name := ""
			//att is something sus on one of the columns, it is a name, it is supposed to be like the given credits

			att := ""

			s.Find("td").Each(func(k int, td *goquery.Selection) {
				switch k {
				case 0:
					name = utils.CleanAString(td.Text())
				case ATTCREDITSINDEX:
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
