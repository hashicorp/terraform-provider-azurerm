package servicefabricmesh_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabricmesh/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceFabricMeshSecretValueResource struct{}

func TestAccServiceFabricMeshSecretValue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_secret_value", "test")
	r := ServiceFabricMeshSecretValueResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("value"),
	})
}

func TestAccServiceFabricMeshSecretValue_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_mesh_secret_value", "test")
	r := ServiceFabricMeshSecretValueResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("value"),
	})
}

func (r ServiceFabricMeshSecretValueResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SecretValueID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ServiceFabricMesh.SecretValueClient.Get(ctx, id.ResourceGroup, id.SecretName, id.ValueName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Service Fabric Mesh Secret Value %q (Secret %q / Resource Group %q): %+v", id.ValueName, id.SecretName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r ServiceFabricMeshSecretValueResource) basic(data acceptance.TestData) string {
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

func (r ServiceFabricMeshSecretValueResource) multiple(data acceptance.TestData) string {
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
