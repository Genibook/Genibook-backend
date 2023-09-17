package constants

var ConstantLinks = map[string]map[string]map[string]string{
	"Montgomery High School": {
		"base":     {"url": "https://parents.mtsd.k12.nj.us"},
		"endpoint": {"url": "https://parents.mtsd.k12.nj.us/genesis/parents?"},
		"home":     {"url": "https://parents.mtsd.k12.nj.us/genesis/sis/view?gohome=true"},
		"profile": {
			"tab1":   "studentdata",
			"tab2":   "studentsummary",
			"action": "form",
		},
		"login": {
			"url":      "https://parents.mtsd.k12.nj.us/genesis/sis/j_security_check",
			"username": "j_username",
			"password": "j_password",
		},
		"gradebook": {
			"tab1":      "studentdata",
			"tab2":      "gradebook",
			"tab3":      "weeklysummary",
			"studentid": "",
			"action":    "form",
			"mpToView":  "",
		},
		"assignments": {
			"tab1":             "studentdata",
			"tab2":             "gradebook",
			"tab3":             "listassignments",
			"studentid":        "",
			"action":           "form",
			"dateRange":        "",
			"date":             "",
			"courseAndSection": "",
			"status":           "",
		},
		"history": {
			"tab1":      "studentdata",
			"tab2":      "grading",
			"tab3":      "history",
			"action":    "form",
			"studentid": "",
		},
		"currHistory": {
			"tab1":      "studentdata",
			"tab2":      "grading",
			"tab3":      "current",
			"action":    "form",
			"studentid": "",
		},
		"transcript": {
			"tab1":               "studentdata",
			"tab2":               "grading",
			"tab3":               "transcripts",
			"action":             "showPDF",
			"studentid":          "",
			"transcriptTemplate": "DEFAULTU",
			"gpaDefinitionCode":  "",
		},
	},
}

// In Go, you can create a dictionary (also called a map) using the built-in make function. Here's an example:

// go

// // Create an empty map with string keys and int values
// m := make(map[string]int)

// // Add some key-value pairs to the map
// m["apple"] = 1
// m["banana"] = 2
// m["cherry"] = 3

// // Access a value by its key
// fmt.Println("Value of apple:", m["apple"])

// // Check if a key exists in the map
// if _, ok := m["banana"]; ok {
//     fmt.Println("Banana is in the map")
// }

// // Delete a key-value pair from the map
// delete(m, "cherry")
