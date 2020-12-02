package validate

import "testing"

func TestVirtualHubID(t *testing.T) {
	testData := []struct {
		Name  string
		Input string
		Valid bool
	}{
		{
			Name:  "Empty",
			Input: "",
			Valid: false,
		},
		{
			Name:  "No Virtual Hubs Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			Name:  "No Virtual Hubs Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualHubs/",
			Valid: false,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualHubs/example",
			Valid: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		_, errors := VirtualHubID(v.Input, "virtual_hub_id")
		isValid := len(errors) == 0
		if v.Valid != isValid {
			t.Fatalf("Expected %t but got %t", v.Valid, isValid)
		}
	}
}
