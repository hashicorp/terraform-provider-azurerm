package validate

import (
	"testing"
)

func TestLogAnalyticsLinkedStorageAccountWorkspaceName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "your case",
			Expected: true,
		},
		{
			Input:    "your case",
			Expected: false,
		},
	}
	for _, v := range testCases {
		_, errors := LogAnalyticsLinkedStorageAccountWorkspaceName(v.Input, "workspace_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
