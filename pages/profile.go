package pages

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"webscrapper/constants"
	"webscrapper/models"
	"webscrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// return Student(
//
//		age: json['age'],
//		imageUrl: json['img_url'],
//		stateId: json['state_id'],
//		birthday: json['birthday'],
//		scheduleLink: json['schedule_link'],
//		name: json['name'],
//		grade: json['grade'],
//		locker: json['locker'],
//		counselorName: json['counselor_name'],
//		id: json['id'],
//		image64: json['image64'],
//		assignments: json['assignments'] != null
//			? Assignments.fromJson(json['assignments'])
//			: Assignments.fromJson({}),
//		grades: json['grades'] != null
//			? Grades.fromJson(json['grades'])
//			: Grades.fromJson({}),
//	  );
func ProfileData(c *colly.Collector) {
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

	c.OnHTML("table.notecard", func(e *colly.HTMLElement) {
		// table := e.DOM.Find("tbody:nth-child(1) > table:nth-child(1)")
		table := e.DOM.Find("tbody table > tbody")

		schedule_td := table.Find("tr[valign=\"top\"] > td[valign=\"top\"]:nth-child(2)")
		student_demo_and_whereabouts := table.Find("tr[valign=\"top\"] > td[valign=\"top\"]:nth-child(1)")

		//right - student stuff
		src, exists_img_url := student_demo_and_whereabouts.Find("img").Attr("src")
		if exists_img_url {
			student.ImgURL = constants.ConstantLinks["base"]["url"] + src
			exists_an_image = true
		}

		//left - schedule stuff
		top_half := schedule_td.Find("table.list > tbody > tr")

		top_half.Each(func(i int, tr *goquery.Selection) {
			if i == 0 {
				firstName := strings.ReplaceAll(tr.Find("td > span[style=\"font-weight: 100; color: #001E37\"]").Text(), "\n", "")
				lastName := strings.ReplaceAll(strings.ReplaceAll(strings.Replace(tr.Find("td[style=\"font-size: 1.5em;\"]").Text(), firstName, "", 1), "\n", ""), " ", "")
				student.Name = firstName + " " + lastName

				grade_int, err := strconv.Atoi(tr.Find("td > span[style=\"font-size: 2em;\"]").Text())
				if err != nil {
					log.Fatal("profile.go - grade for some reason did not convert to int")
				}
				student.Grade = grade_int

			} else if i == 1 {
				spans := tr.Find("td > span")
				spans.Each(func(i int, s *goquery.Selection) {
					if i == 0 {
						student_id_int, err := strconv.Atoi(s.Text())
						if err != nil {
							student_id_int = 0
							log.Fatal("profile.go - student id for some reason did not convert to int")
						}
						student.ID = student_id_int

					} else if i == 1 {
						state_ID, err := strconv.Atoi(s.Text())
						if err != nil {
							state_ID = 0
							log.Fatal("profile.go - student id for some reason did not convert to int")
						}
						student.StateID = state_ID
					}
				})
			}

		})

		schedule_link, exists_schedule_link := e.DOM.Find(".notecard > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > div:nth-child(3) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > a:nth-child(1)").Attr("href")
		if exists_schedule_link {
			student.ScheduleLink = constants.ConstantLinks["base"]["url"] + "/genesis/" + schedule_link
		}

	})
	// deteaches this thing, can be used later in functions maybe!

	str, err := utils.FormatUrl("profile")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Visit(str)
	c.OnHTMLDetach("table.notecard")

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
	//fmt.Println(student.Image64)
	// fmt.Println(student.ID)
	//fmt.Println(student.Grade)
	// fmt.Println(student.StateID)
}
