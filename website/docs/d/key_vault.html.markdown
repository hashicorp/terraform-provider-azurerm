---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault"
sidebar_current: "docs-azurerm-data-source-key-vault-x"
description: |-
  Gets information about a Key Vault.
---

# Data Source: azurerm_key_vault

Gets information about a Key Vault.

## Example Usage

```hcl
data "azurerm_key_vault" "test" {
  name                = "mykeyvault"
  resource_group_name = "some-resource-group"
}

output "vault_uri" {
  value = "${data.azurerm_key_vault.test.vault_uri}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault.

* `resource_group_name` - (Required) The name of the Resource Group in which the Key Vault exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Vault ID.

* `vault_uri` - The URI of the vault for performing operations on keys and secrets.

* `location` - The Azure Region in which the Key Vault exists.

* `sku` - A `sku` block as described below.

* `tenant_id` - The Azure Active Directory Tenant ID used for authenticating requests to the Key Vault.

* `access_policy` - One or more `access_policy` blocks as defined below.

* `enabled_for_deployment` - Can Azure Virtual Machines retrieve certificates stored as secrets from the Key Vault?

* `enabled_for_disk_encryption` - Can Azure Disk Encryption retrieve secrets from the Key Vault?

* `enabled_for_template_deployment` - Can Azure Resource Manager retrieve secrets from the Key Vault?

* `tags` - A mapping of tags assigned to the Key Vault.

A `sku` block exports the following:

* `name` - The name of the SKU used for this Key Vault.

`access_policy` supports the following:

* `tenant_id` - The Azure Active Directory Tenant ID used to authenticate requests for this Key Vault.

* `object_id` - An Object ID of a User, Service Principal or Security Group.

* `application_id` - The Object ID of a Azure Active Directory Application.

* `certificate_permissions` - A list of certificate permissions applicable to this Access Policy.

* `key_permissions` - A list of key permissions applicable to this Access Policy.

* `secret_permissions` - A list of secret permissions applicable to this Access Policy.