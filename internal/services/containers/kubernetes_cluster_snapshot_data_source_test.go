package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/managedclustersnapshots"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesClusterSnapshotDataSource struct{}

func TestAccDataSourceKubernetesClusterSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster_snapshot", "test")
	r := KubernetesClusterSnapshotDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.snapshotSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					client := clients.Containers.ManagedClusterSnapshotClient
					clusterId, err := managedclusters.ParseManagedClusterID(state.ID)
					if err != nil {
						return err
					}
					id := managedclustersnapshots.NewManagedClusterSnapshotID(clusterId.SubscriptionId, clusterId.ResourceGroupName, data.RandomString)
					snapshot := managedclustersnapshots.ManagedClusterSnapshot{
						Location: data.Locations.Primary,
						Properties: &managedclustersnapshots.ManagedClusterSnapshotProperties{
							SnapshotType: utils.ToPtr(managedclustersnapshots.SnapshotTypeManagedCluster),
							CreationData: &managedclustersnapshots.CreationData{
								SourceResourceId: utils.String(clusterId.ID()),
							},
						},
					}
					_, err = client.CreateOrUpdate(ctx, id, snapshot)
					if err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
					return nil
				}, "azurerm_kubernetes_cluster.source"),
			),
		},
		{
			Config: r.snapshotRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source_cluster_id").IsNotEmpty(),
			),
		},
		{
			Config: r.snapshotSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					client := clients.Containers.ManagedClusterSnapshotClient
					clusterId, err := managedclusters.ParseManagedClusterID(state.ID)
					if err != nil {
						return err
					}
					id := managedclustersnapshots.NewManagedClusterSnapshotID(clusterId.SubscriptionId, clusterId.ResourceGroupName, data.RandomString)
					_, err = client.Delete(ctx, id)
					if err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
					return nil
				}, "azurerm_kubernetes_cluster.source"),
			),
		},
	})
}

func (KubernetesClusterSnapshotDataSource) snapshotSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "source" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
  }
  identity {
    type = "SystemAssigned"
  }
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesClusterSnapshotDataSource) snapshotRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "source" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
  }
  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_kubernetes_cluster_snapshot" "test" {
  name                = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
}
 `, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
