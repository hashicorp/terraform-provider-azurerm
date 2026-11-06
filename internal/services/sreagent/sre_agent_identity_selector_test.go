// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestSreAgentIdentitySelectorsRequireUserAssignedIdentity(t *testing.T) {
	cases := []struct {
		name  string
		value string
		valid bool
	}{
		{"user-assigned identity", testIdentityID, true},
		{"separate user-assigned identity", testIdentityID + "-resources", true},
		{"resource group", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example", false},
		{"virtual network", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/example", false},
		{"agent resource", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.App/agents/agent", false},
		{"principal ID", "11111111-1111-1111-1111-111111111111", false},
		{"system sentinel outside setter scope", "system", false},
		{"empty", "", false},
	}
	arguments := SreAgentResource{}.Arguments()
	for _, block := range []string{"action_configuration", "resources_configuration"} {
		t.Run(block, func(t *testing.T) {
			validate := arguments[block].Elem.(*pluginsdk.Resource).Schema["identity_id"].ValidateFunc
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					_, errs := validate(tc.value, block+".0.identity_id")
					if (len(errs) == 0) != tc.valid {
						t.Fatalf("identity selector %q: valid=%t errors=%v", tc.value, tc.valid, errs)
					}
				})
			}
		})
	}
}
