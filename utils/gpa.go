package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"webscrapper/constants"
)

func round(num float64) float64 {
	return math.Round(num*100) / 100
}

func convertToFloat(val interface{}) (float64, error) {
	// Type assertion to float64
	f, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("unable to convert %v to float64", val)
	}
	return f, nil
}

// PRECONDITION: make sure that NO DUPLICATES
// returns last occuring course using CourseName string
func locateHistoricalCourse(courses map[string]map[string]map[string]interface{}, courseName string) (map[string]interface{}, error) {

	course := map[string]interface{}{}
	found := false
	for year := range courses {
		aYearsCourses := courses[year]
		//fmt.Println("year: " + year)
		for coursee := range aYearsCourses {
			//fmt.Printf("coursee: %v\n", coursee)
			courseee := aYearsCourses[coursee]
			if courseee["Name"] == courseName {
				found = true
				course = courseee
			}
		}

	}
	if !found {
		return course, fmt.Errorf("[locateCourse]: cannot locate a course with the name " + courseName)
	}
	return course, nil
}

func GPAsOfMiddleSchoolers(grades map[string]map[string]string) (gpa float64, err error) {

	sum := 0.0
	length := len(grades)

	gpa = 0.0
	for key := range grades {
		grade := grades[key]
		theGrade := grade["grade"]
		numGrade, err := strconv.ParseFloat(theGrade, 64)
		if err != nil {
			return gpa, err
		}
		sum += numGrade
	}
	gpa = round(sum / float64(length))
	return gpa, nil
}

func GimmeGPAS(grades map[string]map[string]string, courses map[string]map[string]map[string]interface{}) (unweighted float64, weighted float64, err error) {
	sumOfGradesWeighted := 0.0
	sumOfGradesUnWeighted := 0.0
	sumOfCredits := 0.0

	unweighted = 0.0
	weighted = 0.0

	for key := range grades {
		grade := grades[key]
		theGrade := grade["grade"]
		class := key
		aClassDict, err := locateHistoricalCourse(courses, class)
		if err != nil {
			return unweighted, weighted, err
		}
		numGrade, err := strconv.ParseFloat(theGrade, 64)
		if err != nil {
			return unweighted, weighted, err
		}
		matched, err := regexp.MatchString(constants.DetermineWeightedOrNotRegrexPattern, strings.ToLower(class))
		if err != nil {
			return unweighted, weighted, err
		}

		cred, err := convertToFloat(aClassDict["Earned"])
		if err != nil {
			return unweighted, weighted, err
		}
		sumOfCredits += cred

		if matched {
			sumOfGradesUnWeighted += numGrade
			numGrade += 5
			sumOfGradesWeighted += numGrade

		} else {
			sumOfGradesWeighted += numGrade
			sumOfGradesUnWeighted += numGrade
		}

	}

	unweighted = round(sumOfGradesUnWeighted / sumOfCredits)
	weighted = round(sumOfGradesWeighted / sumOfCredits)

	return round(unweighted), round(weighted), nil
}

func GimmeCurrGPAS(grades map[string]map[string]string, courses map[string]string) (unweighted float64, weighted float64, err error) {
	sumOfGradesWeighted := 0.0
	sumOfGradesUnWeighted := 0.0
	sumOfCredits := 0.0

	unweighted = 0.0
	weighted = 0.0

	for key := range grades {
		grade := grades[key]
		theGrade := grade["grade"]
		class := key
		numGrade, err := strconv.ParseFloat(theGrade, 64)
		numGrade = math.Round(numGrade)
		if err != nil {
			return unweighted, weighted, err
		}
		matched, err := regexp.MatchString(constants.DetermineWeightedOrNotRegrexPattern, strings.ToLower(class))
		if err != nil {
			return unweighted, weighted, err
		}

		cred, err := strconv.ParseFloat(courses[class], 64)
		if err != nil {
			return unweighted, weighted, err
		}
		sumOfCredits += cred

		if matched {
			sumOfGradesUnWeighted += numGrade * cred
			numGrade += 5
			sumOfGradesWeighted += numGrade * cred

		} else {
			sumOfGradesWeighted += numGrade * cred
			sumOfGradesUnWeighted += numGrade * cred
		}

		// fmt.Printf("class: %v\n", class)
		// fmt.Printf("%f * %f\n", numGrade, cred)
		// fmt.Printf("courses: %v\n", courses)
		// fmt.Printf("sumOfCredits: %v\n", sumOfCredits)
		// fmt.Printf("sumOfGradesUnWeighted: %v\n", sumOfGradesUnWeighted)
		// fmt.Printf("sumOfGradesWeighted: %v\n", sumOfGradesWeighted)

	}

	unweighted = round(sumOfGradesUnWeighted / sumOfCredits)
	weighted = round(sumOfGradesWeighted / sumOfCredits)

	return round(unweighted), round(weighted), nil
}
