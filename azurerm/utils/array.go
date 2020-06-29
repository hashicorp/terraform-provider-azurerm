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
