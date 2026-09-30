// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package discovery_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DiscoveryBookshelfResource struct{}

func TestAccDiscoveryBookshelf_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf", "test")
	r := DiscoveryBookshelfResource{}

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

func TestAccDiscoveryBookshelf_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf", "test")
	r := DiscoveryBookshelfResource{}

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

func TestAccDiscoveryBookshelf_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf", "test")
	r := DiscoveryBookshelfResource{}

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

func TestAccDiscoveryBookshelf_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf", "test")
	r := DiscoveryBookshelfResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiscoveryBookshelf_customerManagedKeyRequiresLogAnalyticsCluster(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf", "test")
	r := DiscoveryBookshelfResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.customerManagedKeyWithoutLogAnalyticsCluster(data),
			ExpectError: regexp.MustCompile("`log_analytics_cluster_id` must be specified when `customer_managed_key` is set"),
		},
	})
}

func (r DiscoveryBookshelfResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := bookshelves.ParseBookshelfID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Discovery.Bookshelves.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DiscoveryBookshelfResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "private_endpoint" {
  name                 = "privateEndpointSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_subnet" "search" {
  name                 = "searchSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/23"]

  delegation {
    name = "Microsoft.App/environments"

    service_delegation {
      name    = "Microsoft.App/environments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DiscoveryBookshelfResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_discovery_bookshelf" "test" {
  name                       = "acctest-%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  private_endpoint_subnet_id = azurerm_subnet.private_endpoint.id
  search_subnet_id           = azurerm_subnet.search.id
  workload_identity_ids      = [azurerm_user_assigned_identity.test.id]
}
`, r.template(data), data.RandomInteger)
}

func (r DiscoveryBookshelfResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_discovery_bookshelf" "import" {
  name                       = azurerm_discovery_bookshelf.test.name
  resource_group_name        = azurerm_discovery_bookshelf.test.resource_group_name
  location                   = azurerm_discovery_bookshelf.test.location
  private_endpoint_subnet_id = azurerm_discovery_bookshelf.test.private_endpoint_subnet_id
  search_subnet_id           = azurerm_discovery_bookshelf.test.search_subnet_id
  workload_identity_ids      = azurerm_discovery_bookshelf.test.workload_identity_ids
}
`, r.basic(data))
}

func (r DiscoveryBookshelfResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_discovery_bookshelf" "test" {
  name                          = "acctest-%[2]d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  private_endpoint_subnet_id    = azurerm_subnet.private_endpoint.id
  public_network_access_enabled = false
  search_subnet_id              = azurerm_subnet.search.id
  workload_identity_ids         = [azurerm_user_assigned_identity.test.id]

  tags = {
    environment = "Testing"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DiscoveryBookshelfResource) customerManagedKeyWithoutLogAnalyticsCluster(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_discovery_bookshelf" "test" {
  name                       = "acctest-%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  private_endpoint_subnet_id = azurerm_subnet.private_endpoint.id
  search_subnet_id           = azurerm_subnet.search.id
  workload_identity_ids      = [azurerm_user_assigned_identity.test.id]

  customer_managed_key {
    identity_client_id = "00000000-0000-0000-0000-000000000000"
    key_name           = "acctestkey"
    key_vault_uri      = "https://acctestkv-%[2]d.vault.azure.net/"
  }
}
`, r.template(data), data.RandomInteger)
}
