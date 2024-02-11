package pages_v2

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	constants "webscrapper/constants/v2"
	"webscrapper/models"
	"webscrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// user - 1 is userselector!
func ProfileData(c *colly.Collector, studentId string, school string) (models.Student, error) {
	exists_an_image := false

	student := models.Student{
		Age:           0,
		ImgURL:        "N/A",
		StateID:       0,
		Birthday:      "N/A",
		ScheduleLink:  "N/A",
		Name:          "N/A",
		Grade:         0,
		Locker:        "N/A",
		CounselorName: "N/A",
		ID:            0,
		Image64:       "N/A",
		LunchBalance:  "N/A",
	}

	c.OnHTML("body", func(h *colly.HTMLElement) {
		notecards := h.DOM.Find("table.notecard")
		//fmt.Println(notecards.Length())
		notecards.Each(func(i int, notecard *goquery.Selection) {

			table := notecard.Find("tbody tr:nth-child(2) > td > table > tbody")
			schedule_td := table.Find("tr > td:nth-child(1)")
			trs := schedule_td.Find("table.list > tbody > tr")
			trs.Each(func(i int, tr *goquery.Selection) {
				switch i {
				case 0:
					tds := tr.Find("td")
					tds.Each(func(k int, td *goquery.Selection) {
						switch k {
						case 0:
							src, exists_img_url := td.Find("img").Attr("src")
							if exists_img_url {
								student.ImgURL = constants.ConstantLinks[school]["base"]["url"] + src
								exists_an_image = true
							}
						case 1:
							firstName := utils.CleanAString(td.Find("div > div > span[style=\"font-weight: 100; color: #001E37\"]").Text())
							lastName := strings.ReplaceAll(utils.CleanAString(strings.Replace(td.Find("div > div").Text(), firstName, "", 1)), " ", "")
							student.Name = firstName + " " + lastName

						case 2:
							grade_int, err := strconv.Atoi(td.Find("span[style=\"font-size: 2em;\"]").Text())
							if err != nil {
								log.Println("[ERROR ProfileData() profile.go]- grade did not convert to int")
								grade_int = 0
							}
							student.Grade = grade_int
						}

					})

				case 2:
					divs := tr.Find("td > div")
					divs.Each(func(i int, s *goquery.Selection) {
						switch i {
						case 0:
							id := utils.CleanAString(s.Find("span").Text())
							student_id_int, err := strconv.Atoi(id)
							if err != nil {
								student_id_int = 0
								log.Println("[ERROR ProfileData() profile.go] - student id did not convert to int")

							}
							student.ID = student_id_int
						case 1:
							id := utils.CleanAString(s.Find("span").Text())
							student_id_int, err := strconv.Atoi(id)
							if err != nil {
								student_id_int = 0
								log.Println("[ERROR ProfileData() profile.go] - state id did not convert to int")
							}
							student.StateID = student_id_int
						}

					})

				}

			})

			//Schedule link stuff.
			schedule_link, exists_schedule_link := notecard.Find("tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(3) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > a:nth-child(1)").Attr("href")
			if exists_schedule_link {
				student.ScheduleLink = constants.ConstantLinks[school]["base"]["url"] + "/genesis/" + schedule_link
			}

			//right - student stuff
			table = notecard.Find("tbody tr:nth-child(2) > td > table.list > tbody")
			trs = table.Find("tr.listrow")
			idxs := map[string]int{
				"counselor":     constants.CounselorTRindex,
				"age":           constants.AgeTRindex,
				"birthdate":     constants.BirthdateTRindex,
				"locker":        constants.LockerTRindex,
				"lunch balance": constants.LunchBalanceTrindex,
			}
			trs.Each(func(i int, tr *goquery.Selection) {
				rowName := strings.ToLower(strings.ReplaceAll(utils.CleanAString(tr.Find("td:nth-child(1)").Text()), ":", ""))
				//fmt.Printf("rowName: %v\n", rowName)

				idxs[strings.ToLower(rowName)] = i

			})
			//fmt.Println(idxs)

			trs.Each(func(i int, tr *goquery.Selection) {

				switch i {
				case idxs["counselor"]:
					student.CounselorName = utils.CleanAString(tr.Find("span[style=\"font-weight: 600;\"]").Text())
				case idxs["age"]:
					age, err := strconv.Atoi(tr.Find("td:nth-child(2)").Text())
					if err != nil {
						age = 0
						log.Println("[ERROR ProfileData() profile.go] - age of student did not convert to int")
					}
					student.Age = age
				case idxs["birthdate"]:
					student.Birthday = tr.Find("td:nth-child(2)").Text()
				case idxs["locker"]:
					student.Locker = tr.Find("td:nth-child(2)").Text()
				case idxs["lunch balance"]:
					student.LunchBalance = strings.ReplaceAll(tr.Find("td:nth-child(2)").Text(), "$", "")
				}

			})

		})
	})

	data := utils.CreateQueryMapCopy("profile", school)
	data["studentid"] = studentId
	profile_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		log.Printf("[ERROR ProfileData() profile.go]: %v\n", err.Error())
		return student, err
	}
	err = c.Visit(profile_url)
	if err != nil {
		fmt.Println("Couldn't visit profile url, function: ProfileData, file: profile.go")
		return student, err
	}
	c.OnHTMLDetach("body")

	c.OnResponse(func(r *colly.Response) {
		test := r.Body
		student.Image64 = base64.StdEncoding.EncodeToString(test)
	})
	if exists_an_image {
		c.Visit(student.ImgURL)
	}
	c.OnResponse(func(r *colly.Response) {})

	return student, nil
}

