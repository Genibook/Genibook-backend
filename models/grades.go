package models

import "encoding/json"

type OneGrade struct {
	// Subject      string  `json:"subject"`
	Grade        float32 `json:"grade"`
	TeacherName  string  `json:"teacher_name"`
	TeacherEmail string  `json:"teacher_email"`
}

func (s *OneGrade) ToJson() (string, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// mathGrade := Grade{
// 	Subject:      "Math",
// 	Grade:        85.0,
// 	TeacherName:  "John Smith",
// 	TeacherEmail: "john.smith@example.com",
// }

// jsonString, err := json.Marshal(mathGrade)
// if err != nil {
// 	// handle error
// }
// fmt.Println(string(jsonString))
