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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
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

	resp, err := clients.ContainerService.V20250301.AutoUpgradeProfiles.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r KubernetesFleetAutoUpgradeProfileTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  name                = "acctestfleet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                = "default"
  resource_group_name = azurerm_resource_group.test.name
  fleet_name          = azurerm_kubernetes_fleet_manager.test.name
  channel             = "Stable"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  name                = "acctestfleet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                    = "default"
  resource_group_name     = azurerm_resource_group.test.name
  fleet_name              = azurerm_kubernetes_fleet_manager.test.name
  channel                 = "Rapid"
  node_image_upgrade_type = "Latest"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
