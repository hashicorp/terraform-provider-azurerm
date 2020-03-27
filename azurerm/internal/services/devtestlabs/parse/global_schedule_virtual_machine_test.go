package parse

import "testing"

func TestGlobalScheduleVirtualMachine(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *GlobalScheduleVirtualMachineId
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
			Name:     "No virtual machine name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/",
			Expected: nil,
		},
		{
			Name:     "Case incorrect in path element",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/VirtualMachines/machine1",
			Expected: nil,
		},
		{
			Name:  "Valid",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/machine1",
			Expected: &GlobalScheduleVirtualMachineId{
				ResourceGroup: "myGroup1",
				Name:          "machine1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := GlobalScheduleVirtualMachineID(v.Input)
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
