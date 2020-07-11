package utils

func Coalesce(val bool, first string, second string) string {
	if val {
		return first
	}

	return second
}
