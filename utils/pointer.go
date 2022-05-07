package utils

func Bool(input bool) *bool {
	return &input
}

func Int(input int) *int {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func Int64(input int64) *int64 {
	return &input
}

func Float(input float64) *float64 {
	return &input
}

func String(input string) *string {
	return &input
}

// generic function to get the value of a pointer
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type float interface {
	~float32 | ~float64
}

type Number interface {
	Integer | UInteger | float
}

type BasicType interface {
	Number | ~bool | ~string
}

// Val2Ptr generic function to get the value of a pointer
func Val2Ptr[T BasicType](input T) *T {
	return &input
}

// Ptr2ValOrNil generic function to get the value of a pointer, for default value return nil
func RealVal2Ptr[T BasicType](input T) (t *T) {
	if input == *new(T) {
		return nil
	}
	return Val2Ptr(input)
}

// Ptr2Val generic function to get the value of a pointer
func Ptr2Val[T BasicType](input *T) (t T) {
	if input != nil {
		return *input
	}
	return
}
