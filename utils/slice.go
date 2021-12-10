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

func SliceContainsValue(input []string, value string) bool {
	for _, v := range input {
		if v == value {
			return true
		}
	}

	return false
}
