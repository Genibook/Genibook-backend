package api_v1

import (
	"fmt"
	"log"
	"net/http"
	"webscrapper/constants"
	"webscrapper/models"
	"webscrapper/pages"
	"webscrapper/utils"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

func GetIDs(userSelector int, c *colly.Collector, highSchool string, w http.ResponseWriter) ([]string, error) {
	IDS, err := pages.StudentIds(c, highSchool)
	if err != nil {
		return make([]string, 0), err
	}
	if userSelector > len(IDS) {
		log.Printf("User selector index > len(available IDS) Length: %d\n", len(IDS))
		http.Error(w, fmt.Sprintf("User selector index > len(available IDS) Length: %d", len(IDS)), http.StatusNotAcceptable)
		return make([]string, 0), http.ErrBodyNotAllowed
	}
	return IDS, nil
}

func GetProfile(w http.ResponseWriter, functionName string, email string, password string, highSchool string, userSelector int) (models.Student, error) {
	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)

	if e != nil {
		return models.Student{
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
		}, e
	}
	profile, err := pages.ProfileData(c, userSelector, highSchool)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func GetGrade(w http.ResponseWriter, functionName string, email string, password string, highSchool string, userSelector int) (int, error) {

	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)

	grades, e := pages.WhatGradeIsStudent(c, highSchool)

	if e != nil {
		return 0, e
	}

	return grades[userSelector-1], nil
}

func GetListOfStudentGrade(w http.ResponseWriter, functionName string, email string, password string, highSchool string, userSelector int) ([]int, error) {
	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)

	grades, e := pages.WhatGradeIsStudent(c, highSchool)

	if e != nil {
		return make([]int, 0), e
	}

	return grades, nil
}

func GetMPs(w http.ResponseWriter, r *http.Request, functionName string, email string, password string, highSchool string, userSelector int) ([]string, error) {
	//GimmeMPs
	mps := make([]string, 0)

	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)

	IDS, err := GetIDs(userSelector, c, highSchool, w)
	if err != nil {
		return mps, err
	}

	//TODO: this doesn't work rn without grades
	mps, e = pages.GimmeMPs(c, IDS[userSelector-1], highSchool)
	if e != nil {
		return mps, err
	}

	return mps, nil
}

func GetGrades(w http.ResponseWriter, r *http.Request, functionName string, email string, password string, highSchool string, userSelector int) (map[string]map[string]string, error) {
	grades := map[string]map[string]string{}
	mp, err := GetMP(w, r)
	if err != nil {
		return grades, err
	}

	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)

	IDS, err := GetIDs(userSelector, c, highSchool, w)
	if err != nil {
		return grades, err
	}

	weeklySumData, err := pages.GradebookData(c, IDS[userSelector-1], mp, highSchool)
	if err != nil {
		return grades, err
	}

	for key := range weeklySumData {
		oneGrade := weeklySumData[key]
		grades[key] = oneGrade.ToDict()
	}

	return grades, nil

}

func functionForGpashandlerV1(c *gin.Context, w http.ResponseWriter, r *http.Request, email string, password string, highSchool string, userSelector int, student_grade int) (gpas map[string]float64, err error) {
	gpas = map[string]float64{}
	functionName := "Func GPAshandlerV1"
	grades, err := GetGrades(w, r, functionName, email, password, highSchool, userSelector)
	if err != nil {
		utils.APIPrintSpecificError("["+functionName+"]  GetGrades error", w, err, http.StatusInternalServerError)
		return
	}

	history, err := GetCurrentGradeHistory(w, r, functionName, email, password, highSchool, userSelector)
	if err != nil {
		utils.APIPrintSpecificError("["+functionName+"]  GetCurrentGradeHistory error", w, err, http.StatusInternalServerError)
		return
	}

	if student_grade < 9 {
		gpa, err := utils.GPAsOfMiddleSchoolers(grades)

		if err != nil {
			utils.APIPrintSpecificError("["+functionName+"]  GPAsOfMiddleSchoolers error", w, err, http.StatusInternalServerError)
			return gpas, err
		}

		gpas["weighted"] = gpa
		gpas["unweighted"] = gpa

	} else if student_grade >= 9 {
		unweighted, weighted, err := utils.GimmeCurrGPAS(grades, history)

		if err != nil {
			utils.APIPrintSpecificError("["+functionName+"]  GimmeCurrGPAS error", w, err, http.StatusInternalServerError)
			return gpas, err
		}

		gpas["weighted"] = weighted
		gpas["unweighted"] = unweighted

	}

	return gpas, nil
}

