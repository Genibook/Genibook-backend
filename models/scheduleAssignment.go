package models

import "encoding/json"

type ScheduleAssignment struct {
	Assignment  string `json:"assignment"`
	Category    string `json:"category"`
	CourseName  string `json:"course_name"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Points      string `json:"points"`
}

func (a *ScheduleAssignment) ToJson() (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
