---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace"
description: |-
  Manages a Service Bus Namespace Customer Managed Key.
---

# azurerm_servicebus_namespace_customer_managed_key

Manages a Service Bus Namespace Customer Managed Key.

!> **Note:** It is not possible to remove the Customer Managed Key from the Service Bus Namespace once it's been added. To remove the Customer Managed Key, the parent Service Bus Namespace must be deleted and recreated.

-> **Note:** This resource should only be used to create a Customer Managed Key for Service Bus Namespaces with System Assigned identities. The `customer_managed_key` block in `azurerm_servicebus_namespace` should be used to create a Customer Managed Key for a Service Bus Namespace with a User Assigned identity.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "example" {
  name                         = "example-servicebus-namespace"
  location                     = azurerm_resource_group.example.location
  resource_group_name          = azurerm_resource_group.example.name
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = 1

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = [customer_managed_key]
  }
}


data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                        = "example-key-vault"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_servicebus_namespace.example.identity[0].tenant_id
    object_id = azurerm_servicebus_namespace.example.identity[0].principal_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "example-key-vault-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_servicebus_namespace_customer_managed_key" "example" {
  namespace_id     = azurerm_servicebus_namespace.example.id
  key_vault_key_id = azurerm_key_vault_key.example.id
}

```

## Argument Reference

The following arguments are supported:

* `namespace_id` - (Required) The ID of the Service Bus namespace. Changing this forces a new resource to be created.

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this Service Bus Namespace.

* `infrastructure_encryption_enabled` - (Optional) Used to specify whether enable Infrastructure Encryption. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Service Bus Namespace ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Bus Namespace Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Bus Namespace Customer Managed Key.
* `update` - (Defaults to 30 minutes) Used when updating the Service Bus Namespace Customer Managed Key.
* `delete` - (Defaults to 5 minutes) Used when deleting the Service Bus Namespace Customer Managed Key.

## Import

Service Bus Namespace Customer Managed Key can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceBus/namespaces/sbns1
```
