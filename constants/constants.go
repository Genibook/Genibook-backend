package constants

var ConstantLinks = map[string]map[string]string{
	"base":     {"url": "https://parents.mtsd.k12.nj.us"},
	"endpoint": {"url": "https://parents.mtsd.k12.nj.us/genesis/parents?"},
	"home":     {"url": "https://parents.mtsd.k12.nj.us/genesis/sis/view?gohome=true"},
	"profile": {
		"tab1":   "studentdata",
		"tab2":   "studentsummary",
		"action": "form",
	},
	"login": {
		"url":      "https://parents.mtsd.k12.nj.us/genesis/sis/j_security_check",
		"username": "j_username",
		"password": "j_password",
	},
	"gradebook": {
		"tab1":      "studentdata",
		"tab2":      "gradebook",
		"tab3":      "weeklysummary",
		"studentid": "",
		"action":    "form",
		"mpToView":  "",
	},
	"assignments": {
		"tab1":          "studentdata",
		"tab2":          "gradebook",
		"tab3":          "coursesummary",
		"studentid":     "",
		"action":        "form",
		"mpToView":      "",
		"courseCode":    "",
		"courseSection": "",
	},
}

var CourseSummaryAssignmentNameIndex = 4
var CourseSummaryRowLength = 9
var CourseSummaryCategoryIndex = 3
var CourseSummaryMPIndex = 0
var CourseSummaryDueIndex = 1
var CourseSummaryTeacherIndex = 2
var CourseSummaryGradeIndex = 5
var CourseSummaryCommentIndex = 6
var CourseSummaryPrevIndex = 7
var CourseSummaryDocsIndex = 8

//Cell types

/*
Goofy eetash cell

<div>
x0.5
</div>
10/10
<div>
100.0%
<div>

Normal Cell
10/10
<div>
100.0%
</div>

"ungraded" cell
<div>

	<subdiv1>
	not graded
	</subdiv1>
	<subdiv2>
	assignment points: 2
	</subdiv2>

<div>
*/
var SussyEetashGradeCellDivCount = 2
var ANotSussyGradeCellDivCount = 1
var GradeCellThatHasNotGradedSubDivCount = 2
var GradeCellThatIsJustNormaSubDivCount = 0

var CourseSummaryNameAssignment = "assignment"
var CourseSummaryNameDescription = "description"
var CourseSummaryNameCategory = "category"

var NotGradedString = "Not Graded"
var AssignmentPtsString = "Assignment Pts: "

// In Go, you can create a dictionary (also called a map) using the built-in make function. Here's an example:

// go

// // Create an empty map with string keys and int values
// m := make(map[string]int)

// // Add some key-value pairs to the map
// m["apple"] = 1
// m["banana"] = 2
// m["cherry"] = 3

// // Access a value by its key
// fmt.Println("Value of apple:", m["apple"])

// // Check if a key exists in the map
// if _, ok := m["banana"]; ok {
//     fmt.Println("Banana is in the map")
// }

// // Delete a key-value pair from the map
// delete(m, "cherry")
