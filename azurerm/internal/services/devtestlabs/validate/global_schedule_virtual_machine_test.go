package validate

import "testing"

func TestGlobalScheduleVirtualMachineID(t *testing.T) {
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
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/virtualMachines/",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/VirtualMachines/machine1",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/machine1",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.ID)
		_, errors := GlobalScheduleVirtualMachineID(tc.ID, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
