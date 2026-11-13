// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccApiCenterEnvironment_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_environment", "test")
	r := ApiCenterEnvironmentResource{}

	checkedFields := map[string]struct{}{
		"workspace_name":      {},
		"name":                {},
		"resource_group_name": {},
		"service_name":        {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_api_center_environment.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_api_center_environment.test", tfjsonpath.New("workspace_name"), knownvalue.StringExact("default")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_center_environment.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_environment.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("service_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_environment.test", tfjsonpath.New("service_name"), tfjsonpath.New("service_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_environment.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("service_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
