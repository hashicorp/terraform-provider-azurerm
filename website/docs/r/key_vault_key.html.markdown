---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_key"
description: |-
  Manages a Key Vault Key.

---

# azurerm_key_vault_key

Manages a Key Vault Key.

## Example Usage

~> **Note:** To use this resource, your client should have RBAC roles with permissions like `Key Vault Crypto Officer` or `Key Vault Administrator` or an assigned Key Vault Access Policy with permissions `Create`,`Delete`,`Get`,`Purge`,`Recover`,`Update` and `GetRotationPolicy` for keys without Rotation Policy. Include `SetRotationPolicy` for keys with Rotation Policy.

~> **Note:** The Azure Provider includes a Feature Toggle which will purge a Key Vault Key resource on destroy, rather than the default soft-delete. See [`purge_soft_deleted_keys_on_destroy`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block#purge_soft_deleted_keys_on_destroy) for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_keys_on_destroy = true
      recover_soft_deleted_keys          = true
    }
  }
}

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
      "Delete",
      "Get",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy"
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

  rotation_policy {
    automatic {
      time_before_expiry = "P30D"
    }

    expire_after         = "P90D"
    notify_before_expiry = "P29D"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Key. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of the Key Vault where the Key should be created. Changing this forces a new resource to be created.

* `key_type` - (Required) Specifies the Key Type to use for this Key Vault Key. Possible values are `EC` (Elliptic Curve), `EC-HSM`, `RSA` and `RSA-HSM`. Changing this forces a new resource to be created.

* `key_size` - (Optional) Specifies the Size of the RSA key to create in bytes. For example, 1024 or 2048. *Note*: This field is required if `key_type` is `RSA` or `RSA-HSM`. Changing this forces a new resource to be created.

* `curve` - (Optional) Specifies the curve to use when creating an `EC` key. Possible values are `P-256`, `P-256K`, `P-384`, and `P-521`. This field will be required in a future release if `key_type` is `EC` or `EC-HSM`. The API will default to `P-256` if nothing is specified. Changing this forces a new resource to be created.

* `key_opts` - (Required) A list of JSON web key operations. Possible values include: `decrypt`, `encrypt`, `sign`, `unwrapKey`, `verify` and `wrapKey`. Please note these values are case sensitive.

* `not_before_date` - (Optional) Key not usable before the provided UTC datetime (Y-m-d'T'H:M:S'Z').

~> **Note:** Once `expiration_date` is set, it's not possible to unset the key even if it is deleted & recreated as underlying Azure API uses the restore of the purged key.

* `expiration_date` - (Optional) Expiration UTC datetime (Y-m-d'T'H:M:S'Z').

~> **Note:** Removing this field from the config forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `rotation_policy` - (Optional) A `rotation_policy` block as defined below.

---

A `rotation_policy` block supports the following:

* `expire_after` - (Optional) Expire a Key Vault Key after given duration as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

* `automatic` - (Optional) An `automatic` block as defined below.

* `notify_before_expiry` - (Optional) Notify at a given duration before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

---

An `automatic` block supports the following:

* `time_after_creation` - (Optional) Rotate automatically at a duration after create as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

* `time_before_expiry` - (Optional) Rotate automatically at a duration before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Key Vault Key ID.
* `resource_id` - The (Versioned) ID for this Key Vault Key. This property points to a specific version of a Key Vault Key, as such using this won't auto-rotate values if used in other Azure Services.
* `resource_versionless_id` - The Versionless ID of the Key Vault Key. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Key is updated.
* `version` - The current version of the Key Vault Key.
* `versionless_id` - The Base ID of the Key Vault Key.
* `n` - The RSA modulus of this Key Vault Key.
* `e` - The RSA public exponent of this Key Vault Key.
* `x` - The EC X component of this Key Vault Key.
* `y` - The EC Y component of this Key Vault Key.
* `public_key_pem` - The PEM encoded public key of this Key Vault Key.
* `public_key_openssh` - The OpenSSH encoded public key of this Key Vault Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Key.
* `read` - (Defaults to 30 minutes) Used when retrieving the Key Vault Key.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Key.

## Import

Key Vault Key which is Enabled can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_key.example "https://example-keyvault.vault.azure.net/keys/example/fdf067c93bbb4b22bff4d8b7a9a56217"
```
