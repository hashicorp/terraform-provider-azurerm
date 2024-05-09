// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2024-05-01/hdinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"testing"
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
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_cluster_pools", "test")
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
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_cluster_pools", "test")
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

func TestAccHDInsightClusterPools_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_cluster_pools", "test")
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

resource "azurerm_hdinsight_cluster_pools" "test" {
  name 			  = "acctestpool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location 		  = azurerm_resource_group.test.location
  
  compute_profile {
	vm_size = "Standard_D3_v2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ClusterPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_virtual_network" "test" {
	name                = "acctestvnet-%[3]d"
	resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	address_space       = ["10.0.0.0/16"]

	subnet {
	  name           = "acctestsubnet-%[3]d"
	  address_prefix = "10.0.2.0/24"
	}
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[2]s-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_hdinsight_cluster_pools" "test" {
	  name 			  = "acctestpool-%[3]d"
	  resource_group_name = azurerm_resource_group.test.name
	  location 		  = azurerm_resource_group.test.location

	  compute_profile {
		vm_size = "Standard_D3_v2"
	 }
	 cluster_pool_profile {
		cluster_version = "1.2"
	}
	log_analytics_profile {
		log_analytics_profile_enabled = true
		workspace_id = azurerm_log_analytics_workspace.test.id
	}
	network_profile {
		subnet_id = azurerm_virtual_network.test.subnet[0].id
		outbound_type = "loadBalancer"
		private_api_services_enabled = true
	}
	 managed_resource_group_name = "df423rg"
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r ClusterPoolResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_hdinsight_cluster_pools" "import" {
	  name 			  = azurerm_hdinsight_cluster_pools.test.name
	  resource_group_name = azurerm_hdinsight_cluster_pools.test.resource_group_name
	  location 		  = azurerm_hdinsight_cluster_pools.test.location

	  compute_profile {
		vm_size = "Standard_D3_v2"
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
`, data.RandomInteger, data.Locations.Primary)
}
