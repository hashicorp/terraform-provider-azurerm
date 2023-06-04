package parse

import "testing"

func TestRoleAssignmentMarketplaceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *RoleAssignmentMarketplaceId
	}{
		{
			Input: "",
			Error: true,
		},

		{
			Input: "/",
			Error: true,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/",
			Error: true,
		},

		{
			Input: "/providers/Microsoft.Subscription/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Error: true,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentMarketplaceId{
				Name: "23456781-2349-8764-5631-234567890121",
			},
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|12345678-1234-5678-1234-567890123456",
			Expected: &RoleAssignmentMarketplaceId{
				Name:     "23456781-2349-8764-5631-234567890121",
				TenantId: "12345678-1234-5678-1234-567890123456",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RoleAssignmentMarketplaceID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("expected a value but got an error: %+v", err)
		}

		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Role Assignment Name", v.Expected.Name, actual.Name)
		}
	}
}
