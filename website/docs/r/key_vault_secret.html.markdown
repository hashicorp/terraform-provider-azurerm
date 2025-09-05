---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_secret"
description: |-
  Manages a Key Vault Secret.

---

# azurerm_key_vault_secret

Manages a Key Vault Secret.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

~> **Note:** The Azure Provider includes a Feature Toggle which will purge a Key Vault Secret resource on destroy, rather than the default soft-delete. See [`purge_soft_deleted_secrets_on_destroy`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block#purge_soft_deleted_secrets_on_destroy) for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_secrets_on_destroy = true
      recover_soft_deleted_secrets          = true
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
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "example" {
  name         = "secret-sauce"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Secret. Changing this forces a new resource to be created.

* `value` - (Optional) Specifies the value of the Key Vault Secret. Changing this will create a new version of the Key Vault Secret.

* `value_wo` - (Optional, Write-Only) Specifies the value of the Key Vault Secret. Changing this will create a new version of the Key Vault Secret.

* ~> **Note:** One of `value` or `value_wo` must be specified.

* `value_wo_version` - (Optional) An integer value used to trigger an update for `value_wo`. This property should be incremented when updating `value_wo`.

~> **Note:** Key Vault strips newlines. To preserve newlines in multi-line secrets try replacing them with `\n` or by base 64 encoding them with `replace(file("my_secret_file"), "/\n/", "\n")` or `base64encode(file("my_secret_file"))`, respectively.

* `key_vault_id` - (Required) The ID of the Key Vault where the Secret should be created. Changing this forces a new resource to be created.

* `content_type` - (Optional) Specifies the content type for the Key Vault Secret.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `not_before_date` - (Optional) Key not usable before the provided UTC datetime (Y-m-d'T'H:M:S'Z').

* `expiration_date` - (Optional) Expiration UTC datetime (Y-m-d'T'H:M:S'Z').

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Key Vault Secret ID.
* `resource_id` - The (Versioned) ID for this Key Vault Secret. This property points to a specific version of a Key Vault Secret, as such using this won't auto-rotate values if used in other Azure Services.
* `resource_versionless_id` - The Versionless ID of the Key Vault Secret. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Secret is updated.
* `version` - The current version of the Key Vault Secret.
* `versionless_id` - The Base ID of the Key Vault Secret.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Secret.
* `read` - (Defaults to 30 minutes) Used when retrieving the Key Vault Secret.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Secret.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Secret.

## Import

Key Vault Secrets which are Enabled can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_secret.example "https://example-keyvault.vault.azure.net/secrets/example/fdf067c93bbb4b22bff4d8b7a9a56217"
```
