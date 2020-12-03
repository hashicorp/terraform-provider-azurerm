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
			Input:    "invalid_Workspace_Name",
			Expected: false,
		},
		{
			Name:     "Invalid characters space",
			Input:    "invalid Workspace Name",
			Expected: false,
		},
		{
			Name:     "Invalid name starts with hyphen",
			Input:    "-invalidWorkspacetName",
			Expected: false,
		},
		{
			Name:     "Invalid name ends with hyphen",
			Input:    "invalidWorkspaceName-",
			Expected: false,
		},
		{
			Name:     "Invalid name too long",
			Input:    "thisIsToLoooooooooooooooooooooooooooooooooooongForAWorkspaceName",
			Expected: false,
		},
		{
			Name:     "Valid name",
			Input:    "validWorkspaceName",
			Expected: true,
		},
		{
			Name:     "Valid name with hyphen",
			Input:    "validWorkspaceName-2",
			Expected: true,
		},
		{
			Name:     "Valid name max length",
			Input:    "thisIsTheLooooooooooooooooooooooongestValidWorkspaceNameThereIs",
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
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
