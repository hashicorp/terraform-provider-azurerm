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

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                = "examplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "get",
    ]

    secret_permissions = [
      "set",
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Key. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of the Key Vault where the Key should be created. Changing this forces a new resource to be created.

* `key_type` - (Required) Specifies the Key Type to use for this Key Vault Key. Possible values are `EC` (Elliptic Curve), `EC-HSM`, `Oct` (Octet), `RSA` and `RSA-HSM`. Changing this forces a new resource to be created.

* `key_size` - (Optional) Specifies the Size of the RSA key to create in bytes. For example, 1024 or 2048. *Note*: This field is required if `key_type` is `RSA` or `RSA-HSM`. Changing this forces a new resource to be created.

* `curve` - (Optional) Specifies the curve to use when creating an `EC` key. Possible values are `P-256`, `P-384`, `P-521`, and `SECP256K1`. This field will be required in a future release if `key_type` is `EC` or `EC-HSM`. The API will default to `P-256` if nothing is specified. Changing this forces a new resource to be created.

* `key_opts` - (Required) A list of JSON web key operations. Possible values include: `decrypt`, `encrypt`, `sign`, `unwrapKey`, `verify` and `wrapKey`. Please note these values are case sensitive.

* `not_before_date` - (Optional) Key not usable before the provided UTC datetime (Y-m-d'T'H:M:S'Z').

* `expiration_date` - (Optional) Expiration UTC datetime (Y-m-d'T'H:M:S'Z').

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Key ID.
* `version` - The current version of the Key Vault Key.
* `n` - The RSA modulus of this Key Vault Key.
* `e` - The RSA public exponent of this Key Vault Key.
* `x` - The EC X component of this Key Vault Key.
* `y` - The EC Y component of this Key Vault Key.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Key.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Key.

## Import

Key Vault Key which is Enabled can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_key.example "https://example-keyvault.vault.azure.net/keys/example/fdf067c93bbb4b22bff4d8b7a9a56217"
```
