// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_maintenance_configuration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccContainerAppEnvironmentMaintenanceConfiguration_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_maintenance_configuration", "test")
	r := ContainerAppEnvironmentMaintenanceConfigurationResource{}

	checkedFields := map[string]struct{}{
		"name":                     {},
		"managed_environment_name": {},
		"resource_group_name":      {},
		"subscription_id":          {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_container_app_environment_maintenance_configuration.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_container_app_environment_maintenance_configuration.test", tfjsonpath.New("name"), knownvalue.StringExact("default")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_container_app_environment_maintenance_configuration.test", tfjsonpath.New("managed_environment_name"), tfjsonpath.New("container_app_environment_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_container_app_environment_maintenance_configuration.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("container_app_environment_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_container_app_environment_maintenance_configuration.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("container_app_environment_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
