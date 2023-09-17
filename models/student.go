package models

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
	LunchBalance  string `json:"lunchbalance"`
	Image64       string `json:"image64"`
}

func (s *Student) ToDict() map[string]interface{} {
	ret := map[string]interface{}{}
	ret["age"] = s.Age
	ret["img_url"] = s.ImgURL
	ret["state_id"] = s.StateID
	ret["birthday"] = s.Birthday
	ret["schedule_link"] = s.ScheduleLink
	ret["name"] = s.Name
	ret["grade"] = s.Grade
	ret["locker"] = s.Locker
	ret["counselor_name"] = s.CounselorName
	ret["id"] = s.ID
	ret["lunchbalance"] = s.LunchBalance
	ret["image64"] = s.Image64

	return ret
}
