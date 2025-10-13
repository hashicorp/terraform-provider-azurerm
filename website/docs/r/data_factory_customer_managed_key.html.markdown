---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_customer_managed_key"
description: |-
  Manages a Customer Managed Key for a Data Factory.
---

# azurerm_data_factory_customer_managed_key

Manages a Customer Managed Key for a Data Factory.

~> **Note:** The Customer Managed Key cannot be removed from the Data Factory once added. To remove the Customer Managed Key delete and recreate the parent Data Factory.

## Example Usage with System Assigned Identity

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                     = "example-key-vault"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_key" "example" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "unwrapKey",
    "wrapKey"
  ]
}

resource "azurerm_key_vault_access_policy" "current_client_policy" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Delete",
    "Get",
    "Set",
  ]
}

resource "azurerm_data_factory" "example" {
  name                = "example_data_factory"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }
}

resource "azurerm_key_vault_access_policy" "datafactory" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_data_factory.example.identity[0].tenant_id
  object_id    = azurerm_data_factory.example.identity[0].principal_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "GetRotationPolicy",
    "WrapKey",
    "UnwrapKey",
  ]

  secret_permissions = [
    "Delete",
    "Get",
    "Set",
  ]
}

resource "azurerm_data_factory_customer_managed_key" "example" {
  data_factory_id         = azurerm_data_factory.example.id
  customer_managed_key_id = azurerm_key_vault_key.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `data_factory_id` - (Required) The ID of the Data Factory Resource the Customer Managed Key will be associated with. Changing this forces a new resource to be created.

* `customer_managed_key_id` - (Required) The ID the of the Customer Managed Key to associate with the Data Factory.

* `user_assigned_identity_id` - (Optional) The User Assigned Identity ID that will be used to access Key Vaults that contain the encryption keys.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Customer Managed Key.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Customer Managed Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Customer Managed Key.

## Import

Data Factory Customer Managed Keys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataFactory` - 2018-06-01
