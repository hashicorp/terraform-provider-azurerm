package tags

func FromTypedObject(input map[string]string) map[string]*string {
	output := make(map[string]*string, len(input))

	for k, v := range input {
		// Validate should have ignored this error already
		value, _ := TagValueToString(v)
		output[k] = &value
	}

	return output
}

func ToTypedObject(input map[string]*string) map[string]string {
	output := make(map[string]string)

	for k, v := range input {
		if v == nil {
			continue
		}

		output[k] = *v
	}

	return output
}