func GetGradeHistory(w http.ResponseWriter, r *http.Request, functionName string, email string, password string, highSchool string, userSelector int) (map[string]map[string]map[string]interface{}, error) {
	history := map[string]map[string]map[string]interface{}{}

	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)

	IDS, err := GetIDs(userSelector, c, highSchool, w)
	if err != nil {
		return history, err
	}

	historyData, err := pages.GradeHistoryData(c, IDS[userSelector-1], highSchool)
	if err != nil {
		return history, err
	}

	for year := range historyData {
		history[year] = map[string]map[string]interface{}{}

		for _, aYearsCourseDict := range historyData[year] {
			for courseName := range aYearsCourseDict {
				aCourseInaYear := aYearsCourseDict[courseName]
				history[year][courseName] = aCourseInaYear.ToDict()
			}
		}
	}

	//fmt.Printf("history: %v\n", history)

	return history, nil
}

func GetCurrentGradeHistory(w http.ResponseWriter, r *http.Request, functionName string, email string, password string, highSchool string, userSelector int) (map[string]string, error) {
	history := map[string]string{}

	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)
	IDS, err := GetIDs(userSelector, c, highSchool, w)
	if err != nil {
		return history, err
	}
	history, err = pages.CurrentGradeHistoryData(c, IDS[userSelector-1], highSchool)
	if err != nil {
		return history, err
	}

	return history, nil
}

func GetAssignments(w http.ResponseWriter, r *http.Request, functionName string, email string, password string, highSchool string, userSelector int) (map[string][]models.Assignment, error) {
	courseAssignments := map[string][]models.Assignment{}
	mp, err := GetMP(w, r)
	if err != nil {
		return courseAssignments, err
	}

	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)

	IDS, err := GetIDs(userSelector, c, highSchool, w)
	if err != nil {
		return courseAssignments, err
	}
	codesAndSections, err := pages.GimmeCourseCodes(c, IDS[userSelector-1], mp, highSchool)
	if err != nil {
		return courseAssignments, err
	}
	aCoursesAssignments := make([]models.Assignment, 0)
	//fmt.Println(pages.GimmeCourseCodes(c, IDS[userSelector-1], mp, highSchool))
	for courseName := range codesAndSections {
		aCoursesDict := codesAndSections[courseName]
		aCoursesAssignments, err = pages.AssignmentsDataForACourse(c, IDS[userSelector-1], mp, aCoursesDict["code"], aCoursesDict["section"], courseName, highSchool)
		if err != nil {
			return courseAssignments, err
		}
		courseAssignments[courseName] = aCoursesAssignments
	}

	return courseAssignments, nil
}

func GetSchedule(w http.ResponseWriter, r *http.Request, functionName string, email string, password string, highSchool string, userSelector int) (map[string][]models.ScheduleAssignment, error) {
	scheduleAssignments := map[string][]models.ScheduleAssignment{
		constants.ScheduleAssignmentsName: make([]models.ScheduleAssignment, 0),
	}
	mp, err := GetMP(w, r)
	if err != nil {
		return scheduleAssignments, err
	}

	c, e := utils.InitAndLogin(email, password, highSchool)
	utils.APIPrintSpecificError(functionName+": Couldn't init/login", w, e, http.StatusInternalServerError)

	IDS, err := GetIDs(userSelector, c, highSchool, w)
	if err != nil {
		return scheduleAssignments, err
	}
	codesAndSections, err := pages.GimmeCourseCodes(c, IDS[userSelector-1], mp, highSchool)
	if err != nil {
		return scheduleAssignments, err
	}
	for courseName := range codesAndSections {
		aCoursesDict := codesAndSections[courseName]
		aScheduleAssignments, err := pages.ScheduleDataForACourse(c, IDS[userSelector-1], mp, aCoursesDict["code"], aCoursesDict["section"], courseName, highSchool)
		if err != nil {
			return scheduleAssignments, err
		}
		// scheduleAssignments[courseName] = aScheduleAssignments
		scheduleAssignments[constants.ScheduleAssignmentsName] = append(scheduleAssignments[constants.ScheduleAssignmentsName], aScheduleAssignments...)
	}

	return scheduleAssignments, nil
}
