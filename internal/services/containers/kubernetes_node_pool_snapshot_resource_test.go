package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type kubernetesNodePoolSnapshotResource struct{}

func TestAcckubernetesNodePoolSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_node_pool_snapshot", "test")
	r := kubernetesNodePoolSnapshotResource{}

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

func TestAcckubernetesNodePoolSnapshot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_node_pool_snapshot", "test")
	r := kubernetesNodePoolSnapshotResource{}

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

func TestAcckubernetesNodePoolSnapshot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_node_pool_snapshot", "test")
	r := kubernetesNodePoolSnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_node_pool_snapshot"),
		},
	})
}

func (t kubernetesNodePoolSnapshotResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NodePoolSnapshotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.SnapshotsClient.Get(ctx, id.ResourceGroup, id.SnapshotName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r kubernetesNodePoolSnapshotResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_node_pool_snapshot" "test" {
  name                = "acctestss%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  node_pool_id        = azurerm_kubernetes_cluster_node_pool.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r kubernetesNodePoolSnapshotResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_node_pool_snapshot" "test" {
  name                = "acctestss%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  node_pool_id        = azurerm_kubernetes_cluster_node_pool.test.id
  tags = {
    environment = "production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r kubernetesNodePoolSnapshotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_node_pool_snapshot" "import" {
  name                = azurerm_kubernetes_node_pool_snapshot.test.name
  location            = azurerm_kubernetes_node_pool_snapshot.test.location
  resource_group_name = azurerm_kubernetes_node_pool_snapshot.test.resource_group_name
  node_pool_id        = azurerm_kubernetes_node_pool_snapshot.test.node_pool_id
}
`, r.basic(data))
}

func (kubernetesNodePoolSnapshotResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
