---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection_aad"
description: |-
  Manages a Cognitive Services Account Connection with AAD authentication.
---

# azurerm_cognitive_account_connection_aad

Manages a Cognitive Services Account Connection with AAD (Azure Active Directory) authentication.

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

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "examplecsc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_cognitive_account_connection_aad" "example" {
  name                 = "example-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
  category             = "AzureBlob"
  target               = azurerm_storage_account.example.primary_blob_endpoint

  metadata = {
    accountName   = azurerm_storage_account.example.name
    containerName = azurerm_storage_container.example.name
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Services Account Connection. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Required) The ID of the Cognitive Services Account. Changing this forces a new resource to be created.

* `category` - (Required) The category of the connection. Possible values include `AzureBlob`, `AzureOpenAI`, and other supported connection categories. Changing this forces a new resource to be created.

* `metadata` - (Required) A mapping of metadata key-value pairs for the connection.

* `target` - (Required) The target endpoint or resource for the connection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Services Account Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Services Account Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services Account Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Services Account Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Services Account Connection.

## Import

Cognitive Services Account Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_connection_aad.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1/connections/connection1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2025-06-01
