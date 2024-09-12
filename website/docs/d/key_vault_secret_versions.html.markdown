---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_secret_versions"
description: |-
  Gets information about an existing Key Vault Secret Versions.
---

# Data Source: azurerm_key_vault_secret_versions

Use this data source to access information about an existing Key Vault Secret Versions.

## Example Usage

```hcl
data "azurerm_key_vault_secret_versions" "example" {
  name = "existing"
  key_vault_id = "TODO"
}

output "id" {
  value = data.azurerm_key_vault_secret_versions.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required) The ID of the TODO.

* `name` - (Required) The name of this Key Vault Secret Versions.

---

* `max_results` - (Optional) TODO. Defaults to `25`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Key Vault Secret Versions.

* `versions` - A `versions` block as defined below.

---

A `versions` block exports the following:

* `created_date` - TODO.

* `enabled` - Is the TODO enabled?

* `expiration_date` - TODO.

* `id` - TODO.

* `not_before_date` - TODO.

* `updated_date` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Secret Versions.