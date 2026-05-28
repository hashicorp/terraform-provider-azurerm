// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccCognitiveAccountConnectionCustomKeys_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"resource_group_name": {},
		"account_name":        {},
		"name":                {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_cognitive_account_connection_custom_keys.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_cognitive_account_connection_custom_keys.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cognitive_account_connection_custom_keys.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("cognitive_account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cognitive_account_connection_custom_keys.test", tfjsonpath.New("account_name"), tfjsonpath.New("cognitive_account_id")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_cognitive_account_connection_custom_keys.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(true),
		data.ImportBlockWithIDStep(true),
	}, false)
}
