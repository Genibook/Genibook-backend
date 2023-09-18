package constants

var DueIndex = 0
var CourseIndex = 1
var AssignmentIndex = 2
var GradeIndex = 3

var CourseSummaryAssignmentNameIndex = 4
var CourseSummaryRowLength = 5
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

var NormalCell = 1
var EttashCell = 2
var UngradedCell = 3
var GradeCellThatIsJustNormaSubDivCount = 0
var GradeCellThatHasNotGradedSubDivCount = 2

var CounselorTRindex = 1
var AgeTRindex = 2
var BirthdateTRindex = 3
var LunchBalanceTrindex = 4
var LockerTRindex = 5
