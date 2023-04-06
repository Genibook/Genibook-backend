package pages

import (
	"fmt"
	"os"
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

	fmt.Println("yoo getting profile data")

	c.OnHTML("table.notecard", func(e *colly.HTMLElement) {
		// table := e.DOM.Find("tbody:nth-child(1) > table:nth-child(1)")
		table := e.DOM.Find("tbody table > tbody")
		html, _ := table.Html()
		os.WriteFile("table.html", []byte(html), 0600)

		schedule_td := table.Find("tr[valign=\"top\"] > td[valign=\"top\"]:nth-child(2)")

		student_demo_and_whereabouts := table.Find("tr[valign=\"top\"] > td[valign=\"top\"]:nth-child(1)")

		student_demo_and_whereabouts.Find("img").Each(func(i int, s *goquery.Selection) {
			src, exists := s.Attr("src")
			if exists {
				fmt.Println("Image source:", src)
			}

		})

		schedule_td = schedule_td.Find("table.list > tbody")

	})
	// deteaches this thing, can be used later in functions maybe!
	defer c.OnHTMLDetach("table.notecard")

	str, err := utils.FormatUrl("profile")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Visit(str)

}
