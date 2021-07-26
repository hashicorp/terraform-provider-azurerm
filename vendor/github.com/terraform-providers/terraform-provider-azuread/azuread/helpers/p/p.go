package p // or maybe ptr?

// helper functions to convert to a pointer
func Bool(input bool) *bool {
	return &input
}

func BoolI(i interface{}) *bool {
	b := i.(bool)
	return &b
}

func Int32(input int32) *int32 {
	return &input
}

func Int32I(i interface{}) *int32 {
	i32 := i.(int32)
	return &i32
}

func String(input string) *string {
	return &input
}

func StringI(i interface{}) *string {
	s := i.(string)
	return &s
}