func StudentIds(c *colly.Collector, school string) ([]string, error) {
	info := make([]string, 0)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		students := h.DOM.Find("div > div.selectStudent > div > ul > li.selectStudentItem")
		//fmt.Println(h.DOM.Html())
		log.Print(students.Length())
		log.Println(" <- ID's length")
		//fmt.Println(students.Length())
		students.Each(func(i int, s *goquery.Selection) {
			//fmt.Println(s.Html())
			val, _ := s.Find("a").Attr("onclick")
			//fmt.Println(val)
			studentId := strings.Split(val, "'")
			if len(studentId) == 3 {
				info = append(info, studentId[1])
			} // } else {
			// 	//fmt.Println(studentId)
			// }
		})
	})

	profile_url, err := utils.FormatUrl("profile", school)
	if err != nil {
		log.Printf("[ERROR StudentIds() profile.go]: %v\n", err.Error())
		return info, err
	}
	err = c.Visit(profile_url)
	if err != nil {
		log.Println("[ERROR StudentIds() profile.go]: Couldn't visit profile url")
	}
	c.OnHTMLDetach("body")

	//fmt.Printf("info: %v\n", info)

	return info, err

}

func WhatGradeIsStudent(c *colly.Collector, school string, studentIds []string) ([]int, error) {
	// we gonna do some trolling
	grades := make([]int, 0)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		notecards := h.DOM.Find("table.notecard")
		//fmt.Println(notecards.Length())
		notecards.Each(func(i int, notecard *goquery.Selection) {
			table := notecard.Find("tbody tr:nth-child(2) > td > table > tbody")
			schedule_td := table.Find("tr > td:nth-child(1)")
			top_half := schedule_td.Find("table.list > tbody > tr")
			top_half.Each(func(i int, tr *goquery.Selection) {
				switch i {
				case 0:
					span := tr.Find("td:nth-child(3) > span[style='font-size: 2em;']")
					grade, err := strconv.Atoi(utils.CleanAString(span.Text()))
					if err != nil {
						grades = append(grades, 0)
					} else {
						grades = append(grades, grade)
					}

				}

			})
		})

	})

	for _, id := range studentIds {
		data := utils.CreateQueryMapCopy("profile", school)
		data["studentid"] = id
		profile_url, err := utils.FormatDynamicUrl(data, school)
		if err != nil {
			log.Printf("[ERROR WhatGradeIsStudent() profile.go]: %v\n", err.Error())
			return grades, err
		}
		err = c.Visit(profile_url)
		if err != nil {
			log.Println("[ERROR WhatGradeIsStudent() profile.go]: Couldn't visit profile url")
		}
	}

	//fmt.Println(grades)

	c.OnHTMLDetach("body")

	return grades, nil
}
