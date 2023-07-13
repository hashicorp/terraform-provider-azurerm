// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vmware_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VmwareClusterResource struct{}

func TestAccVmwareCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_cluster", "test")
	r := VmwareClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cluster_number").Exists(),
				check.That(data.ResourceName).Key("hosts.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVmwareCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_cluster", "test")
	r := VmwareClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVmwareCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_cluster", "test")
	r := VmwareClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cluster_number").Exists(),
				check.That(data.ResourceName).Key("hosts.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cluster_number").Exists(),
				check.That(data.ResourceName).Key("hosts.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cluster_number").Exists(),
				check.That(data.ResourceName).Key("hosts.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (VmwareClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := clusters.ParseClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Vmware.ClusterClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r VmwareClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vmware_cluster" "test" {
  name               = "acctest-Cluster-%d"
  vmware_cloud_id    = azurerm_vmware_private_cloud.test.id
  cluster_node_count = 3
  sku_name           = "av36"
}
`, VmwarePrivateCloudResource{}.basic(data), data.RandomInteger)
}

func (r VmwareClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vmware_cluster" "import" {
  name               = azurerm_vmware_cluster.test.name
  vmware_cloud_id    = azurerm_vmware_cluster.test.vmware_cloud_id
  cluster_node_count = azurerm_vmware_cluster.test.cluster_node_count
  sku_name           = azurerm_vmware_cluster.test.sku_name
}
`, r.basic(data))
}

func (r VmwareClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vmware_cluster" "test" {
  name               = "acctest-Cluster-%d"
  vmware_cloud_id    = azurerm_vmware_private_cloud.test.id
  cluster_node_count = 4
  sku_name           = "av36"
}
`, VmwarePrivateCloudResource{}.basic(data), data.RandomInteger)
}
