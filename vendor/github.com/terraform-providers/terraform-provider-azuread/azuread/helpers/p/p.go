package p // or maybe ptr?

// helper functions to convert to a pointer
func Bool(input bool) *bool {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func String(input string) *string {
	return &input
}
