package policyinsights

import "testing"

func TestValidateRemediationName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "hello",
			expected: true,
		},
		{
			// start with an underscore
			input:    "_hello",
			expected: true,
		},
		{
			// end with a hyphen
			input:    "hello-",
			expected: true,
		},
		{
			// can contain an exclamation mark
			input:    "hello!",
			expected: true,
		},
		{
			// dash in the middle
			input:    "malcolm-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: true,
		},
		{
			// can't contain %
			input:    "hello%world",
			expected: false,
		},
		{
			// can't contain ^
			input:    "hello^world",
			expected: false,
		},
		{
			// can't contain #
			input:    "hello#world",
			expected: false,
		},
		{
			// can't contain ?
			input:    "hello?world",
			expected: false,
		},
		{
			// 260 chars
			input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij",
			expected: true,
		},
		{
			// 261 chars
			input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijk",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := validateRemediationName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidatePolicyAssignmentID(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// policy assignment in resource group
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Authorization/policyAssignments/assignment1",
			expected: true,
		},
		{
			// policy assignment in resource group but no name
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Authorization/policyAssignments/",
			expected: false,
		},
		{
			// the returned value of policy assignment id may not keep its casing
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.authorization/policyassignments/assignment1",
			expected: true,
		},
		{
			// policy assignment in subscription
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/assignment1",
			expected: true,
		},
		{
			// policy assignment in subscription but no name
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/",
			expected: false,
		},
		{
			// policy assignment in management group
			input:    "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/assignment1",
			expected: true,
		},
		{
			// policy assignment in management group but no name
			input:    "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/",
			expected: false,
		},
		{
			// policy assignment in resource but no name
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Authorization/policyAssignments/assignment1",
			expected: true,
		},
		{
			// policy assignment in resource but no name
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Authorization/policyAssignments/",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := validatePolicyAssignmentID(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidatePolicyDefinitionID(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// regular policy definition
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			expected: true,
		},
		{
			// regular policy definition but no name
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/",
			expected: false,
		},
		{
			// policy definition in management group
			input:    "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			expected: true,
		},
		{
			// policy definition in management group (without keep casing)
			input:    "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			expected: true,
		},
		{
			// policy definition in management group but no name
			input:    "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := validatePolicyDefinitionID(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
