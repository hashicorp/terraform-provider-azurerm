package compute

import "testing"

func TestValidateLinuxName(t *testing.T) {
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
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// dash in the middle
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// 63 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghij",
			expected: true,
		},
		{
			// 64 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
			expected: true,
		},
		{
			// 65 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkj",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateLinuxName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidateWindowsName(t *testing.T) {
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
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// dash in the middle
			input:    "malcolm-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// 14 chars
			input:    "abcdefghijklmn",
			expected: true,
		},
		{
			// 15 chars
			input:    "abcdefghijklmno",
			expected: true,
		},
		{
			// 16 chars
			input:    "abcdefghijklmnop",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateWindowsName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestParseVirtualMachineScaleSetExtensionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VirtualMachineScaleSetExtensionResourceID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Virtual Machine Scale Set Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Expected: nil,
		},
		{
			Name:     "No Virtual Machine Scale Set Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/",
			Expected: nil,
		},
		{
			Name:     "No Extensions Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/machine1",
			Expected: nil,
		},
		{
			Name:     "No Extensions Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/machine1/extensions/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/machine1/extensions/extension1",
			Expected: &VirtualMachineScaleSetExtensionResourceID{
				Name:               "extension1",
				VirtualMachineName: "machine1",
				ResourceGroup:      "foo",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseVirtualMachineScaleSetExtensionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}

func TestParseVirtualMachineScaleSetID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VirtualMachineScaleSetResourceID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Virtual Machine Scale Set Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Expected: nil,
		},
		{
			Name:     "No Virtual Machine Scale Set Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/example",
			Expected: &VirtualMachineScaleSetResourceID{
				Name:          "example",
				ResourceGroup: "foo",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseVirtualMachineScaleSetID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}

func TestValidateDiskEncryptionSetName(t *testing.T) {
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
			//basic example
			input:    "hello",
			expected: true,
		},
		{
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't start with dot
			input:    ".hello",
			expected: false,
		},
		{
			// dot in middle
			input:    "hello.world",
			expected: true,
		},
		{
			// hyphen in middle
			input:    "hello-world",
			expected: true,
		},
		{
			// can't end with hyphen
			input:    "helloworld-",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// can't end with dot
			input:    "hello.",
			expected: false,
		},
		{
			// underscore at end
			input:    "helloworld_",
			expected: true,
		},
		{
			// 80 characters
			input:    "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			expected: true,
		},
		{
			// 81 characters
			input:    "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdef",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q...", v.input)

		_, errors := validateDiskEncryptionSetName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
