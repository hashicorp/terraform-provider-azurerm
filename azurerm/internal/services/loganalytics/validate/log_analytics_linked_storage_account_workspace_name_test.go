package validate

import (
	"testing"
)

func TestLogAnalyticsLinkedStorageAccountWorkspaceName(t *testing.T) {
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
			Input:    "invalid_Linked_Storage_Account_Name",
			Expected: false,
		},
		{
			Name:     "Invalid characters space",
			Input:    "invalid Linked Storage Account",
			Expected: false,
		},
		{
			Name:     "Invalid name starts with hyphen",
			Input:    "-invalidLinkedStorageAccountName",
			Expected: false,
		},
		{
			Name:     "Invalid name ends with hyphen",
			Input:    "invalidLinkedStorageAccountName-",
			Expected: false,
		},
		{
			Name:     "Invalid name too long",
			Input:    "thisIsToLooooooooooooooooooooooooongForALinkedStorageAccountName",
			Expected: false,
		},
		{
			Name:     "Valid name",
			Input:    "validLinkedStorageAccountName",
			Expected: true,
		},
		{
			Name:     "Valid name with hyphen",
			Input:    "validLinkedStorageAccountName-2",
			Expected: true,
		},
		{
			Name:     "Valid name max length",
			Input:    "thisIsTheLoooooooooooongestValidLinkedStorageAccountNameThereIs",
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

		_, errors := LogAnalyticsLinkedStorageAccountWorkspaceName(v.Input, "workspace_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
