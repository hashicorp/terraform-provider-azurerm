---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_secrets"
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
  for_each = data.azurerm_key_vault_secrets.example.names
  name     = each.key
}

```

## Argument Reference

The following arguments are supported:

* `key_vault_id` - Specifies the ID of the Key Vault instance to fetch secret names from, available on the `azurerm_key_vault` Data Source / Resource.

**NOTE:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

The following attributes are exported:

* `names` - List containing names of secrets that exist in this Key Vault.
* `key_vault_id` - The Key Vault ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Secret.
