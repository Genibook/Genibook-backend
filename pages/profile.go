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
				table := notecard.Find("tbody table > tbody")

				schedule_td := table.Find("tr[valign=\"top\"] > td[valign=\"top\"]:nth-child(2)")
				student_demo_and_whereabouts := table.Find("tr[valign=\"top\"] > td[valign=\"top\"]:nth-child(1) > table.list > tbody")

				//right - student stuff
				src, exists_img_url := student_demo_and_whereabouts.Find("img").Attr("src")
				if exists_img_url {
					student.ImgURL = constants.ConstantLinks[school]["base"]["url"] + src
					exists_an_image = true
				}

				trs := student_demo_and_whereabouts.Find("tr.listroweven")
				trs.Each(func(i int, tr *goquery.Selection) {
					if i == 2 {
						student.CounselorName = strings.Replace(strings.ReplaceAll(strings.ReplaceAll(tr.Find("span[style=\"font-weight: 600;\"]").Text(), "\n", ""), " ", ""), ",", ", ", 1)
					} else if i == trs.Length()-3 {
						age, err := strconv.Atoi(tr.Find("td:nth-child(2)").Text())
						if err != nil {
							age = 0
							log.Println("profile.go - age of student did not convert to int")
						}
						student.Age = age
					} else if i == trs.Length()-2 {
						student.Birthday = tr.Find("td:nth-child(2)").Text()
					} else if i == trs.Length()-1 {
						student.Locker = tr.Find("td:nth-child(2)").Text()
					}
				})

				//left - schedule stuff
				top_half := schedule_td.Find("table.list > tbody > tr")

				top_half.Each(func(i int, tr *goquery.Selection) {
					if i == 0 {
						firstName := strings.ReplaceAll(tr.Find("td > span[style=\"font-weight: 100; color: #001E37\"]").Text(), "\n", "")
						lastName := strings.ReplaceAll(strings.ReplaceAll(strings.Replace(tr.Find("td[style=\"font-size: 1.5em;\"]").Text(), firstName, "", 1), "\n", ""), " ", "")
						student.Name = firstName + " " + lastName

						grade_int, err := strconv.Atoi(tr.Find("td > span[style=\"font-size: 2em;\"]").Text())
						if err != nil {
							log.Println("profile.go - grade did not convert to int")
						}
						student.Grade = grade_int
					} else if i == 1 {
						spans := tr.Find("td > span")
						spans.Each(func(i int, s *goquery.Selection) {
							if i == 0 {
								student_id_int, err := strconv.Atoi(s.Text())
								if err != nil {
									student_id_int = 0
									log.Println("profile.go - student id did not convert to int")
								}
								student.ID = student_id_int

							} else if i == 1 {
								state_ID, err := strconv.Atoi(s.Text())
								if err != nil {
									state_ID = 0
									log.Println("profile.go - state id did not convert to int")
								}
								student.StateID = state_ID
							}
						})
					}

				})

				schedule_link, exists_schedule_link := notecard.Find("tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > div:nth-child(3) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > a:nth-child(1)").Attr("href")
				if exists_schedule_link {
					student.ScheduleLink = constants.ConstantLinks[school]["base"]["url"] + "/genesis/" + schedule_link
				}
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
			table := notecard.Find("tbody table > tbody")
			schedule_td := table.Find("tr[valign=\"top\"] > td[valign=\"top\"]:nth-child(2)")
			top_half := schedule_td.Find("table.list > tbody > tr")
			top_half.Each(func(i int, tr *goquery.Selection) {
				if i == 1 {
					spans := tr.Find("td > span")
					spans.Each(func(i int, s *goquery.Selection) {
						if i == 0 {
							id := utils.CleanAString(s.Text())
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
		log.Println("Couldn't visit profile url, function: StudentIdAndCurrMP, file: profile.go")
	}
	c.OnHTMLDetach("body")

	return info, err

}
