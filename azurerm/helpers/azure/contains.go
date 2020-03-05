package azure

func SliceContainsValue(input []string, value string) bool {
	for _, v := range input {
		if v == value {
			return true
		}
	}

	return false
}
