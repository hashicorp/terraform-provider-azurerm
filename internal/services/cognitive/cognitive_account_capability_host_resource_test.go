// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountcapabilityhost"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountCapabilityHostTestResource struct{}

func TestAccCognitiveAccountCapabilityHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_capability_host", "test")
	r := CognitiveAccountCapabilityHostTestResource{}

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

func TestAccCognitiveAccountCapabilityHost_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_capability_host", "test")
	r := CognitiveAccountCapabilityHostTestResource{}

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

func TestAccCognitiveAccountCapabilityHost_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_capability_host", "test")
	r := CognitiveAccountCapabilityHostTestResource{}

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

func (r CognitiveAccountCapabilityHostTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accountcapabilityhost.ParseCapabilityHostID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cognitive.AccountCapabilityHostClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CognitiveAccountCapabilityHostTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_account_capability_host" "test" {
  name                 = "acctest-ch-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountCapabilityHostTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_capability_host" "import" {
  name                 = azurerm_cognitive_account_capability_host.test.name
  cognitive_account_id = azurerm_cognitive_account_capability_host.test.cognitive_account_id
}
`, r.basic(data))
}

func (r CognitiveAccountCapabilityHostTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[2]d"
  address_space       = ["192.168.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["192.168.1.0/24"]

  delegation {
    name = "app-delegation"

    service_delegation {
      name = "Microsoft.App/environments"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-cosmos-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
    zone_redundant    = false
  }
}

resource "azurerm_search_service" "test" {
  name                         = "acctest-search-%[3]s"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  sku                          = "basic"
  replica_count                = 1
  partition_count              = 1
  hosting_mode                 = "default"
  local_authentication_enabled = true

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestcsc%[3]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_cognitive_account" "aisvc" {
  name                  = "acctest-aisvc-%[2]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  kind                  = "AIServices"
  sku_name              = "S0"
  custom_subdomain_name = "acctest-aisvc-%[2]d"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection" "cosmos" {
  name                 = "acctest-conn-cosmos-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "AAD"
  category             = "CosmosDb"
  target               = azurerm_cosmosdb_account.test.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cosmosdb_account.test.id
    location   = azurerm_resource_group.test.location
  }
}

resource "azurerm_cognitive_account_connection" "storage" {
  name                 = "acctest-conn-storage-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "AAD"
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint
  metadata = {
    containerName = azurerm_storage_container.test.name
    accountName   = azurerm_storage_account.test.name
  }
}

resource "azurerm_cognitive_account_connection" "search" {
  name                 = "acctest-conn-search-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "AAD"
  category             = "CognitiveSearch"
  target               = "https://${azurerm_search_service.test.name}.search.windows.net"

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_search_service.test.id
    location   = azurerm_resource_group.test.location
  }
}

resource "azurerm_cognitive_account_connection" "aisvc" {
  name                 = "acctest-conn-aisvc-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "ApiKey"
  category             = "AIServices"
  target               = azurerm_cognitive_account.aisvc.endpoint
  api_key              = azurerm_cognitive_account.aisvc.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.aisvc.id
    location   = azurerm_resource_group.test.location
  }
}

resource "azurerm_cognitive_account_connection" "aisvc2" {
  name                 = "acctest-conn-aisvc-2-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "ApiKey"
  category             = "AIServices"
  target               = azurerm_cognitive_account.aisvc.endpoint
  api_key              = azurerm_cognitive_account.aisvc.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.aisvc.id
    location   = azurerm_resource_group.test.location
  }
}

resource "azurerm_cognitive_account_capability_host" "test" {
  name                       = "acctest-ch-%[2]d"
  cognitive_account_id       = azurerm_cognitive_account.test.id
  subnet_id                  = azurerm_subnet.test.id
  ai_services_connections    = [azurerm_cognitive_account_connection.aisvc.name]
  storage_connections        = [azurerm_cognitive_account_connection.storage.name]
  thread_storage_connections = [azurerm_cognitive_account_connection.cosmos.name]
  vector_store_connections   = [azurerm_cognitive_account_connection.search.name]
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CognitiveAccountCapabilityHostTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cognitive_account" "test" {
  name                       = "acctest-cogacc-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctest-cogacc-%[1]d"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
