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

type ServiceFabricMeshApplicationResource struct{}

func TestAccServiceFabricMeshApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_application", "test")
	r := ServiceFabricMeshApplicationResource{}

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

func TestAccServiceFabricMeshApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_application", "test")
	r := ServiceFabricMeshApplicationResource{}

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

func (r ServiceFabricMeshApplicationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ApplicationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ServiceFabricMesh.ApplicationClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Service Fabric Mesh Application %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r ServiceFabricMeshApplicationResource) basic(data acceptance.TestData) string {
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

func (r ServiceFabricMeshApplicationResource) update(data acceptance.TestData) string {
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
