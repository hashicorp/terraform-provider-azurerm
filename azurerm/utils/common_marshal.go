package utils

// todo array -> slice ?
func ExpandStringArray(input []interface{}) *[]string {
	result := make([]string, 0)
	for _, item := range input {
		result = append(result, item.(string))
	}
	return &result
}

// todo array -> slice ?
func FlattenStringArray(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}
