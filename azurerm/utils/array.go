package utils

// RemoveFromStringArray removes all matching values from a string array
func RemoveFromStringArray(elements []string, remove string) []string {
	for i, v := range elements {
		if v == remove {
			return append(elements[:i], elements[i+1:]...)
		}
	}
	return elements
}

// ContainsInStringArray checks a string array for a matching string and returns true if so
func ContainsInStringArray(elements []string, contains string) bool {
	for _, v := range elements {
		if v == contains {
			return true
		}
	}
	return false
}
