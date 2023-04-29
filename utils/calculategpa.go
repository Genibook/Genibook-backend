package utils

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"webscrapper/constants"
)

func round(num float64) float64 {
	return math.Round(num*100) / 100
}

func GimmeGPAS(grades map[string]map[string]string) (unweighted float64, weighted float64, err error) {
	sumOfGradesWeighted := 0.0
	sumOfGradesUnWeighted := 0.0

	// wrong! do not divide by courses, divde by credit
	courses := float64(len(grades))

	unweighted = 0.0
	weighted = 0.0

	for key := range grades {
		grade := grades[key]
		theGrade := grade["grade"]
		class := key
		numGrade, err := strconv.ParseFloat(theGrade, 64)
		if err != nil {
			return unweighted, weighted, err
		}
		matched, err := regexp.MatchString(constants.DetermineWeightedOrNotRegrexPattern, strings.ToLower(class))
		if err != nil {
			return unweighted, weighted, err
		} else if matched {
			sumOfGradesUnWeighted += numGrade
			numGrade += 5
			sumOfGradesWeighted += numGrade

		} else {
			sumOfGradesWeighted += numGrade
			sumOfGradesUnWeighted += numGrade
		}

	}

	unweighted = round(sumOfGradesUnWeighted / courses)
	weighted = round(sumOfGradesWeighted / courses)

	return round(unweighted), round(weighted), nil
}
