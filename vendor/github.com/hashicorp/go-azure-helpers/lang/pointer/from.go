package pointer

func FromNilableBool(input *bool) bool {
	if input != nil {
		return *input
	}

	return false
}

func FromNilableFloat(input *float64) float64 {
	if input != nil {
		return *input
	}

	return 0.0
}

func FromNilableInt(input *int) int {
	if input != nil {
		return *input
	}

	return 0
}

func FromNilableInt64(input *int64) int64 {
	if input != nil {
		return *input
	}

	return 0
}

func FromNilableMapOfStringInterfaces(input *map[string]interface{}) map[string]interface{} {
	if input != nil {
		return *input
	}

	return map[string]interface{}{}
}

func FromNilableMapOfStringStrings(input *map[string]string) map[string]string {
	if input != nil {
		return *input
	}

	return map[string]string{}
}

func FromNilableSliceOfStrings(input *[]string) []string {
	if input != nil {
		return *input
	}

	return []string{}
}

func FromNilableString(input *string) string {
	if input != nil {
		return *input
	}

	return ""
}
