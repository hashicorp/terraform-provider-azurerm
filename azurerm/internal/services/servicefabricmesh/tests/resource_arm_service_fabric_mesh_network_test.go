package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabricmesh/parse"
)

func TestAccAzureRMServiceFabricMeshNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricMeshNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricMeshNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricMeshNetwork_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricMeshNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricMeshNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceFabricMeshNetwork_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceFabricMeshNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshNetworkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMServiceFabricMeshNetworkDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceFabricMesh.NetworkClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_service_fabric_mesh_network" {
			continue
		}

		id, err := parse.ServiceFabricMeshNetworkID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Service Fabric Mesh Network still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMServiceFabricMeshNetworkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceFabricMesh.NetworkClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ServiceFabricMeshNetworkID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceFabricMeshNetworksClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Service Fabric Mesh Network %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceFabricMeshNetwork_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfm-%d"
  location = "%s"
}

resource "azurerm_service_fabric_mesh_network" "test" {
  name                = "accTest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  description = "Test Description"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceFabricMeshNetwork_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfm-%d"
  location = "%s"
}

resource "azurerm_service_fabric_mesh_network" "test" {
  name                = "accTest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  
  description = "Test Description"
  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
