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

type kubernetesClusterSnapshotResource struct{}

func TestAcckubernetesClusterSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_snapshot", "test")
	r := kubernetesClusterSnapshotResource{}

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

func TestAcckubernetesClusterSnapshot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_snapshot", "test")
	r := kubernetesClusterSnapshotResource{}

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

func TestAcckubernetesClusterSnapshot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_snapshot", "test")
	r := kubernetesClusterSnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_cluster_snapshot"),
		},
	})
}

func (t kubernetesClusterSnapshotResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ClusterSnapshotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.ManagedClusterSnapshotsClient.Get(ctx, id.ResourceGroup, id.ManagedclustersnapshotName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r kubernetesClusterSnapshotResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_snapshot" "test" {
  name                = "acctestss%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_id          = azurerm_kubernetes_cluster.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r kubernetesClusterSnapshotResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_snapshot" "test" {
  name                = "acctestss%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_id          = azurerm_kubernetes_cluster.test.id
  tags = {
    environment = "production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r kubernetesClusterSnapshotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_snapshot" "import" {
  name                = azurerm_kubernetes_cluster_snapshot.test.name
  resource_group_name = azurerm_kubernetes_cluster_snapshot.test.resource_group_name
  location            = azurerm_kubernetes_cluster_snapshot.test.location
  cluster_id          = azurerm_kubernetes_cluster_snapshot.test.cluster_id
}
`, r.basic(data))
}

func (kubernetesClusterSnapshotResource) template(data acceptance.TestData) string {
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
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
