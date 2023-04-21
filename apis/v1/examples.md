# Examples
```go
	c := utils.Init_colly()
	utils.DebugLogin(c)
	student := pages.ProfileData(c, 1, constants.MontgomeryHighSchoolKeyName)
	fmt.Printf("student: %v\n", student)
	//pages.GimmeCourseCodes(c, 107604, "MP1")
	// assignments := pages.AssignmentsDataForACourse(c, 107604, "MP2", "22000", "2", "US History I")
	// for _, assignment := range assignments {
	// 	fmt.Println(assignment)
	fmt.Printf("pages.StudentIdAndCurrMP(c): %v\n", pages.StudentIdAndCurrMP(c, constants.MontgomeryHighSchoolKeyName))
	// }
	pages.GradebookData(c, 107604, "MP1", constants.MontgomeryHighSchoolKeyName)
	assignments_more := pages.ScheduleDataForACourse(c, 107604, "MP2", "22000", "2", "US History I")
	for _, assignment := range assignments_more {
		fmt.Println(assignment)

	}
	http.HandleFunc("/view/", makeHandler(viewHandler))
```