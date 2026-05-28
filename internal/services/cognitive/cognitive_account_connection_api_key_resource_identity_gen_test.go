// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccCognitiveAccountConnectionApiKey_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"account_name":        {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_cognitive_account_connection_api_key.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_cognitive_account_connection_api_key.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cognitive_account_connection_api_key.test", tfjsonpath.New("account_name"), tfjsonpath.New("cognitive_account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cognitive_account_connection_api_key.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("cognitive_account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cognitive_account_connection_api_key.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("cognitive_account_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(true),
		data.ImportBlockWithIDStep(true),
	}, false)
}
