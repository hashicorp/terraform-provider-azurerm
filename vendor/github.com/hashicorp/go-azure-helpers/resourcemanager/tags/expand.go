package tags

func Expand(input map[string]interface{}) *map[string]string {
	output := make(map[string]string)

	for k, v := range input {
		output[k] = v.(string)
	}

	return &output
}
