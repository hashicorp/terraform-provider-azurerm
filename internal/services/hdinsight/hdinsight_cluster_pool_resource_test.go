// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2024-05-01/hdinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ClusterPoolResource struct{}

func (r ClusterPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := hdinsights.ParseClusterPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.HDInsight2024.Hdinsights.ClusterPoolsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retreiving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccHDInsightClusterPools_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_cluster_pool", "test")
	r := ClusterPoolResource{}

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

func TestAccHDInsightClusterPools_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_cluster_pool", "test")
	r := ClusterPoolResource{}

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

func TestAccHDInsightClusterPools_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_cluster_pool", "test")
	r := ClusterPoolResource{}

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
	})
}

func TestAccHDInsightClusterPools_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_cluster_pool", "test")
	r := ClusterPoolResource{}

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

func (r ClusterPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_hdinsight_cluster_pool" "test" {
  name                        = "accpool-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  managed_resource_group_name = "df433rg"

  cluster_pool_profile {
    version = "1.1"
  }
  compute_profile {
    vm_size = "Standard_F4s_v2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ClusterPoolResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[3]s-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_hdinsight_cluster_pool" "test" {
  name                        = "accpool-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  managed_resource_group_name = "df433rg"

  cluster_pool_profile {
    version = "1.1"
  }
  compute_profile {
    vm_size = "Standard_F4s_v2"
  }
  log_analytics_profile {
    enabled      = true
    workspace_id = azurerm_log_analytics_workspace.test.id
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r ClusterPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%[3]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[2]s-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_hdinsight_cluster_pool" "test" {
  name                = "accpool-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  compute_profile {
    vm_size = "Standard_D4a_v4"
  }
  cluster_pool_profile {
    version = "1.1"
  }
  log_analytics_profile {
    enabled      = true
    workspace_id = azurerm_log_analytics_workspace.test.id
  }
  network_profile {
    subnet_id            = azurerm_subnet.test.id
    outbound_type        = "loadBalancer"
    private_link_enabled = true
  }
  managed_resource_group_name = "df423rg"
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r ClusterPoolResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_hdinsight_cluster_pool" "import" {
  name                        = azurerm_hdinsight_cluster_pool.test.name
  resource_group_name         = azurerm_hdinsight_cluster_pool.test.resource_group_name
  location                    = azurerm_hdinsight_cluster_pool.test.location
  managed_resource_group_name = "df433rg"

  cluster_pool_profile {
    version = "1.1"
  }

  compute_profile {
    vm_size = "Standard_F4s_v2"
  }
}
`, template)
}

func (r ClusterPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, "westus2")
}
