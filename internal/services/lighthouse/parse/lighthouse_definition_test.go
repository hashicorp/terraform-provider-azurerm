package parse

import (
	"testing"
)

func TestLighthouseDefinitionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *LighthouseDefinitionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Missing Scope",
			Input:    "providers/Microsoft.ManagedServices/registrationDefinitions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "Missing LighthouseDefinitionID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationDefinitions",
			Expected: nil,
		},
		{
			Name:  "Lighthouse Definition ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationDefinitions/00000000-0000-0000-0000-000000000000",
			Expected: &LighthouseDefinitionId{
				LighthouseDefinitionID: "00000000-0000-0000-0000-000000000000",
				Scope:                  "/subscriptions/00000000-0000-0000-0000-000000000000",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := LighthouseDefinitionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for Scope", v.Expected.Scope, actual.Scope)
		}

		if actual.LighthouseDefinitionID != v.Expected.LighthouseDefinitionID {
			t.Fatalf("Expected %q but got %q for LighthouseDefinitionID",
				v.Expected.LighthouseDefinitionID, actual.LighthouseDefinitionID)
		}
	}
}
