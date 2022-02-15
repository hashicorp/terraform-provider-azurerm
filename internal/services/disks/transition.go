package disks

func flattenTags(input *map[string]string) map[string]interface{} {
	output := make(map[string]interface{})

	if input != nil {
		for k, v := range *input {
			val := v
			output[k] = val
		}
	}

	return output
}
