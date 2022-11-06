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

resource "azurerm_key_vault_key" "example" {
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
  key_resource_versionless_id = azurerm_key_vault_key.example.resource_versionless_id

  expire_after         = "P61D"
  notify_before_expiry = "P8D"

  automatic {
    time_after_creation = "P31D"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `key_resource_versionless_id` - (Required) Specifies the Key Vault Key to be managed. Changing this forces a new Key Rotation Policy to be created.

* `expire_after` - (Optional) Expire a Key Vault Key after given duration as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). Is required if `notify_before_expiry` is set.

* `notify_before_expiry` - (Optional) Notify at a given time before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). Is required if `expire_after` is set.

* `automatic` - (Optional) An `automatic` block as defined below. Details within this block specify the automatic rotation of the Key Vault Key this policy is applied to.

---

An `automatic` block supports the following:

* `time_after_creation` - (Optional) Rotate automatically at a given duration after create as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

* `time_before_expiry` - (Optional) Rotate automatically at a given duration before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The Key Vault Key Rotation Policy ID.

* `resource_id` - The Resource ID of the Key Rotation Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Rotation Policy.
* `read` - (Defaults to 30 minutes) Used when retrieving the Key Rotation Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Key Rotation Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Rotation Policy.

## Import

Key Rotation Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_key_rotation_policy.example "https://example-keyvault.vault.azure.net/keys/key1/rotationpolicy"
```
