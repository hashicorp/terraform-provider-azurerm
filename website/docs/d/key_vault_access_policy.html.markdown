---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_access_policy"
description: |-
  Get information about the templated Access Policies for Key Vault.
---

# Data Source: azurerm_key_vault_access_policy

Use this data source to access information about the permissions from the Management Key Vault Templates.

## Example Usage

```hcl
data "azurerm_key_vault_access_policy" "contributor" {
  name = "Key Management"
}

output "access_policy_key_permissions" {
  value = data.azurerm_key_vault_access_policy.contributor.key_permissions
}
```

## Argument Reference

* `name` - Specifies the name of the Management Template. Possible values are: `Key Management`,
`Secret Management`, `Certificate Management`, `Key & Secret Management`, `Key & Certificate Management`,
`Secret & Certificate Management`,  `Key, Secret, & Certificate Management`


## Attributes Reference

The following attributes are exported:

* `id` - the ID of the Key Vault Access Policy

* `key_permissions` - the key permissions for the access policy

* `secret_permissions` - the secret permissions for the access policy

* `certificate_permissions` - the certificate permissions for the access policy

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Access Policy.
