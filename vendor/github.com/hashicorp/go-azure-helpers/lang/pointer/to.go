package pointer

func ToBool(input bool) *bool {
	return &input
}

func ToFloat(input float64) *float64 {
	return &input
}

func ToInt(input int) *int {
	return &input
}

func ToInt64(input int64) *int64 {
	return &input
}

func ToMapOfStringInterfaces(input map[string]interface{}) *map[string]interface{} {
	return &input
}

func ToMapOfStringStrings(input map[string]string) *map[string]string {
	return &input
}

func ToSliceOfStrings(input []string) *[]string {
	return &input
}

func ToString(input string) *string {
	return &input
}
