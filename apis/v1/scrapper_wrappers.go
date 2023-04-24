package api_v1

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"webscrapper/constants"
	"webscrapper/models"
	"webscrapper/pages"
	"webscrapper/utils"

	"github.com/gocolly/colly"
)

func GetMP(w http.ResponseWriter, r *http.Request) (string, error) {
	mp := r.URL.Query().Get(constants.MPFormKey)
	if !strings.Contains(mp, "MP") {
		log.Println("Marking Period Not Valid: " + mp)
		http.Error(w, "Marking Period Not Valid: "+mp, http.StatusNotAcceptable)
		return "", http.ErrBodyNotAllowed
	}
	return mp, nil
}

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
