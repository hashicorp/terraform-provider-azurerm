package parse

import (
	"testing"
)

func TestLighthouseAssignmentID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *LighthouseAssignmentId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Missing Scope",
			Input:    "providers/Microsoft.ManagedServices/registrationAssignments/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "Missing Name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationAssignments",
			Expected: nil,
		},
		{
			Name:  "Lighthouse assignment ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationAssignments/00000000-0000-0000-0000-000000000000",
			Expected: &LighthouseAssignmentId{
				Name:  "00000000-0000-0000-0000-000000000000",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := LighthouseAssignmentID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for Scope", v.Expected.Scope, actual.Scope)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
