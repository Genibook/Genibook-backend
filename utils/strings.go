package utils

import "strings"

func CleanAString(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
}
