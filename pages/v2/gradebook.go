package pages_v2

import (
	"log"
	"strconv"
	"strings"
	constants "webscrapper/constants/v2"
	"webscrapper/models"
	"webscrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func GradebookData(c *colly.Collector, studentId string, mpToView string, school string) (map[string]models.OneGrade, error) {

	grades := map[string]models.OneGrade{}

	c.OnHTML("body", func(h *colly.HTMLElement) {
		items := h.DOM.Find("table.notecard > tbody > tr:nth-child(2) > td > div.itemContainer > div.twoColFlex")
		items.Each(func(i int, item *goquery.Selection) {
			aGrade := models.OneGrade{
				Grade:        0,
				TeacherName:  "",
				TeacherEmail: "",
			}
			L := item.Find("div.gradebookGrid")
			R := item.Find("div.twoColGridItem")
			utils.Assert((R.Length()+L.Length()) == 2, "GradeBookDataV2 gradebook.go")

			grade := L.Find("span[style=\"font-size:20pt\"]")
			class := R.Find("span.categorytab")
			name := R.Find("div[style=\"font-size: 12pt\"]")
			email, _ := name.Find("a").Attr("href")

			float_grade, _ := strconv.ParseFloat(utils.CleanAString(strings.ReplaceAll(grade.Text(), "%", "")), 32)
			aGrade.Grade = float_grade

			courseName := utils.CleanAString(class.Text())
			aGrade.TeacherName = utils.CleanAString(name.Text())
			aGrade.TeacherEmail = utils.CleanAString(utils.CleanAString(strings.ReplaceAll(string(email), "mailto:", "")))
			grades[courseName] = aGrade
		})
		// dom := h.DOM
		// rows := dom.Find(".list > tbody>tr")

		// rows.Each(func(i int, row *goquery.Selection) {

		// 	if i != 0 {
		// 		aGrade := models.OneGrade{
		// 			Grade:        0,
		// 			TeacherName:  "",
		// 			TeacherEmail: "",
		// 		}

		// 		courseName := fmt.Sprintf("class %d", i)
		// 		tds := row.Children()
		// 		tds.Each(func(k int, s *goquery.Selection) {
		// 			if k == 0 {
		// 				courseName = strings.TrimSpace(strings.ReplaceAll(s.Text(), "\n", ""))

		// 				//fmt.Println(courseName)

		// 			} else if k == 1 {
		// 				aName := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s.Text(), "Email:", ""), "\n", ""))

		// 				//fmt.Println(aName)
		// 				aGrade.TeacherName = aName
		// 				href, found := s.Find("a").Attr("href")
		// 				if !found {
		// 					log.Println("[gradebook.go] - teacher email not found")
		// 					href = ""
		// 				}
		// 				aGrade.TeacherEmail = strings.ReplaceAll(string(href), "mailto:", "")
		// 				//fmt.Println(aGrade.TeacherEmail)
		// 			} else if k == 2 {
		// 				grade := strings.TrimSpace(strings.ReplaceAll(s.Find("tbody>tr>td:nth-child(1)").Text(), "%", ""))
		// 				grade = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(grade, "\n", ""), "*PROJECTED", ""))

		// 				float_grade, err := strconv.ParseFloat(grade, 32)
		// 				if err != nil {
		// 					log.Println("[gradebook.go] - couldn't convert grade string to float grade")
		// 					log.Println(err)
		// 					float_grade = 0.0
		// 				}
		// 				aGrade.Grade = float64(float_grade)
		// 			}
		// 		})
		// 		grades[courseName] = aGrade

		// 	}

		// })
	})

	data := constants.ConstantLinks[school]["gradebook"]
	data["studentid"] = studentId
	data["mpToView"] = mpToView
	gradebook_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		log.Println(err)
		return grades, err

	}

	err = c.Visit(gradebook_url)
	if err != nil {
		log.Println("Couldn't visit gradebook url: function GradebookData, file gradebook.go")
		return grades, err
	}
	c.OnHTMLDetach("body")

	return grades, nil

}

func GimmeCourseCodes(c *colly.Collector, studentId string, mpToView string, school string) (map[string]map[string]string, error) {
	courseCodes := map[string]map[string]string{}

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM

		//this should be equal to the amount of courses * 2
		divs := dom.Find("div.itemContainer > div.twoColFlex > div")

		divs.Each(func(i int, div *goquery.Selection) {
			// even odd checking for the goToCourseSummary function and showAssignmentsByMpandCourse function
			switch i % 2 {
			case 1:
				courseName := utils.CleanAString(div.Find("div > span").Text())
				onclick, err := div.Attr("onclick")
				if !err {
					log.Printf("Course on index %d does not have onclick attr\n", i)
					onclick = ""
				}
				onclick = strings.ReplaceAll(strings.ReplaceAll(onclick, "goToCourseSummary(", ""), ");", "")

				if onclick != "" {
					data := strings.Split(onclick, ",")
					courseCodes[courseName] = map[string]string{"code": strings.ReplaceAll(data[0], "'", ""), "section": strings.ReplaceAll(data[1], "'", "")}
				}
			}
		})

	})

	data := constants.ConstantLinks[school]["gradebook"]
	data["studentid"] = studentId
	data["mpToView"] = mpToView
	gradebook_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		log.Println(err)
		return courseCodes, err

	}

	err = c.Visit(gradebook_url)
	if err != nil {
		log.Println("Couldn't visit gradebook url: function gimmeCourseCodes,  file gradebook.go")
		return courseCodes, err
	}
	c.OnHTMLDetach("body")

	return courseCodes, nil
}

func GimmeMPs(c *colly.Collector, studentId string, school string) ([]string, error) {
	mps := make([]string, 0)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		dom := h.DOM
		selector := dom.Find("select.fieldvalue:nth-child(1)")
		selector.Children().Each(func(i int, s *goquery.Selection) {
			mp := utils.CleanAString(s.Text())
			mps = append(mps, mp)
		})
	})

	data := constants.ConstantLinks[school]["gradebook"]
	data["studentid"] = studentId
	data["mpToView"] = "MP1"

	gradebook_url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		log.Println(err)
		return mps, err
	}

	err = c.Visit(gradebook_url)
	if err != nil {
		log.Println("Couldn't visit gradebook url: function GimmeMps,  file gradebook.go")
		return mps, err
	}
	c.OnHTMLDetach("body")

	return mps, nil
}

// final grades = Grades.fromJson({
//   'Math': {
//     'grade': 85.0,
//     'teacher_name': 'John Smith',
//     'teacher_email': 'john.smith@example.com'
//   },
//   'English': {
//     'grade': 92.0,
//     'teacher_name': 'Jane Doe',
//     'teacher_email': 'jane.doe@example.com'
//   },
//   'Science': {
//     'grade': 78.0,
//     'teacher_name': 'Bob Johnson',
//     'teacher_email': 'bob.johnson@example.com'
//   },
// });
