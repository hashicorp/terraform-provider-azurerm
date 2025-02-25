package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesFleetManagerTestResource struct{}

func TestAccKubernetesFleetManager_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}
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

func TestAccKubernetesFleetManager_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}
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

func TestAccKubernetesFleetManager_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}
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

func TestAccKubernetesFleetManager_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}
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

func (r KubernetesFleetManagerTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleets.ParseFleetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ContainerService.V20231015.Fleets.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r KubernetesFleetManagerTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]s"
  location = "%[2]s"
}

`, data.RandomString, data.Locations.Primary)
}

func (r KubernetesFleetManagerTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_manager" "test" {
  name                = "acctestkfm-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomString)
}

func (r KubernetesFleetManagerTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_manager" "import" {
  name                = azurerm_kubernetes_fleet_manager.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.basic(data))
}

func (r KubernetesFleetManagerTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_manager" "test" {
  name                = "acctestkfm-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  hub_profile {
    dns_prefix = "acctestkfm-%[2]s"
  }
  tags = {
    environment = "terraform-acctests"
  }
}
`, r.template(data), data.RandomString)
}

func (r KubernetesFleetManagerTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_manager" "test" {
  name                = "acctestkfm-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  hub_profile {
    dns_prefix = "acctestkfm-%[2]s"
  }
  tags = {
    new_environment = "terraform-acctests-updated"
  }
}
`, r.template(data), data.RandomString)
}
