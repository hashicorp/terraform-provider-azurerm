package parse

import "testing"

func TestManagementGroupID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *ManagementGroupId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Management Groups Segment",
			Input: "/providers/Microsoft.Management",
			Error: true,
		},
		{
			Name:  "No Management Group ID",
			Input: "/providers/Microsoft.Management/managementGroups/",
			Error: true,
		},
		{
			Name:  "Management Group ID in UUID",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				GroupID: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management Group ID in Readable ID",
			Input: "/providers/Microsoft.Management/managementGroups/myGroup",
			Expected: &ManagementGroupId{
				GroupID: "myGroup",
			},
		},
		{
			Name:  "Invalid Management group id",
			Input: "/providers/Microsoft.Management/managementGroups/myGroup/another",
			Error: true,
		},
		{
			Name:  "Management Group ID in UUID with wrong casing",
			Input: "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				GroupID: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management Group ID in Readable ID with wrong casing",
			Input: "/providers/microsoft.management/managementgroups/group1",
			Expected: &ManagementGroupId{
				GroupID: "group1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ManagementGroupID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.GroupID != v.Expected.GroupID {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.GroupID, actual.GroupID)
		}
	}
}

func TestValidateManagementGroupName(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: false,
		},
		{
			Name:     "Proper management group name",
			Input:    "Hello",
			Expected: true,
		},
		{
			Name:     "Braces allowed",
			Input:    "Hello(world)",
			Expected: true,
		},
		{
			Name:     "Period allowed",
			Input:    "Hello.world",
			Expected: true,
		},
		{
			Name:     "Hyphen allowed",
			Input:    "Hello-world",
			Expected: true,
		},
		{
			Name:     "Underscore allowed",
			Input:    "Hello_world",
			Expected: true,
		},
		{
			Name:     "Asterisk not allowed",
			Input:    "hello*world",
			Expected: false,
		},
		{
			Name:     "Comma not allowed",
			Input:    "Hello,world",
			Expected: false,
		},
		{
			Name:     "Space not allowed",
			Input:    "Hello world",
			Expected: false,
		},
		{
			Name:     "90 characters",
			Input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij",
			Expected: true,
		},
		{
			Name:     "91 characters",
			Input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijk",
			Expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q...", v.Name)

		_, errors := ValidateManagementGroupName(v.Input, "name")
		actual := len(errors) == 0
		if v.Expected != actual {
			t.Fatalf("Expected %t but got %t", v.Expected, actual)
		}
	}
}
