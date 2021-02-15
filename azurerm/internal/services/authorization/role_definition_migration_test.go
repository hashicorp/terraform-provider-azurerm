package authorization

import (
	"testing"
)

func TestRoleDefinitionMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion    int
		InputAttributes map[string]interface{}
		ExpectedNewID   string
	}{
		"subscription_scope": {
			StateVersion: 0,
			InputAttributes: map[string]interface{}{
				"id":                 "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/11111111-1111-1111-1111-111111111111",
				"name":               "roleName",
				"role_definition_id": "11111111-1111-1111-1111-111111111111",
				"scope":              "/subscriptions/00000000-0000-0000-0000-000000000000",
			},
			ExpectedNewID: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/11111111-1111-1111-1111-111111111111|/subscriptions/00000000-0000-0000-0000-000000000000",
		},
		"managementGroup_scope": {
			StateVersion: 0,
			InputAttributes: map[string]interface{}{
				"id":                 "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/11111111-1111-1111-1111-111111111111",
				"name":               "roleName",
				"role_definition_id": "11111111-1111-1111-1111-111111111111",
				"scope":              "/providers/Microsoft.Management/managementGroups/22222222-2222-2222-2222-222222222222",
			},
			ExpectedNewID: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/11111111-1111-1111-1111-111111111111|/providers/Microsoft.Management/managementGroups/22222222-2222-2222-2222-222222222222",
		},
	}

	for _, tc := range cases {
		newID, _ := resourceArmRoleDefinitionStateUpgradeV0(tc.InputAttributes, nil)

		if newID["id"].(string) != tc.ExpectedNewID {
			t.Fatalf("ID migration failed, expected %q, got: %q", tc.ExpectedNewID, newID["id"].(string))
		}
	}
}
