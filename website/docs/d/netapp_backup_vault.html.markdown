---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: netapp_backup_vault"
description: |-
  Gets information about an existing NetApp Backup Vault
---

# Data Source: netapp_backup_vault

Use this data source to access information about an existing NetApp Backup Vault.

## NetApp Backup Vault Usage

```hcl
data "azurerm_netapp_backup_vault" "example" {
  resource_group_name = "example-resource-group"
  account_name        = "example-netappaccount"
  name                = "example-backupvault"
}

output "backup_vault_id" {
  value = data.azurerm_netapp_backup_vault.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the NetApp Backup Vault.

* `resource_group_name` - The name of the resource group where the NetApp Backup Vault exists.

* `account_name` - The name of the NetApp Account in which the NetApp Vault exists.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Backup Vault.

## Import

NetApp Backup Vault can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_backup_vault.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/backupPolicies/backupvault1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.NetApp`: 2025-01-01
