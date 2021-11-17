package shim

func mapStringPtrToMapString(input map[string]*string) map[string]string {
	output := make(map[string]string, len(input))

	for k, v := range input {
		if v == nil {
			continue
		}

		output[k] = *v
	}

	return output
}

func mapStringToMapStringPtr(input map[string]string) map[string]*string {
	output := make(map[string]*string, len(input))

	for k, v := range input {
		v := v
		output[k] = &v
	}

	return output
}
