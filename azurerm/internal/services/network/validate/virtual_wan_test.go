package validate

import (
	"testing"
)

func TestValidateVirtualWan(t *testing.T) {
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
			Name:  "No Virtual Wans Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			Name:  "No Virtual Wans Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualWans/",
			Valid: false,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualWans/example",
			Valid: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		_, errors := ValidateVirtualWanID(v.Input, "virtual_wan_id")
		isValid := len(errors) == 0
		if v.Valid != isValid {
			t.Fatalf("Expected %t but got %t", v.Valid, isValid)
		}
	}
}
