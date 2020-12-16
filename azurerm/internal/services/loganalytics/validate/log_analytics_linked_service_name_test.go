package validate

import (
	"testing"
)

func TestLogAnalyticsLinkedServiceName(t *testing.T) {
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
			Input:    "invalid_Linked_Service_Name",
			Expected: false,
		},
		{
			Name:     "Invalid characters space",
			Input:    "invalid Linked Service Name",
			Expected: false,
		},
		{
			Name:     "Invalid name starts with hyphen",
			Input:    "-invalidLinkedServiceName",
			Expected: false,
		},
		{
			Name:     "Invalid name ends with hyphen",
			Input:    "invalidLinkedServicesName-",
			Expected: false,
		},
		{
			Name:     "Invalid name too long",
			Input:    "thisIsToLoooooooooooooooooooooooooooooooongForALinkedServiceName",
			Expected: false,
		},
		{
			Name:     "Valid name",
			Input:    "validLinkedServiceName",
			Expected: true,
		},
		{
			Name:     "Valid name with hyphen",
			Input:    "validLinkedServiceName-2",
			Expected: true,
		},
		{
			Name:     "Valid name max length",
			Input:    "thisIsTheLooooooooooooooooooongestValidLinkedServiceNameThereIs",
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

		_, errors := LogAnalyticsLinkedServiceName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
