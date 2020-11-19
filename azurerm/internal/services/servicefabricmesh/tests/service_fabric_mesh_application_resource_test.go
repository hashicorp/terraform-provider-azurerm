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

func TestAccAzureRMServiceFabricMeshApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricMeshApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricMeshApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshApplicationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricMeshApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricMeshApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricMeshApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshApplicationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceFabricMeshApplication_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshApplicationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceFabricMeshApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricMeshApplicationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMServiceFabricMeshApplicationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceFabricMesh.ApplicationClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_service_fabric_mesh_application" {
			continue
		}

		id, err := parse.ServiceFabricMeshApplicationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Service Fabric Mesh Application still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMServiceFabricMeshApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceFabricMesh.ApplicationClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ServiceFabricMeshApplicationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceFabricMeshApplicationsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Service Fabric Mesh Application %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceFabricMeshApplication_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfm-%d"
  location = "%s"
}

resource "azurerm_service_fabric_mesh_application" "test" {
  name                = "accTest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  service {
    name    = "testservice1"
    os_type = "Linux"

    code_package {
      name       = "testcodepackage1"
      image_name = "seabreeze/sbz-helloworld:1.0-alpine"

      resources {
        requests {
          memory = 1
          cpu    = 1
        }
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceFabricMeshApplication_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfm-%d"
  location = "%s"
}

resource "azurerm_service_fabric_mesh_application" "test" {
  name                = "accTest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  service {
    name    = "testservice1"
    os_type = "Linux"

    code_package {
      name       = "testcodepackage1"
      image_name = "seabreeze/sbz-helloworld:1.0-alpine"

      resources {
        requests {
          memory = 1
          cpu    = 1
        }
      }
    }
  }

  service {
    name    = "testservice2"
    os_type = "Linux"

    code_package {
      name       = "testcodepackage2"
      image_name = "seabreeze/sbz-helloworld:1.0-alpine"

      resources {
        requests {
          memory = 2
          cpu    = 2
        }
      }
    }
  }

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
