package pages

import (
	"encoding/base64"
	"log"
	"strconv"
	"strings"
	"webscrapper/constants"
	"webscrapper/models"
	"webscrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ProfileData(c *colly.Collector, user int, school string) (models.Student, error) {
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
	}

	c.OnHTML("body", func(h *colly.HTMLElement) {
		notecards := h.DOM.Find("table.notecard")
		//fmt.Println(notecards.Length())
		notecards.Each(func(i int, notecard *goquery.Selection) {
			if i == user-1 {

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
									log.Println("profile.go - grade did not convert to int")
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
									log.Println("profile.go - student id did not convert to int")
								}
								student.ID = student_id_int
							case 1:
								id := utils.CleanAString(s.Find("span").Text())
								student_id_int, err := strconv.Atoi(id)
								if err != nil {
									student_id_int = 0
									log.Println("profile.go - student id did not convert to int")
								}
								student.StateID = student_id_int
							}

						})

					}

				})

				//Schedule link stuff.
				schedule_link, exists_schedule_link := notecard.Find("tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > div:nth-child(3) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > a:nth-child(1)").Attr("href")
				if exists_schedule_link {
					student.ScheduleLink = constants.ConstantLinks[school]["base"]["url"] + "/genesis/" + schedule_link
				}

				//right - student stuff
				table = notecard.Find("tbody tr:nth-child(2) > td > table.list > tbody")
				trs = table.Find("tr.listrow")
				trs.Each(func(i int, tr *goquery.Selection) {
					//TODO we can migrate to a map that has keys as the names of the first child in td:nth child
					// then we can like have the values as functions that handles the goquery selection!~!
					switch i {
					case constants.CounselorTRindex:
						student.CounselorName = utils.CleanAString(tr.Find("span[style=\"font-weight: 600;\"]").Text())
					case constants.AgeTRindex:
						age, err := strconv.Atoi(tr.Find("td:nth-child(2)").Text())
						if err != nil {
							age = 0
							log.Println("profile.go - age of student did not convert to int")
						}
						student.Age = age
					case constants.BirthdateTRindex:
						student.Birthday = tr.Find("td:nth-child(2)").Text()
					case constants.LockerTRindex:
						student.Locker = tr.Find("td:nth-child(2)").Text()
					}

				})

			}
		})
	})
	// deteaches this thing, can be used later in functions maybe!

	profile_url, err := utils.FormatUrl("profile", school)
	if err != nil {
		log.Println(err)
		return student, err
	}
	err = c.Visit(profile_url)
	if err != nil {
		log.Println("Couldn't visit profile url, function: ProfileData, file: profile.go")
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

	// fmt.Println(student.ScheduleLink)
	// fmt.Println(student.ImgURL)
	// fmt.Println(student.Image64)
	// fmt.Println(student.ID)
	// fmt.Println(student.Grade)
	// fmt.Println(student.StateID)
	return student, nil
}

func StudentIds(c *colly.Collector, school string) ([]string, error) {
	info := make([]string, 0)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		notecards := h.DOM.Find("table.notecard")
		notecards.Each(func(i int, notecard *goquery.Selection) {
			table := notecard.Find("tbody tr:nth-child(2) > td > table > tbody")
			schedule_td := table.Find("tr > td:nth-child(1)")
			top_half := schedule_td.Find("table.list > tbody > tr")
			top_half.Each(func(i int, tr *goquery.Selection) {
				if i == 2 {
					divs := tr.Find("td > div")
					divs.Each(func(i int, s *goquery.Selection) {
						if i == 0 {
							id := utils.CleanAString(s.Find("span").Text())
							info = append(info, id)
						}
					})
				}
			})
		})

	})
	profile_url, err := utils.FormatUrl("profile", school)
	if err != nil {
		log.Println(err)
		return info, err
	}
	err = c.Visit(profile_url)
	if err != nil {
		log.Println("[StudentIds]: Couldn't visit profile url, file: profile.go")
	}
	c.OnHTMLDetach("body")

	return info, err

}

func WhatGradeIsStudent(c *colly.Collector, school string) ([]int, error) {
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
						return
					}
					grades = append(grades, grade)
				}

			})
		})

	})
	profile_url, err := utils.FormatUrl("profile", school)
	if err != nil {
		log.Println(err)
		return grades, err
	}
	err = c.Visit(profile_url)
	if err != nil {
		log.Println("[WhatGradeIsStudent]: Couldn't visit profile url, file: profile.go")
	}
	c.OnHTMLDetach("body")

	return grades, err
}
