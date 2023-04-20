package utils

import (
	"strings"
	"webscrapper/constants"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func DayClassifier(day string) string {
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

func BasicDataExtractor(row *goquery.Selection, courseName string) map[string]string {
	data := map[string]string{}

	data["course_name"] = courseName
	data[constants.CourseSummaryNameAssignment] = ""
	data[constants.CourseSummaryNameDescription] = ""
	data[constants.CourseSummaryNameCategory] = ""

	tds := row.Children()
	tds.Each(func(i int, s *goquery.Selection) {
		if i == constants.CourseSummaryAssignmentNameIndex {
			data[constants.CourseSummaryNameAssignment] = s.Find("b").Text()

			data[constants.CourseSummaryNameDescription] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s.Find("div[style=\"font-style:italic;padding-left:10px;\"]").Text(), data[constants.CourseSummaryNameAssignment], ""), "\n", ""))
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

func ProcessGradeCellForAssignment(s *goquery.Selection) (string, string) {

	gradeNum := ""
	gradePercent := ""

	nodes := s.Contents()
	for i, node := range nodes.Nodes {
		if node.Type == html.TextNode {
			// this one is ALWAYS the x / z
			if i%2 == 0 {
				if CleanAString(node.Data) != "" {
					gradeNum = strings.ReplaceAll(CleanAString(node.Data), " ", "")
					//fmt.Printf("gradeNum: %v\n", gradeNum)
				}

			}

		}

	}

	divs := s.Find("div")
	lenDivs := divs.Length()

	if lenDivs == constants.ANotSussyGradeCellDivCount {
		subDivs := divs.Find("div")
		if subDivs.Length() == constants.GradeCellThatIsJustNormaSubDivCount {
			//normal cell
			gradePercent = CleanAString(divs.Text())
		} else {
			// ungraded cell
			subDivs.Each(func(i int, s *goquery.Selection) {
				if i == 1 {
					gradeNum = constants.NotGradedString
					gradePoints := CleanAString(s.Text())

					gradePercent = strings.ReplaceAll(gradePoints, constants.AssignmentPtsString, "")
				}
			})

		}

	} else if lenDivs == 2 {
		//sussy eetash cell
		divs.Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				x_modifier := CleanAString(s.Find("span").Text())
				gradeNum = x_modifier + " " + gradeNum
			} else if i == 1 {
				// we replace the % and then add it in the front end when neccessary
				gradePercent = strings.ReplaceAll(CleanAString(s.Text()), "%", "")
			}
		})

	}

	return gradeNum, gradePercent
}

func ProcessGradeCellForSchedule(s *goquery.Selection) (string, bool) {
	gradePoints := ""
	notGraded := false

	divs := s.Find("div")
	lenDivs := divs.Length()

	if lenDivs == constants.ANotSussyGradeCellDivCount {
		subDivs := divs.Find("div")
		if subDivs.Length() == constants.GradeCellThatHasNotGradedSubDivCount {
			// ungraded cell
			subDivs.Each(func(i int, s *goquery.Selection) {
				if i == 1 {
					gradePoints = strings.ReplaceAll(CleanAString(s.Text()), constants.AssignmentPtsString, "")
					notGraded = true
				}
			})

		}

	}

	return gradePoints, notGraded
}

func ProcessDueCell(s *goquery.Selection) (dayname string, date string) {
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

// data = {
// 	"course_name": course_namee,
// 	"date": date,
// 	"points": grade_points,
// 	constants.CourseSummaryNameCategory: category,
// 	constants.CourseSummaryNameAssignment: assignment,
// 	constants.CourseSummaryNameDescription: description,
// }
