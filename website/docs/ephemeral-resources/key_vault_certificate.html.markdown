---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate"
description: |-
  Gets information about an existing Key Vault Certificate.
---

# Ephemeral: azurerm_key_vault_certificate

~> **Note:** Ephemeral Resources are supported in Terraform 1.10 and later.

Use this to access information about an existing Key Vault Certificate.

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "examplekv"
  resource_group_name = "some-resource-group"
}

ephemeral "azurerm_key_vault_certificate" "example" {
  name         = "secret-sauce"
  key_vault_id = data.azurerm_key_vault.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Certificate.

* `key_vault_id` - (Required) Specifies the ID of the Key Vault instance where the Certificate resides, available on the `azurerm_key_vault` Data Source / Resource.

* `version` - (Optional) Specifies the version of the Key Vault Certificate. Defaults to the current version of the Key Vault Certificate.

~> **Note:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

The following attributes are exported:

* `hex` - The raw Key Vault Certificate data represented as a hexadecimal string.

* `pem` - The Key Vault Certificate in PEM format.

* `key` - The Key Vault Certificate Key.

* `expiration_date` - The date and time at which the Key Vault Certificate expires and is no longer valid.

* `not_before_date` - The earliest date at which the Key Vault Certificate can be used.
