---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_protection_backup_vault"
description: |-
  Manages a Backup Vault.
---

# Data Source: azurerm_data_protection_backup_vault

Use this data source to access information about an existing Backup Vault.

## Example Usage

```hcl
data "azurerm_data_protection_backup_vault" "example" {
  name                = "existing-backup-vault"
  resource_group_name = "existing-resource-group"
}

output "azurerm_data_protection_backup_vault_id" {
  value = data.azurerm_vpn_gateway.example.id
}

output "azurerm_data_protection_backup_vault_principal_id" {
  value = data.azurerm_data_protection_backup_vault.example.identity.0.principal_id
}
```

## Arguments Reference

* `name` - (Required) Specifies the name of the Backup Vault.

* `resource_group_name` - (Required) The name of the Resource Group where the Backup Vault exists.

## Attributes Reference

* `id` - The ID of the Backup Vault.

* `location` -  The Azure Region where the Backup Vault exists.

* `datastore_type` - Specifies the type of the data store.

* `redundancy` -  Specifies the backup storage redundancy.

* `identity` -  A `identity` block as defined below.

* `tags` -  A mapping of tags which are assigned to the Backup Vault.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Backup Vault.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Backup Vault.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Backup Vault.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Vault.
