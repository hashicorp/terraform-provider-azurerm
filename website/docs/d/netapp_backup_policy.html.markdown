---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: netapp_backup_policy"
description: |-
  Gets information about an existing NetApp Backup Policy
---

# Data Source: netapp_backup_policy

Use this data source to access information about an existing NetApp Backup Vault.

## NetApp Backup Policy Usage

```hcl
data "azurerm_netapp_backup_policy" "example" {
  resource_group_name = "example-resource-group"
  account_name        = "example-netappaccount"
  name                = "example-backuppolicy"
}

output "backup_policy_id" {
  value = data.azurerm_netapp_backup_policy.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the NetApp Backup Policy.

* `resource_group_name` - The name of the resource group where the NetApp Backup Policy exists.

* `account_name` - The name of the NetApp Account in which the NetApp Policy exists.

## Attributes Reference

The following attributes are exported:

* `location` - NetApp Backup Policy location.

* `account_name` - The name of the NetApp account in which the NetApp Policy exists.

* `daily_backups_to_keep` - The number of daily backups to keep.

* `weekly_backups_to_keep` - The number of weekly backups to keep.

* `monthly_backups_to_keep` - The number of monthly backups to keep.

* `enabled` - Whether the Backup Policy is enabled.

* `tags` - List of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Backup Policy.

## Import

NetApp Backup Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_backup_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/backupPolicies/backuppolicy1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.NetApp`: 2025-01-01
