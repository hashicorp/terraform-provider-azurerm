---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_capability_host"
description: |-
  Manages a Cognitive Account Capability Host.
---

# azurerm_cognitive_account_capability_host

Manages a Cognitive Account Capability Host.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmos"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
    zone_redundant    = false
  }
}

resource "azurerm_search_service" "example" {
  name                         = "example-search"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  sku                          = "basic"
  replica_count                = 1
  partition_count              = 1
  hosting_mode                 = "default"
  local_authentication_enabled = true

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacct"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-account"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "example-account"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account" "aisvc" {
  name                  = "example-aisvc"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  kind                  = "AIServices"
  sku_name              = "S0"
  custom_subdomain_name = "example-aisvc"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection_entra_id" "cosmos" {
  name                 = "example-conn-cosmos"
  cognitive_account_id = azurerm_cognitive_account.example.id
  category             = "CosmosDb"
  target               = azurerm_cosmosdb_account.example.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cosmosdb_account.example.id
    location   = azurerm_resource_group.example.location
  }
}

resource "azurerm_cognitive_account_connection_account_managed_identity" "storage" {
  name                 = "example-conn-storage"
  cognitive_account_id = azurerm_cognitive_account.example.id
  category             = "AzureStorageAccount"
  target               = azurerm_storage_account.example.primary_blob_endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_storage_account.example.id
    location   = azurerm_resource_group.example.location
  }
}

resource "azurerm_cognitive_account_connection_entra_id" "search" {
  name                 = "example-conn-search"
  cognitive_account_id = azurerm_cognitive_account.example.id
  category             = "CognitiveSearch"
  target               = "https://${azurerm_search_service.example.name}.search.windows.net"

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_search_service.example.id
    location   = azurerm_resource_group.example.location
  }
}

resource "azurerm_cognitive_account_connection_api_key" "aisvc" {
  name                 = "example-conn-aisvc"
  cognitive_account_id = azurerm_cognitive_account.example.id
  category             = "AIServices"
  target               = azurerm_cognitive_account.aisvc.endpoint
  api_key              = azurerm_cognitive_account.aisvc.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.aisvc.id
    location   = azurerm_resource_group.example.location
  }
}

resource "azurerm_cognitive_account_capability_host" "example" {
  name                       = "example-capability-host"
  cognitive_account_id       = azurerm_cognitive_account.example.id
  ai_services_connections    = [azurerm_cognitive_account_connection_api_key.aisvc.name]
  storage_connections        = [azurerm_cognitive_account_connection_account_managed_identity.storage.name]
  thread_storage_connections = [azurerm_cognitive_account_connection_entra_id.cosmos.name]
  vector_store_connections   = [azurerm_cognitive_account_connection_entra_id.search.name]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Account Capability Host. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Required) The ID of the Cognitive Account. Changing this forces a new resource to be created.

* `ai_services_connections` - (Optional) A list of AI Services connection names. Changing this forces a new resource to be created.

~> **Note:** A maximum of one AI Services connection can be specified.

* `storage_connections` - (Optional) A list of Storage connection names. Changing this forces a new resource to be created.

~> **Note:** A maximum of one Storage connection can be specified.

* `thread_storage_connections` - (Optional) A list of Thread Storage connection names. Changing this forces a new resource to be created.

~> **Note:** A maximum of one Thread Storage connection can be specified.

* `vector_store_connections` - (Optional) A list of Vector Store connection names. Changing this forces a new resource to be created.

~> **Note:** A maximum of one Vector Store connection can be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Account Capability Host.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Account Capability Host.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Account Capability Host.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Account Capability Host.

## Import

A Cognitive Account Capability Host can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_capability_host.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CognitiveServices/accounts/account1/capabilityHosts/capabilityHost1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2026-03-01
