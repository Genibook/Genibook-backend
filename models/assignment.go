package models

import (
	"encoding/json"
)

type Assignment struct {
	CourseName   string `json:"course_name"`
	MP           string `json:"mp"`
	DayName      string `json:"dayname"`
	FullDayName  string `json:"full_dayname"`
	Date         string `json:"date"`
	FullDate     string `json:"full_date"`
	Teacher      string `json:"teacher"`
	Category     string `json:"category"`
	Assignment   string `json:"assignment"`
	Description  string `json:"description"`
	GradePercent string `json:"grade_percent"`
	GradeNum     string `json:"grade_num"`
	Comment      string `json:"comment"`
	Prev         string `json:"prev"`
	Docs         string `json:"docs"`
}

func (a *Assignment) toJson() (string, error) {
	jsonData, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
