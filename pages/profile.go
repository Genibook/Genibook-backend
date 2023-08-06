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
	counselorTRindex := 2
	ageTRindex := 4
	birthdateTRindex := 5
	lockerTRindex := 6

	c.OnHTML("body", func(h *colly.HTMLElement) {
		notecards := h.DOM.Find("table.notecard")
		//fmt.Println(notecards.Length())
		notecards.Each(func(i int, notecard *goquery.Selection) {
			if i == user-1 {
				//TODO: update format issues
				table := notecard.Find("tbody table > tbody")
				//TODO: so go to the profile page, rn its kidna like
				/*
					table := notecard.Find("tbody tr:nth-child(2) > td > table > tbody")

					using this, we can see that until tbdoy theres a tr and then comments. My old code kinda gets the second and first child which is
					 correct based on the comments
					rn the thing i sthat there is only the "schedule" section, and that is the second nth child

					theres also no more valigns so now its kinda like:

					schedule_td :=  table.Find("tr > td:nth-child(1)")
					student_demo_and_whereabouts := table.Find("tr> td > table.list > tbody")

					we'll see at the start of the school year what they changed 💀

				*/
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
					if i == counselorTRindex {
						student.CounselorName = strings.Replace(strings.ReplaceAll(strings.ReplaceAll(tr.Find("span[style=\"font-weight: 600;\"]").Text(), "\n", ""), " ", ""), ",", ", ", 1)
					} else if i == ageTRindex {
						age, err := strconv.Atoi(tr.Find("td:nth-child(2)").Text())
						if err != nil {
							age = 0
							log.Println("profile.go - age of student did not convert to int")
						}
						student.Age = age
					} else if i == birthdateTRindex {
						student.Birthday = tr.Find("td:nth-child(2)").Text()
					} else if i == lockerTRindex {
						student.Locker = tr.Find("td:nth-child(2)").Text()
					}
				})

				//left - schedule stuff
				// the top half in the html code is correct, we are finding the first one, so i'll just use .Find not findAll()
				top_half := schedule_td.Find("table.list > tbody > tr")

				top_half.Each(func(i int, tr *goquery.Selection) {
					// TODO: update thses cells.
					/*
							It's like three <tr></trs>s now.

							the first one has three tds, and no longer spans, and then the first td is the profile picture (yeah so ig no longer
							in the student demo row/col? (is that even a thing now?)),
						 	second is name and barcode?, and third is the grade thingy

							first tr will look smth like:

							tds = tr.Find("td") // we can use this cuz i chedk the html and nothing rly appears to have td othe rhtan the children
							tds.Each(func(k int, td *goquery.Selection){
								switch k{

								}case 0{
									// image getter code
								}case 1{
									//code to get name
									// first name -> "td > div > div > span[style=\"font-weight: 100; color: #001E37\"]" text
									// last name = td > div > div " text replace firstname

								}case 2{
									"td > span[style=\"font-size: 2em;\"]"
									// same thing in the current code with i== 0 grade_int thingy
								}
							})


							second tr is blank

							third tr is student id and state id

							strcutre tree:
							td
								div
									student id
										span style font-weight bold: literal id
								div
									state id
										span style font-weight bold: literal id


					*/

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

// TODO update this too...
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
