// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/snapshots"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesNodePoolSnapshotDataSource struct{}

func TestAccDataSourceKubernetesNodePoolSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_node_pool_snapshot", "test")
	r := KubernetesNodePoolSnapshotDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.snapshotSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					if _, ok := ctx.Deadline(); !ok {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
						defer cancel()
					}
					client := clients.Containers.SnapshotClient
					poolId, err := agentpools.ParseAgentPoolID(state.ID)
					if err != nil {
						return err
					}
					id := snapshots.NewSnapshotID(poolId.SubscriptionId, poolId.ResourceGroupName, data.RandomString)
					snapshot := snapshots.Snapshot{
						Location: data.Locations.Primary,
						Properties: &snapshots.SnapshotProperties{
							CreationData: &snapshots.CreationData{
								SourceResourceId: utils.String(poolId.ID()),
							},
						},
					}
					_, err = client.CreateOrUpdate(ctx, id, snapshot)
					if err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
					return nil
				}, "azurerm_kubernetes_cluster_node_pool.source"),
			),
		},
		{
			Config: r.snapshotRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source_node_pool_id").IsNotEmpty(),
			),
		},
		{
			Config: r.snapshotSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					if _, ok := ctx.Deadline(); !ok {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
						defer cancel()
					}
					client := clients.Containers.SnapshotClient
					poolId, err := agentpools.ParseAgentPoolID(state.ID)
					if err != nil {
						return err
					}
					id := snapshots.NewSnapshotID(poolId.SubscriptionId, poolId.ResourceGroupName, data.RandomString)
					_, err = client.Delete(ctx, id)
					if err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
					return nil
				}, "azurerm_kubernetes_cluster_node_pool.source"),
			),
		},
	})
}

func (KubernetesNodePoolSnapshotDataSource) snapshotSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "source" {
  name                  = "source"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesNodePoolSnapshotDataSource) snapshotRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "source" {
  name                  = "source"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
}


data "azurerm_kubernetes_node_pool_snapshot" "test" {
  name                = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
}
 `, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
