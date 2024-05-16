package storage

func expandStringSlice(input []interface{}) []string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		}
	}
	return result
}
