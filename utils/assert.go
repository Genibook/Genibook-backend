package utils

func Assert(pred bool, location string) {
	if !pred {
		panic("Assertion failed: " + location)
	}
}
