package utils

func ExpandStringSlice(input []interface{}) *[]string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return &result
}

func ExpandMapStringPtrString(input map[string]interface{}) map[string]*string {
	result := make(map[string]*string)
	for k, v := range input {
		result[k] = String(v.(string))
	}
	return result
}

func FlattenStringSlice(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func FlattenMapStringPtrString(input map[string]*string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range input {
		if v == nil {
			result[k] = ""
		} else {
			result[k] = *v
		}
	}
	return result
}
