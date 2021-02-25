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
      "create",
      "get",
    ]

    secret_permissions = [
      "set",
      "get",
      "delete",
      "purge",
      "recover"
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

* `value` - (Required) Specifies the value of the Key Vault Secret.

~> **Note:** Key Vault strips newlines. To preserve newlines in multi-line secrets try replacing them with `\n` or by base 64 encoding them with `replace(file("my_secret_file"), "/\n/", "\n")` or `base64encode(file("my_secret_file"))`, respectively.

* `key_vault_id` - (Required) The ID of the Key Vault where the Secret should be created.

* `content_type` - (Optional) Specifies the content type for the Key Vault Secret.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `not_before_date` - (Optional) Key not usable before the provided UTC datetime (Y-m-d'T'H:M:S'Z').

* `expiration_date` - (Optional) Expiration UTC datetime (Y-m-d'T'H:M:S'Z').

## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Secret ID.
* `version` - The current version of the Key Vault Secret.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Secret.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Secret.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Secret.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Secret.

## Import

Key Vault Secrets which are Enabled can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_secret.example "https://example-keyvault.vault.azure.net/secrets/example/fdf067c93bbb4b22bff4d8b7a9a56217"
```
