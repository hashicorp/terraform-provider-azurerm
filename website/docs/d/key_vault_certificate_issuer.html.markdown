---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_certificate_issuer"
description: |-
  Gets information about an existing Key Vault Certificate Issuer.
---

# Data Source: azurerm_key_vault_certificate_issuer

Use this data source to access information about an existing Key Vault Certificate Issuer.

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "mykeyvault"
  resource_group_name = "some-resource-group"
}

data "azurerm_key_vault_certificate_issuer" "example" {
  name         = "existing"
  key_vault_id = data.azurerm_key_vault.example.id
}

output "id" {
  value = data.azurerm_key_vault_certificate_issuer.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required) The ID of the Key Vault in which to locate the Certificate Issuer.

* `name` - (Required) The name of the Key Vault Certificate Issuer.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Key Vault Certificate Issuer.

* `account_id` - The account number with the third-party Certificate Issuer.

* `org_id` - The organization ID with the third-party Certificate Issuer.

* `admin` - A list of `admin` blocks as defined below.

* `provider_name` - The name of the third-party Certificate Issuer.

---

An `admin` block exports the following:

* `email_address` - E-mail address of the admin.

* `first_name` - First name of the admin.

* `last_name` - Last name of the admin.

* `phone` - Phone number of the admin.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate Issuer.
