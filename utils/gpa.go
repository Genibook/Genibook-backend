package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	constants "webscrapper/constants/v2"
)

func round(num float64) float64 {
	return math.Round(num*100) / 100
}

func convertToFloat(val interface{}) (float32, error) {
	// Type assertion to float64
	f, ok := val.(float32)
	if !ok {
		return 0, fmt.Errorf("unable to convert %v to float32", val)
	}
	return f, nil
}

func convertToString(val interface{}) (string, error) {
	// Type assertion to float64
	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("unable to convert %v to string", val)
	}
	return s, nil
}

func convertToInt(val interface{}) (int, error) {
	// Type assertion to float64
	i, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("unable to convert %v to string", val)
	}
	return i, nil
}

func GPAsOfMiddleSchoolers(grades map[string]map[string]string) (gpa float64, err error) {

	sum := 0.0
	length := len(grades)

	gpa = 0.0
	for key := range grades {
		grade := grades[key]
		theGrade := grade["grade"]

		if theGrade == "0" {
			length--
			continue
		}

		numGrade, err := strconv.ParseFloat(theGrade, 64)
		if err != nil {
			return gpa, err
		}
		sum += numGrade
	}
	gpa = round(sum / float64(length))
	return gpa, nil
}

func GimmeHistoryGPAS(courses map[string]map[string]map[string]interface{}) (gpaHistory map[string]map[string]float64, err error) {
	gpaHistory = map[string]map[string]float64{}

	for year := range courses {
		if len(courses[year]) < 7 {
			continue
		}

		sumOfGradesWeighted := 0.0
		sumOfGradesUnWeighted := 0.0
		sumOfCredits := 0.0

		//fmt.Printf("courses[year]: %v\n", courses[year])

		for key := range courses[year] {
			grade := courses[year][key]

			theGrade := grade["FG"]
			s, e := convertToString(theGrade)
			if e != nil {
				return gpaHistory, err
			}

			numGrade, err := strconv.ParseFloat(s, 64)
			if err != nil {
				fmt.Printf("[GimmeHistoryGPAs utils/gpa.go]: WARNING - numGrade converted falsely:  %v\n", err.Error())
				continue
				// ngl we just skipping all that stuff skull

			}
			// if theGrade == "P" {
			// 	continue
			// } else if theGrade == "AP" {
			// 	continue
			// } else if theGrade == "PRO" {
			// 	continue
			// }

			class := key

			schoolGrade, err := convertToInt(grade["Grade"])
			if err != nil {
				return gpaHistory, err
			}

			matched, err := regexp.MatchString(constants.DetermineWeightedOrNotRegrexPattern, strings.ToLower(class))
			if err != nil {
				return gpaHistory, err
			}

			credd, err := convertToFloat(grade["Earned"])
			if err != nil {
				return gpaHistory, err
			}
			cred := float64(credd)
			sumOfCredits += float64(cred)

			if matched {
				if schoolGrade >= 9 {
					sumOfGradesUnWeighted += numGrade * cred
					numGrade += 5
					sumOfGradesWeighted += numGrade * cred
				} else {
					sumOfGradesWeighted += numGrade * cred
					sumOfGradesUnWeighted += numGrade * cred
				}

			} else {
				sumOfGradesWeighted += numGrade * cred
				sumOfGradesUnWeighted += numGrade * cred
			}

			// fmt.Printf("grade: %v\n", grade)
			// fmt.Printf("numGrade: %v\n", numGrade)
			// fmt.Printf("cred: %v\n", cred)

		}

		// fmt.Printf("sumOfCredits: %v\n", sumOfCredits)
		// fmt.Printf("sumOfGradesUnWeighted: %v\n", sumOfGradesUnWeighted)
		// fmt.Printf("sumOfGradesWeighted: %v\n", sumOfGradesWeighted)

		unweighted := 0.0
		weighted := 0.0

		unweighted = round(sumOfGradesUnWeighted / sumOfCredits)
		weighted = round(sumOfGradesWeighted / sumOfCredits)

		gpaHistory[year] = map[string]float64{
			"unweighted": unweighted,
			"weighted":   weighted,
		}

	}

	return gpaHistory, nil

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
		if theGrade == "0" {
			continue
		}
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
