package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesFleetAutoUpgradeProfileTestResource struct{}

func TestAccKubernetesFleetAutoUpgradeProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetAutoUpgradeProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileTestResource{}

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
	r := KubernetesFleetAutoUpgradeProfileTestResource{}

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
	r := KubernetesFleetAutoUpgradeProfileTestResource{}

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

func (r KubernetesFleetAutoUpgradeProfileTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autoupgradeprofiles.ParseAutoUpgradeProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.FleetAutoUpgradeProfilesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r KubernetesFleetAutoUpgradeProfileTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                        = "acctestfaup-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Stable"
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "import" {
  name                        = azurerm_kubernetes_fleet_auto_upgrade_profile.test.name
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_auto_upgrade_profile.test.kubernetes_fleet_manager_id
  channel                     = "Stable"
}
`, r.basic(data))
}

func (r KubernetesFleetAutoUpgradeProfileTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
