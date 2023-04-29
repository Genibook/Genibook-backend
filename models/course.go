package models

type Course struct {
	SchoolYear string  `json:"school_year"`
	Grade      int     `json:"grade"`
	Name       string  `json:"description"`
	School     string  `json:"school"`
	FG         string  `json:"fg"`
	Attempted  float32 `json:"attempted"`
	Earned     float32 `json:"earned"`
}

func (c *Course) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"SchoolYear": c.SchoolYear,
		"Grade":      c.Grade,
		"Name":       c.Name,
		"School":     c.School,
		"FG":         c.FG,
		"Attempted":  c.Attempted,
		"Earned":     c.Earned,
	}
}
