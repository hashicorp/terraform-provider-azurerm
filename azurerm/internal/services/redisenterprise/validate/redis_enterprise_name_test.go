package validate

import (
	"testing"
)

func TestRedisEnterpriseName(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{
			Name:     "Invalid empty string",
			Input:    "",
			Expected: false,
		},
		{
			Name:     "Invalid To short",
			Input:    "in",
			Expected: false,
		},
		{
			Name:     "Invalid characters underscores",
			Input:    "invalid_Exports_Name",
			Expected: false,
		},
		{
			Name:     "Invalid characters space",
			Input:    "invalid Redis Enterprise name",
			Expected: false,
		},
		{
			Name:     "Invalid name starts with hyphen",
			Input:    "-invalidRedisEnterpriseName",
			Expected: false,
		},
		{
			Name:     "Invalid name ends with hyphen",
			Input:    "invalidRedisEnterpriseName-",
			Expected: false,
		},
		{
			Name:     "Invalid name with consecutive hyphens",
			Input:    "validStorageInsightConfigName--2",
			Expected: false,
		},
		{
			Name:     "Invalid name over max length",
			Input:    "thisIsToLoooooooooooooooooooooooooooooongForARedisEnterpriseName",
			Expected: false,
		},
		{
			Name:     "Valid name max length",
			Input:    "thisIsTheLooooooooooooooooooooooooooooongestRedisEnterpriseName",
			Expected: true,
		},
		{
			Name:     "Valid name",
			Input:    "validRedisEnterpriseName",
			Expected: true,
		},
		{
			Name:     "Valid name with hyphen",
			Input:    "validStorageInsightConfigName-2",
			Expected: true,
		},
		{
			Name:     "Valid name min length",
			Input:    "val",
			Expected: true,
		},
	}
	for _, v := range testCases {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		_, errors := RedisEnterpriseName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %v but got %v (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
