package extendedlocation_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"
)

type CustomLocationResource struct{}

func (r CustomLocationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := customlocations.ParseCustomLocationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ExtendedLocation.CustomLocationsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccExtendedLocationCustomLocations_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_extended_custom_locations", "test")
	r := CustomLocationResource{}

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

func (r CustomLocationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment" "admin" {
  scope                = azurerm_kubernetes_cluster.test.id
  role_definition_name = "Azure Kubernetes Service RBAC Cluster Admin"
  principal_id         = "51dfe1e8-70c6-4de5-a08e-e18aff23d815"
}

resource "azurerm_extended_custom_locations" "test" {
  name = "acctestcustomlocation%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  cluster_extension_ids = [
	"${azurerm_kubernetes_cluster.test.id}/providers/Microsoft.KubernetesConfiguration/extensions/foo"
  ]
  display_name = "customlocation%[2]d"
  namespace = "namespace%[2]d"
  host_resource_id = azurerm_kubernetes_cluster.test.id
}
`, template, data.RandomInteger)
}

func (r CustomLocationResource) template(data acceptance.TestData) string {
	data.Locations.Primary = "westus2"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-%[1]d"
  location = "%[2]s"
}

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"
  kubernetes_version  = "1.24.3"
  run_command_enabled = true

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_ds2_v2"
  }
  identity {
    type = "SystemAssigned"
  }
  tags = {
    ENV = "Test1"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
