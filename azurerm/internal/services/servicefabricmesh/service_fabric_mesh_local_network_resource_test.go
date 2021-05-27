package servicefabricmesh_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabricmesh/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceFabricMeshLocalNetworkResource struct{}

func TestAccServiceFabricMeshLocalNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_local_network", "test")
	r := ServiceFabricMeshLocalNetworkResource{}

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

func TestAccServiceFabricMeshLocalNetwork_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_local_network", "test")
	r := ServiceFabricMeshLocalNetworkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ServiceFabricMeshLocalNetworkResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ServiceFabricMesh.NetworkClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Service Fabric Mesh Local Network %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r ServiceFabricMeshLocalNetworkResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfm-%d"
  location = "%s"
}

resource "azurerm_service_fabric_mesh_local_network" "test" {
  name                   = "accTest-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  network_address_prefix = "10.0.0.0/22"

  description = "Test Description"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricMeshLocalNetworkResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfm-%d"
  location = "%s"
}

resource "azurerm_service_fabric_mesh_local_network" "test" {
  name                   = "accTest-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  network_address_prefix = "10.1.0.0/22"
  description            = "Test Description"

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
