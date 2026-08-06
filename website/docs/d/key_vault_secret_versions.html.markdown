---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_secret_versions"
description: |-
  Get a list of versions for an existing Key Vault Secret.
---

# Data Source: azurerm_key_vault_secret_versions

Use this data source to access information about an existing Key Vault Secret's versions. The secret version values is not included. The `key_vault_secret` data source can be used to retrieve the value of a given secret version using it's `id`.

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "mykeyvault"
  resource_group_name = "some-resource-group"
}

data "azurerm_key_vault_secret_versions" "example" {
  name         = "mysecret"
  key_vault_id = data.azurerm_key_vault.example.id
}

output "versions" {
  value = data.azurerm_key_vault_secret_versions.example.versions
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required) The ID of the Key Vault containing the secret.

* `name` - (Required) The name of the Key Vault Secret to retrieve versions from.

---

* `max_results` - (Optional) Maximum number of versions to retrieve. Defaults to `25`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Key Vault Secret.

* `versions` - A `versions` list as defined below. The list of versions are sorted by `created_date` descending, meaning the most recently created secret version will be first in the list.

---

The `versions` entries in the list export the following:

* `created_date` - The date and time when the Key Vault Secret version was created.

* `enabled` - Is the version enabled? Returns a `bool` value.

* `expiration_date` - The date and time at which the Key Vault Secret version expires and is no longer valid.

* `id` - The Key Vault Secret version ID.

* `not_before_date` - The earliest date and time at which the Key Vault Secret version can be used.

* `updated_date` - The date and time when the Key Vault Secret version was last updated.

* `uri` - The full URI of the secret version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Secret Versions.
