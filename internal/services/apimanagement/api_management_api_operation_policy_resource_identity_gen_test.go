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

func TestAccApiManagementApiOperationPolicy_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")
	r := ApiManagementApiOperationPolicyResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"api_id":              {},
		"operation_id":        {},
		"resource_group_name": {},
		"service_name":        {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.rawXml(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_api_management_api_operation_policy.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_api_management_api_operation_policy.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_management_api_operation_policy.test", tfjsonpath.New("api_id"), tfjsonpath.New("api_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_management_api_operation_policy.test", tfjsonpath.New("operation_id"), tfjsonpath.New("operation_id")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_management_api_operation_policy.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_management_api_operation_policy.test", tfjsonpath.New("service_name"), tfjsonpath.New("api_management_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
