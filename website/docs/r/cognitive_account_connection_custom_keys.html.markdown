---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection_custom_keys"
description: |-
  Manages a Cognitive Services (Microsoft Foundry) Account Connection with Custom Keys authentication.
---

# azurerm_cognitive_account_connection_custom_keys

Manages a Cognitive Services (Microsoft Foundry) Account Connection with Custom Keys authentication.

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

resource "azurerm_cognitive_account" "openai" {
  name                = "example-openai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "OpenAI"
  sku_name            = "S0"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection_custom_keys" "example" {
  name                 = "example-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai.endpoint

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai.id
    location   = azurerm_cognitive_account.openai.location
  }

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai.secondary_access_key
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Services Account Connection. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Required) The ID of the Cognitive Services Account. Changing this forces a new resource to be created.

* `category` - (Required) The category of the connection. Possible values include `AzureBlob`, `AzureOpenAI`, `CustomKeys`, and other supported connection categories. Changing this forces a new resource to be created.

* `metadata` - (Required) A mapping of metadata key-value pairs for the connection.

* `target` - (Required) The target endpoint or resource for the connection.

* `custom_keys` - (Required) A mapping of custom keys for authentication. All values in this map are sensitive.

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
terraform import azurerm_cognitive_account_connection_custom_keys.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1/connections/connection1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2025-06-01
