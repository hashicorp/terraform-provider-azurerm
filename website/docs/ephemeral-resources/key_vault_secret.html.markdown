---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_secret"
description: |-
  Gets information about an existing Key Vault Secret.
---

# Ephemeral: azurerm_key_vault_secret

~> **Note:** Ephemeral Resources are supported in Terraform 1.10 and later.

Use this to access information about an existing Key Vault Secret.

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "examplekv"
  resource_group_name = "some-resource-group"
}

ephemeral "azurerm_key_vault_secret" "example" {
  name         = "secret-sauce"
  key_vault_id = data.azurerm_key_vault.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Secret.

* `key_vault_id` - (Required) Specifies the ID of the Key Vault instance where the Secret resides, available on the `azurerm_key_vault` Data Source / Resource.

* `version` - (Optional) Specifies the version of the Key Vault Secret. Defaults to the current version of the Key Vault Secret.

~> **Note:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

The following attributes are exported:

* `expiration_date` - The date and time at which the Key Vault Secret expires and is no longer valid.

* `not_before_date` - The earliest date at which the Key Vault Secret can be used.

* `value` - The Key Vault Secret value.
