package utils

func Bool(input bool) *bool {
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

func SliceContainsString(haystack []interface{}, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func StringSliceContainsStringSlice(leftSide, rightSide []interface{}) bool {
	for _, s := range leftSide {
		if !SliceContainsString(rightSide, s.(string)) {
			return false
		}
	}
	return true
}
