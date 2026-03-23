---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_online_endpoint"
description: |-
  Manages a Machine Learning Online Endpoint.
---

# azurerm_machine_learning_online_endpoint

Manages a Machine Learning Online Endpoint.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-application-insights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_key_vault" "example" {
  name                = "example-key-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-machine-learning-workspac"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_machine_learning_online_endpoint" "example" {
  name                          = "example-machine-learning-online"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.example.id
  location                      = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Machine Learning Online Endpoint. Changing this forces a new resource to be created.

-> **Note:** The `name` must start with a letter, contain only letters, digits, or dashes, end with a letter or digit, and be between 3 and 32 characters.

* `machine_learning_workspace_id` - (Required) The ID of the Machine Learning Workspace in which the Online Endpoint should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Machine Learning Online Endpoint should exist. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new resource to be created.

* `authentication_mode` - (Optional) The authentication mode for the Machine Learning Online Endpoint. Possible values are `AADToken`, `AMLToken`, and `Key`. Defaults to `Key`.

* `description` - (Optional) The description of the Machine Learning Online Endpoint.

~> **Note:** The `description` cannot be changed after the Machine Learning Online Endpoint has been created.

* `properties` - (Optional) A mapping of properties to assign to the Machine Learning Online Endpoint. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Should `public network access` be enabled? Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Online Endpoint.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Machine Learning Online Endpoint. Possible values are `SystemAssigned` and `UserAssigned`. Changing this forces a new resource to be created.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Machine Learning Online Endpoint. Changing this forces a new resource to be created.

~> **Note:** This is required when `type` is set to `UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Machine Learning Online Endpoint.

* `rest_endpoint` - The REST endpoint of the Machine Learning Online Endpoint.

* `swagger_uri` - The Swagger URI of the Machine Learning Online Endpoint.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Online Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Online Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Online Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Online Endpoint.

## Import

A Machine Learning Online Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_online_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/onlineEndpoints/onlineEndpoint1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.MachineLearningServices` - 2025-06-01
