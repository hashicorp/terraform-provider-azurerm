---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_secrets"
description: |-
  Gets a list of secret names from an existing Key Vault Secret.
---

# Data Source: azurerm_key_vault_secrets

Use this data source to retrieve a list of secret names from an existing Key Vault Secret.

## Example Usage

```hcl
data "azurerm_key_vault_secrets" "example" {
  key_vault_id = data.azurerm_key_vault.existing.id
}

data "azurerm_key_vault_secret" "example" {
  for_each     = toset(data.azurerm_key_vault_secrets.example.names)
  name         = each.key
  key_vault_id = data.azurerm_key_vault.existing.id
}

```

## Argument Reference

The following arguments are supported:

* `key_vault_id` - (Required) Specifies the ID of the Key Vault instance to fetch secret names from, available on the `azurerm_key_vault` Data Source / Resource.

**NOTE:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

In addition to the Argument listed above - the following Attributes are exported:

* `names` - List containing names of secrets that exist in this Key Vault.

* `secrets` - One or more `secrets` blocks as defined below.

---

A `secrets` block supports following:

* `name` - The name of secret.

* `enabled` - Whether this secret is enabled.

* `id` - The ID of this secret.

* `tags` - The tags of this secret.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Secret.
