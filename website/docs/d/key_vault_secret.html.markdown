---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_secret"
description: |-
  Gets information about an existing Key Vault Secret.
---

# Data Source: azurerm_key_vault_secret

Use this data source to access information about an existing Key Vault Secret.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_key_vault_secret" "example" {
  name         = "secret-sauce"
  key_vault_id = data.azurerm_key_vault.existing.id
}

output "secret_value" {
  value     = data.azurerm_key_vault_secret.example.value
  sensitive = true
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required)  Specifies the ID of the Key Vault instance to fetch secret names from, available on the `azurerm_key_vault` Data Source / Resource.

* `name` - (Required) Specifies the name of the Key Vault Secret.

* `version` - (Optional) Specifies the version of the Key Vault Secret. Defaults to the current version of the Key Vault Secret.

**NOTE:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Key Vault Secret ID.

* `content_type` - The content type for the Key Vault Secret.

* `resource_id` - The (Versioned) ID for this Key Vault Secret. This property points to a specific version of a Key Vault Secret, as such using this won't auto-rotate values if used in other Azure Services.

* `resource_versionless_id` - The Versionless ID of the Key Vault Secret. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Secret is updated.

* `tags` - Any tags assigned to this resource.

* `value` - The value of the Key Vault Secret.

* `versionless_id` - The Versionless ID of the Key Vault Secret. This can be used to always get latest secret value, and enable fetching automatically rotating secrets.

* `not_before_date` - The earliest date at which the Key Vault Secret can be used.

* `expiration_date` - The date and time at which the Key Vault Secret expires and is no longer valid.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the Key Vault Secret.
