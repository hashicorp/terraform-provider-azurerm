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

func TestAccAzureRMServiceFabricMeshSecretValue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_secret_value", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricMeshSecretValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricMeshSecretValue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshSecretValueExists(data.ResourceName),
				),
			},
			data.ImportStep("value"),
		},
	})
}

func TestAccAzureRMServiceFabricMeshSecretValue_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_secret_value", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricMeshSecretValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricMeshSecretValue_multiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshSecretValueExists(data.ResourceName),
				),
			},
			data.ImportStep("value"),
		},
	})
}

func testCheckAzureRMServiceFabricMeshSecretValueDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceFabricMesh.SecretValueClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_service_fabric_mesh_secret_value" {
			continue
		}

		id, err := parse.ServiceFabricMeshSecretValueID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.SecretName, id.Name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Service Fabric Mesh Secret Value still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMServiceFabricMeshSecretValueExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceFabricMesh.SecretValueClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ServiceFabricMeshSecretValueID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.SecretName, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceFabricMeshSecretValuesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Service Fabric Mesh Secret Value %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceFabricMeshSecretValue_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfm-%d"
  location = "%s"
}

resource "azurerm_service_fabric_mesh_secret" "test" {
  name                = "a"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  content_type        = "string"
}

resource "azurerm_service_fabric_mesh_secret_value" "test" {
  name                          = "accTest-%d"
  service_fabric_mesh_secret_id = azurerm_service_fabric_mesh_secret.test.id
  location                      = azurerm_resource_group.test.location
  value                         = "testValue"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceFabricMeshSecretValue_multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfm-%d"
  location = "%s"
}

resource "azurerm_service_fabric_mesh_secret" "test" {
  name                = "%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  content_type        = "string"
}

resource "azurerm_service_fabric_mesh_secret_value" "test" {
  name                          = "accTest-%d"
  service_fabric_mesh_secret_id = azurerm_service_fabric_mesh_secret.test.id
  location                      = azurerm_resource_group.test.location
  value                         = "testValue"
}

resource "azurerm_service_fabric_mesh_secret_value" "test2" {
  name                          = "accTest2-%d"
  service_fabric_mesh_secret_id = azurerm_service_fabric_mesh_secret.test.id
  location                      = azurerm_resource_group.test.location
  value                         = "testValue2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
