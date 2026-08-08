---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_project_connection_account_key"
description: |-
  Manages a Cognitive Services (Microsoft Foundry) Project Connection with Account Key authentication.
---

# azurerm_cognitive_account_project_connection_account_key

Manages a Cognitive Services (Microsoft Foundry) Project Connection with Account Key authentication.

-> **Note:** In the new Foundry portal experience, "Project Connections" are shown as "Connected resources" in a project.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-aiservices"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "exampleaiservices"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project" "example" {
  name                 = "example-project"
  cognitive_account_id = azurerm_cognitive_account.example.id
  location             = azurerm_resource_group.example.location

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

resource "azurerm_cognitive_account_project_connection_account_key" "example" {
  name                         = "example-connection"
  cognitive_account_project_id = azurerm_cognitive_account_project.example.id
  category                     = "AzureStorageAccount"
  target                       = azurerm_storage_account.example.primary_blob_endpoint
  account_key                  = azurerm_storage_account.example.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_storage_account.example.id
    Location   = azurerm_storage_account.example.location
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Services Project Connection. Changing this forces a new resource to be created.

* `cognitive_account_project_id` - (Required) The ID of the Cognitive Services Project. Changing this forces a new resource to be created.

* `account_key` - (Required) The account key used for authentication.

* `category` - (Required) The category of the connection. Changing this forces a new resource to be created.

-> **Note:** The only possible value is `AzureStorageAccount`.

* `metadata` - (Required) A mapping of metadata key-value pairs for the connection.

~> **Note:** The `metadata` map must include a non-empty `ResourceId` value. To check any additional metadata returned by Azure, create an equivalent connection in the Foundry portal and inspect it with `az rest --method get --url "{connection_resource_id}?api-version=2026-03-01"`.

* `target` - (Required) The target endpoint or resource for the connection.

~> **Note:** `target` must be an absolute HTTP or HTTPS URL.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Services (Microsoft Foundry) Project Connection with Account Key authentication.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Services (Microsoft Foundry) Project Connection with Account Key authentication.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services (Microsoft Foundry) Project Connection with Account Key authentication.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Services (Microsoft Foundry) Project Connection with Account Key authentication.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Services (Microsoft Foundry) Project Connection with Account Key authentication.

## Import

A Cognitive Services (Microsoft Foundry) Project Connection with Account Key authentication can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_project_connection_account_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CognitiveServices/accounts/account1/projects/project1/connections/connection1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2026-03-01
