// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccKubernetesFleetAutoUpgradeProfile_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"fleet_name":          {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_kubernetes_fleet_auto_upgrade_profile.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_kubernetes_fleet_auto_upgrade_profile.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_kubernetes_fleet_auto_upgrade_profile.test", tfjsonpath.New("fleet_name"), tfjsonpath.New("kubernetes_fleet_manager_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_kubernetes_fleet_auto_upgrade_profile.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("kubernetes_fleet_manager_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_kubernetes_fleet_auto_upgrade_profile.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("kubernetes_fleet_manager_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
