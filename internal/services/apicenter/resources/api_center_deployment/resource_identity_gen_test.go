// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccApiCenterDeployment_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_deployment", "test")
	r := ApiCenterDeploymentResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"api_name":            {},
		"resource_group_name": {},
		"service_name":        {},
		"subscription_id":     {},
		"workspace_name":      {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_api_center_deployment.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_center_deployment.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_deployment.test", tfjsonpath.New("api_name"), tfjsonpath.New("api_center_api_definition_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_deployment.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("api_center_api_definition_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_deployment.test", tfjsonpath.New("service_name"), tfjsonpath.New("api_center_api_definition_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_deployment.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("api_center_api_definition_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_deployment.test", tfjsonpath.New("workspace_name"), tfjsonpath.New("api_center_api_definition_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
