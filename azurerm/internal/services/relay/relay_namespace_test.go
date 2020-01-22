package relay

import (
	"testing"
)

func TestParseRelayNamespaceID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *NamespaceResourceID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing namespaces Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Relay/namespaces/",
			Expected: nil,
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Relay/Namespaces/mynamespace",
			Expected: nil,
		},
		{
			Name:  "Relay Namespace Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Relay/namespaces/mynamespace",
			Expected: &NamespaceResourceID{
				ResourceGroup: "mygroup1",
				Name:          "mynamespace",
			},
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseNamespaceID(v.Input)
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
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}

func TestValidateNamespaceID(t *testing.T) {
	cases := []struct {
		ID    string
		Valid bool
	}{
		{
			ID:    "",
			Valid: false,
		},
		{
			ID:    "nonsense",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Relay/namespaces",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Relay/Namespaces/relay1",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Relay/namespaces/relay1",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing value %s", tc.ID)
		_, errors := ValidateNamespaceID(tc.ID, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
