package constants

var DueIndex = 0
var CourseIndex = 1
var AssignmentIndex = 2
var GradeIndex = 3
var RowLength = 5

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
var SubDivCount = 0
var NotGradedSubDivCount = 2

var CounselorTRindex = 1
var AgeTRindex = 2
var BirthdateTRindex = 3
var LunchBalanceTrindex = 4
var LockerTRindex = 5
