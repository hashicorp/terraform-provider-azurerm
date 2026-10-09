// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesFleetAutoUpgradeProfileResource struct{}

func TestAccKubernetesFleetAutoUpgradeProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_image_selection_type").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetAutoUpgradeProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccKubernetesFleetAutoUpgradeProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetAutoUpgradeProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetAutoUpgradeProfile_selectionRemoval(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.selection(data, "Latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_image_selection_type").HasValue("Latest"),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.selection(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_image_selection_type").HasValue(""),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.selection(data, "Latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_image_selection_type").HasValue("Latest"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetAutoUpgradeProfile_reenable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.selection(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("node_image_selection_type").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetAutoUpgradeProfile_updateStrategyRemoval(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.strategy(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("update_strategy_id").MatchesOtherKey(check.That("azurerm_kubernetes_fleet_update_strategy.test").Key("id")),
			),
		},
		data.ImportStep(),
		{
			Config: r.strategy(data, false),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionReplace),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("update_strategy_id").HasValue(""),
			),
		},
		data.ImportStep(),
		{
			Config: r.strategy(data, true),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionReplace),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("update_strategy_id").MatchesOtherKey(check.That("azurerm_kubernetes_fleet_update_strategy.test").Key("id")),
			),
		},
		data.ImportStep(),
	})
}

func (r KubernetesFleetAutoUpgradeProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autoupgradeprofiles.ParseAutoUpgradeProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.FleetAutoUpgradeProfilesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r KubernetesFleetAutoUpgradeProfileResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                        = "acctestfaup-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Stable"
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "import" {
  name                        = azurerm_kubernetes_fleet_auto_upgrade_profile.test.name
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_auto_upgrade_profile.test.kubernetes_fleet_manager_id
  channel                     = "Stable"
}
`, r.basic(data))
}

func (r KubernetesFleetAutoUpgradeProfileResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_update_strategy" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  stage {
    name = "acctestfus-%[2]d"
    group {
      name = "acctestfus-%[2]d"
    }
  }
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                        = "acctestfaup-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Rapid"
  node_image_selection_type   = "Latest"
  update_strategy_id          = azurerm_kubernetes_fleet_update_strategy.test.id
  enabled                     = true
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_update_strategy" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  stage {
    name = "acctestfus-%[2]d"
    group {
      name = "acctestfus-%[2]d"
    }
  }
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                        = "acctestfaup-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Stable"
  node_image_selection_type   = "Consistent"
  update_strategy_id          = azurerm_kubernetes_fleet_update_strategy.test.id
  enabled                     = false
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestkfm-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileResource) selection(data acceptance.TestData, selection string) string {
	selectionConfig := ""
	if selection != "" {
		selectionConfig = fmt.Sprintf("node_image_selection_type = %q", selection)
	}
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                        = "acctestfaup-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Stable"
  enabled                     = false
  %[3]s
}
`, r.template(data), data.RandomInteger, selectionConfig)
}

func (r KubernetesFleetAutoUpgradeProfileResource) strategy(data acceptance.TestData, useStrategy bool) string {
	strategyConfig := ""
	if useStrategy {
		strategyConfig = "update_strategy_id = azurerm_kubernetes_fleet_update_strategy.test.id"
	}
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_update_strategy" "test" {
  name                        = "acctestfus-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  stage {
    name = "test"
    group {
      name = "test"
    }
  }
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                        = "acctestfaup-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Stable"
  enabled                     = false
  %[3]s
}
`, r.template(data), data.RandomInteger, strategyConfig)
}
