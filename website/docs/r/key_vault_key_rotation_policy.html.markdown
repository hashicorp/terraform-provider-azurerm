---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_key_rotation_policy"
description: |-
  Manages a Key Rotation Policy.
---

# azurerm_key_vault_key_rotation_policy

Manages a Key Rotation Policy.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                       = "examplekeyvault"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
      "Purge",
      "Recover"
    ]

    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "generated" {
  name         = "generated-certificate"
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

resource "azurerm_key_vault_key_rotation_policy" "example" {
  key_vault_id = azurerm_key_vault_key.generated.id
  key_name     = azurerm_key_vault_key.generated.name

  expiry_time       = "P61D"
  notification_time = "P8D"

  auto_rotation {
    time_after_create = "P31D"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `key_name` - (Required) Specifies the name of the Key Vault Key managed. Changing this forces a new Key Rotation Policy to be created.

* `key_vault_id` - (Required) The ID of the Key Vault where the Key Vault Key Policy is applied. Changing this forces a new Key Rotation Policy to be created.

* `expiry_time` - (Required) Expire a Key Vault Key after given time as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

---

* `auto_rotation` - (Optional) A `auto_rotation` block as defined below.

* `notification_time` - (Optional) Notify at a given time before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

---

A `auto_rotation` block supports the following:

* `time_after_create` - (Optional) Rotate automatically at a given time after create as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

* `time_before_expiry` - (Optional) Rotate automatically at a given time before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Key Rotation Policy.

* `resource_id` - The ID of the TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Rotation Policy.
* `read` - (Defaults to 30 minutes) Used when retrieving the Key Rotation Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Key Rotation Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Rotation Policy.

## Import

Key Rotation Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_key_rotation_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/keys/key1/rotationpolicy
```
