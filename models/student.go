package models

import (
	"encoding/json"
)

type Student struct {
	Age           int    `json:"age"`
	ImgURL        string `json:"img_url"`
	StateID       int    `json:"state_id"`
	Birthday      string `json:"birthday"`
	ScheduleLink  string `json:"schedule_link"`
	Name          string `json:"name"`
	Grade         int    `json:"grade"`
	Locker        string `json:"locker"`
	CounselorName string `json:"counselor_name"`
	ID            int    `json:"id"`
	Image64       string `json:"image64"`
}

func (s *Student) toJson() (string, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
