// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccApiCenterApiDefinition_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api_definition", "test")
	r := ApiCenterApiDefinitionResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"api_name":            {},
		"resource_group_name": {},
		"service_name":        {},
		"subscription_id":     {},
		"version_name":        {},
		"workspace_name":      {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_api_center_api_definition.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_center_api_definition.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_api_definition.test", tfjsonpath.New("api_name"), tfjsonpath.New("api_version_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_api_definition.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("api_version_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_api_definition.test", tfjsonpath.New("service_name"), tfjsonpath.New("api_version_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_api_definition.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("api_version_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_api_definition.test", tfjsonpath.New("version_name"), tfjsonpath.New("api_version_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_api_definition.test", tfjsonpath.New("workspace_name"), tfjsonpath.New("api_version_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
