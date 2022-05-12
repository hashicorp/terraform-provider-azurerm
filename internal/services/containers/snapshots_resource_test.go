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

type SnapshotsResource struct{}

func TestAccSnapshots_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshots", "test")
	r := SnapshotsResource{}

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

func TestAccSnapshots_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshots", "test")
	r := SnapshotsResource{}

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

func TestAccSnapshots_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_snapshots", "test")
	r := SnapshotsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_snapshots"),
		},
	})
}

func (t SnapshotsResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SnapshotsID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.SnapshotsClient.Get(ctx, id.ResourceGroup, id.SnapshotName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r SnapshotsResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_snapshots" "test" {
  name                = "acctestss%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  source_resource_id  = azurerm_kubernetes_cluster_node_pool.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r SnapshotsResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_snapshots" "test" {
  name                = "acctestss%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  source_resource_id  = azurerm_kubernetes_cluster_node_pool.test.id
  snapshot_type       = "NodePool"
}
`, r.template(data), data.RandomInteger)
}

func (r SnapshotsResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_snapshots" "import" {
  name                = azurerm_snapshots.test.name
  location            = azurerm_snapshots.test.location
  resource_group_name = azurerm_snapshots.test.resource_group_name
  source_resource_id  = azurerm_snapshots.test.source_resource_id
}
`, r.basic(data))
}

func (SnapshotsResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
