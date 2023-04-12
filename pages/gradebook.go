package pages

import (
	"log"
	"strconv"
	"webscrapper/constants"
	"webscrapper/utils"

	"github.com/gocolly/colly"
)

func GradebookData(c *colly.Collector, studentId int, mpToView string) {

	data := constants.ConstantLinks["gradebook"]
	data["studentid"] = strconv.Itoa(studentId)
	data["mpToView"] = mpToView
	gradebook_url, err := utils.FormatDynamicUrl(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	c.Visit(gradebook_url)

}

//'grades': {
// 'Math': {
// 	'grade': 85.0,
// 	'teacher_name': 'John Smith',
// 	'teacher_email': 'john.smith@example.com'
//   },
//   'English': {
// 	'grade': 92.0,
// 	'teacher_name': 'Jane Doe',
// 	'teacher_email': 'jane.doe@example.com'
//   },
//   'Science': {
// 	'grade': 78.0,
// 	'teacher_name': 'Bob Johnson',
// 	'teacher_email': 'bob.johnson@example.com'
//   },
// }
