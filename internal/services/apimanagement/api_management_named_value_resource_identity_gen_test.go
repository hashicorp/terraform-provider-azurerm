// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccApiManagementNamedValue_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")
	r := ApiManagementNamedValueResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"named_value_id":      {},
		"resource_group_name": {},
		"service_name":        {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_api_management_named_value.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_api_management_named_value.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_management_named_value.test", tfjsonpath.New("named_value_id"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_management_named_value.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_management_named_value.test", tfjsonpath.New("service_name"), tfjsonpath.New("api_management_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
