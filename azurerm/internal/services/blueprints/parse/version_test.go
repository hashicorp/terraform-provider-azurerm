package parse

import "testing"

func TestVersionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *VersionId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "Invalid scope",
			Input: "/managementGroups/testAccManagementGroup",
			Error: true,
		},
		// We have two valid possibilities to check for
		{
			Name:  "Valid subscription scoped",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/simpleBlueprint/versions/v1-test",
			Expected: &VersionId{
				Scope:     "subscriptions/00000000-0000-0000-0000-000000000000",
				Blueprint: "simpleBlueprint",
				Name:      "v1-test",
			},
		},
		{
			Name:  "Valid management group scoped",
			Input: "/providers/Microsoft.Management/managementGroups/testAccManagementGroup/providers/Microsoft.Blueprint/blueprints/simpleBlueprint/versions/v1-test",
			Expected: &VersionId{
				Scope:     "providers/Microsoft.Management/managementGroups/testAccManagementGroup",
				Blueprint: "simpleBlueprint",
				Name:      "v1-test",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VersionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Scope", v.Expected.Name, actual.Name)
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for Scope", v.Expected.Scope, actual.Scope)
		}
	}
}
