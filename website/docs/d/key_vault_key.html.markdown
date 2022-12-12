---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_key"
description: |-
  Gets information about an existing Key Vault Key.

---

# Data Source: azurerm_key_vault_key

Use this data source to access information about an existing Key Vault Key.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_key_vault_key" "example" {
  name         = "secret-sauce"
  key_vault_id = data.azurerm_key_vault.existing.id
}

output "key_type" {
  value = data.azurerm_key_vault_key.example.key_type
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Key Vault Key.

* `key_vault_id` - Specifies the ID of the Key Vault instance where the Secret resides, available on the `azurerm_key_vault` Data Source / Resource.

**NOTE:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Key Vault Key.

* `curve` - The EC Curve name of this Key Vault Key.

* `e` - The RSA public exponent of this Key Vault Key.

* `key_type` - Specifies the Key Type of this Key Vault Key

* `key_size` - Specifies the Size of this Key Vault Key.

* `key_opts` - A list of JSON web key operations assigned to this Key Vault Key

* `n` - The RSA modulus of this Key Vault Key.

* `public_key_pem` - The PEM encoded public key of this Key Vault Key.

* `public_key_openssh` - The OpenSSH encoded public key of this Key Vault Key.

* `resource_id` - The (Versioned) ID for this Key Vault Key. This property points to a specific version of a Key Vault Key, as such using this won't auto-rotate values if used in other Azure Services.

* `resource_versionless_id` - The Versionless ID of the Key Vault Key. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Key is updated.

* `tags` - A mapping of tags assigned to this Key Vault Key.

* `version` - The current version of the Key Vault Key.

* `versionless_id` - The Base ID of the Key Vault Key.

* `x` - The EC X component of this Key Vault Key.

* `y` - The EC Y component of this Key Vault Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the Key Vault Key.
