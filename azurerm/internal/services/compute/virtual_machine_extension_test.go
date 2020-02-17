package compute

import "testing"

func TestParseVirtualMachineExtensionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VirtualMachineExtensionID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No virtual machine segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/",
			Expected: nil,
		},
		{
			Name:     "No extension name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/machine1/extension/",
			Expected: nil,
		},
		{
			Name:     "Case incorrect in path element",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/machine1/Extensions/extName",
			Expected: nil,
		},
		{
			Name:  "Valid",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.Compute/virtualMachines/machine1/extensions/extName",
			Expected: &VirtualMachineExtensionID{
				ResourceGroup:  "myGroup1",
				Name:           "extName",
				VirtualMachine: "machine1",
			},
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseVirtualMachineExtensionID(v.Input)
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

func TestValidateVirtualMachineExtensionID(t *testing.T) {
	cases := []struct {
		ID    string
		Valid bool
	}{
		{
			ID:    "",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/machine1/extension/",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/machine1/Extensions/extName",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.Compute/virtualMachines/machine1/extensions/extName",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.ID)
		_, errors := ValidateVirtualMachineExtensionID(tc.ID, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
