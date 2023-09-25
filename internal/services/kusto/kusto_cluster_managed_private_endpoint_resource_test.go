// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/managedprivateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KustoClusterManagedPrivateEndpointResource struct{}

func TestAccKustoClusterManagedPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_managed_private_endpoint", "test")
	r := KustoClusterManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("group_id").HasValue("blob"),
			),
		},
		data.ImportStep(),
	},
	)
}

func TestAccKustoClusterManagedPrivateEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_managed_private_endpoint", "test")
	r := KustoClusterManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("group_id").HasValue("blob"),
				check.That(data.ResourceName).Key("request_message").HasValue("Please Approve"),
			),
		},
		data.ImportStep(),
	},
	)
}

func (KustoClusterManagedPrivateEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.ClusterManagedPrivateEndpointClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("response model is empty")
	}
	return utils.Bool(resp.Model.Properties != nil), nil
}

func (r KustoClusterManagedPrivateEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-kusto-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_kusto_cluster_managed_private_endpoint" "test" {
  name                     = "acctestmpe%d"
  resource_group_name      = azurerm_resource_group.rg.name
  cluster_name             = azurerm_kusto_cluster.test.name
  private_link_resource_id = azurerm_storage_account.test.id
  group_id                 = "blob"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}

func (r KustoClusterManagedPrivateEndpointResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-kusto-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_kusto_cluster_managed_private_endpoint" "test" {
  name                         = "acctestmpe%d"
  resource_group_name          = azurerm_resource_group.rg.name
  cluster_name                 = azurerm_kusto_cluster.test.name
  private_link_resource_id     = azurerm_storage_account.test.id
  private_link_resource_region = azurerm_storage_account.test.location
  group_id                     = "blob"
  request_message              = "Please Approve"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}
