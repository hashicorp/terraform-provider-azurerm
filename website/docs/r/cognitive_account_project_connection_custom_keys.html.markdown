---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_project_connection_custom_keys"
description: |-
  Manages a Cognitive Services (Microsoft Foundry) Account Project Connection with Custom Keys authentication.
---

# azurerm_cognitive_account_project_connection_custom_keys

Manages a Cognitive Services (Microsoft Foundry) Account Project Connection with Custom Keys authentication.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-cognitive-account"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "example-cognitive-account"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project" "example" {
  name                 = "example-cognitive-account-project"
  cognitive_account_id = azurerm_cognitive_account.example.id
  location             = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project_connection_custom_keys" "example" {
  name                         = "example-project-connection-keys"
  cognitive_account_project_id = azurerm_cognitive_account_project.example.id
  category                     = "CustomKeys"
  target                       = azurerm_cognitive_account.example.endpoint

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.example.primary_access_key
    secondaryKey = azurerm_cognitive_account.example.secondary_access_key
  }

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.example.id
    Location   = azurerm_cognitive_account.example.location
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Services Account Project Connection. Changing this forces a new resource to be created.

* `cognitive_account_project_id` - (Required) The ID of the Cognitive Services Account Project. Changing this forces a new resource to be created.

* `category` - (Required) The category of the connection. Possible values are `CustomKeys`, `RemoteA2A`, and `RemoteTool`. Changing this forces a new resource to be created.

* `custom_keys` - (Required) A mapping of custom keys for authentication. All values in this map are sensitive.

~> **Note:** Custom keys are stored in the Terraform state. Azure does not return these credentials when reading the connection.

* `target` - (Required) The target endpoint or resource for the connection.

* `metadata` - (Optional) A mapping of metadata key-value pairs for the connection.

-> **Note:** Only configured metadata keys are tracked. Additional metadata keys returned by Azure are not added to the Terraform state.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Services Account Project Connection with Custom Keys authentication.

* `authentication_type` - The authentication type of the connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Services Account Project Connection with Custom Keys authentication.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services Account Project Connection with Custom Keys authentication.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Services Account Project Connection with Custom Keys authentication.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Services Account Project Connection with Custom Keys authentication.

## Import

A Cognitive Services Account Project Connection with Custom Keys authentication can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_project_connection_custom_keys.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CognitiveServices/accounts/account1/projects/project1/connections/connection1
```

-> **Note:** Only connections using Custom Keys authentication can be imported. Configure `custom_keys` after import, as these values are not populated in the imported state.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2026-03-01
