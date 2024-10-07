---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificates"
description: |-
  Gets a list of certificate names from an existing Key Vault.
---

# Data Source: azurerm_key_vault_certificates

Use this data source to retrieve a list of certificate names from an existing Key Vault.

## Example Usage

```hcl
data "azurerm_key_vault_certificates" "example" {
  key_vault_id = data.azurerm_key_vault.existing.id
}

data "azurerm_key_vault_certificate" "example" {
  for_each     = toset(data.azurerm_key_vault_certificates.example.names)
  name         = each.key
  key_vault_id = data.azurerm_key_vault.existing.id
}

```

## Argument Reference

The following arguments are supported:

* `key_vault_id` - Specifies the ID of the Key Vault instance to fetch certificate names from, available on the `azurerm_key_vault` Data Source / Resource.

**NOTE:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

* `include_pending` - Specifies whether to include certificates which are not completely provisioned. Defaults to true.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `names` - List containing names of certificates that exist in this Key Vault.

* `key_vault_id` - The Key Vault ID.
 
* `certificates` - One or more `certificates` blocks as defined below.

---

A `certificates` block supports following:

* `name` - The name of certificate.

* `enabled` - Whether this certificate is enabled.

* `id` - The ID of this certificate.

* `tags` - The tags of this certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificates.
