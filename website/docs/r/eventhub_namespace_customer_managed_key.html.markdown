---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace_customer_managed_key"
description: |-
  Manages a Customer Managed Key for a EventHub Namespace.
---

# azurerm_eventhub_namespace_customer_managed_key

Manages a Customer Managed Key for a EventHub Namespace.

!> **Note:** In 2.x versions of the Azure Provider during deletion this resource will **delete and recreate the parent EventHub Namespace which may involve data loss** as it's not possible to remove the Customer Managed Key from the EventHub Namespace once it's been added. Version 3.0 of the Azure Provider will change this so that the Delete operation is a noop, requiring the parent EventHub Namespace is deleted/recreated to remove the Customer Managed Key.

## Example Usage with System Assigned Identity

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Dedicated_1"
}

resource "azurerm_eventhub_namespace" "example" {
  name                 = "example-namespace"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  sku                  = "Standard"
  dedicated_cluster_id = azurerm_eventhub_cluster.example.id

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                     = "examplekv"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "example" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_eventhub_namespace.example.identity[0].tenant_id
  object_id    = azurerm_eventhub_namespace.example.identity[0].principal_id

  key_permissions = ["Get", "UnwrapKey", "WrapKey"]
}

resource "azurerm_key_vault_access_policy" "example2" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "List",
    "Purge",
    "Recover",
    "GetRotationPolicy",
  ]
}

resource "azurerm_key_vault_key" "example" {
  name         = "examplekvkey"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.example,
    azurerm_key_vault_access_policy.example2,
  ]
}

resource "azurerm_eventhub_namespace_customer_managed_key" "example" {
  eventhub_namespace_id = azurerm_eventhub_namespace.example.id
  key_vault_key_ids     = [azurerm_key_vault_key.example.id]
}
```

## Example Usage with User Assigned Identity

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Dedicated_1"
}

resource "azurerm_user_assigned_identity" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_eventhub_namespace" "example" {
  name                 = "example-namespace"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  sku                  = "Standard"
  dedicated_cluster_id = azurerm_eventhub_cluster.example.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                     = "examplekv"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "example" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = ["Get", "UnwrapKey", "WrapKey"]
}

resource "azurerm_key_vault_access_policy" "example2" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "List",
    "Purge",
    "Recover",
    "GetRotationPolicy",
  ]
}

resource "azurerm_key_vault_key" "example" {
  name         = "examplekvkey"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.example,
    azurerm_key_vault_access_policy.example2,
  ]
}

resource "azurerm_eventhub_namespace_customer_managed_key" "example" {
  eventhub_namespace_id     = azurerm_eventhub_namespace.example.id
  key_vault_key_ids         = [azurerm_key_vault_key.example.id]
  user_assigned_identity_id = azurerm_user_assigned_identity.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `eventhub_namespace_id` - (Required) The ID of the EventHub Namespace. Changing this forces a new resource to be created.

* `key_vault_key_ids` - (Required) The list of keys of Key Vault.

* `infrastructure_encryption_enabled` - (Optional) Whether to enable Infrastructure Encryption (Double Encryption). Changing this forces a new resource to be created.

* `user_assigned_identity_id` - (Optional) The ID of a User Managed Identity that will be used to access Key Vaults that contain the encryption keys.

~> **Note:** If using `user_assigned_identity_id`, ensure the User Assigned Identity is also assigned to the parent Event Hub.

~> **Note:** If using `user_assigned_identity_id`, make sure to assign the identity the appropriate permissions to access the Key Vault key. Failure to grant `Get, UnwrapKey, and WrapKey` will cause this resource to fail to apply.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the EventHub Namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Namespace Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Namespace Customer Managed Key.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Namespace Customer Managed Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Namespace Customer Managed Key.

## Import

Customer Managed Keys for a EventHub Namespace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_namespace_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1
```
