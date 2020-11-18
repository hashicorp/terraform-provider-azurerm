package validate

import (
	"testing"
)

func TestLogAnalyticsWorkspaceName(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{
			Name:     "Too short",
			Input:    "inv",
			Expected: false,
		},
		{
			Name:     "Invalid characters underscores",
			Input:    "invalid_Log_Analytics_Workspace_Name",
			Expected: false,
		},
		{
			Name:     "Invalid characters space",
			Input:    "invalid Log Analytics Workspace Name",
			Expected: false,
		},
		{
			Name:     "Invalid name starts with hyphen",
			Input:    "-invalidLogAnalyticsWorkspaceName",
			Expected: false,
		},
		{
			Name:     "Invalid name ends with hyphen",
			Input:    "invalidLogAnalyticsWorkspaceName-",
			Expected: false,
		},
		{
			Name:     "Invalid name too long",
			Input:    "thisIsToLooooooooooooooooooooongestForALogAnalyticsWorkspaceName",
			Expected: false,
		},
		{
			Name:     "Valid name",
			Input:    "validLogAnalyticsWorkspaceName",
			Expected: true,
		},
		{
			Name:     "Valid name with hyphen",
			Input:    "validLogAnalyticsWorkspaceName-2",
			Expected: true,
		},
		{
			Name:     "Valid name max length",
			Input:    "thisIsTheLooooooooooongestValidLogAnalyticsWorkspaceNameThereIs",
			Expected: true,
		},
		{
			Name:     "Valid name min length",
			Input:    "vali",
			Expected: true,
		},
	}
	for _, v := range testCases {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		_, errors := LogAnalyticsWorkspaceName(v.Input, "workspace_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %v but got %v (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
