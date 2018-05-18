package utils

import "strings"

func Bool(input bool) *bool {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func Int64(input int64) *int64 {
	return &input
}

func String(input string) *string {
	return &input
}

func AzureRMNormalizeCollation(input string, find string, replace string, count int) *string {
	collation := strings.Replace(input, find, replace, count)
	return &collation
}
