---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_backup_policy"
description: |-
  Manages a NetApp Backup Policy.
---

# azurerm_netapp_backup_policy

Manages a NetApp Backup Policy.

## NetApp Backup Policy Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_backup_policy" "example" {
  name                = "example-netappbackuppolicy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  account_name        = azurerm_netapp_account.example.name
  enabled             = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Backup Policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Backup Policy should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account in which the NetApp Policy should be created under. Changing this forces a new resource to be created.

* `daily_backups_to_keep` - (Optional) Provides the number of daily backups to keep, defaults to `2` which is the minimum, maximum is 1019.

* `weekly_backups_to_keep` - (Optional) Provides the number of weekly backups to keep, defaults to `1` which is the minimum, maximum is 1019.

* `monthly_backups_to_keep` - (Optional) Provides the number of monthly backups to keep, defaults to `1` which is the minimum, maximum is 1019.

~> **Note:** Currently, the combined (daily + weekly + monthy) retention counts cannot exceed 1019.

* `enabled` - (Optional) Whether the Backup Policy is enabled. Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the NetApp Backup Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Backup Policy.
* `update` - (Defaults to 2 hours) Used when updating the NetApp Backup Policy.
* `delete` - (Defaults to 2 hours) Used when deleting the NetApp Backup Policy.

## Import

NetApp Backup Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_backup_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/backupPolicies/backuppolicy1
```
