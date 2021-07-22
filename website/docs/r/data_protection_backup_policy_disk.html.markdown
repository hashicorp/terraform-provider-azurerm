---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_disk"
description: |-
  Manages a Backup Policy Disk.
---

# azurerm_data_protection_backup_policy_disk

Manages a Backup Policy Disk.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-backup-vault"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}

resource "azurerm_data_protection_backup_policy_disk" "example" {
  name     = "example-backup-policy"
  vault_id = azurerm_data_protection_backup_vault.example.id

  backup_repeating_time_intervals = ["R/2021-05-19T06:33:16+00:00/PT4H"]
  default_retention_duration      = "P7D"

  retention_rule {
    name     = "Daily"
    duration = "P7D"
    priority = 25
    criteria {
      absolute_criteria = "FirstOfDay"
    }
  }

  retention_rule {
    name     = "Weekly"
    duration = "P7D"
    priority = 20
    criteria {
      absolute_criteria = "FirstOfWeek"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:
* `name` - (Required) The name which should be used for this Backup Policy Disk. Changing this forces a new Backup Policy Disk to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Backup Policy Disk should exist. Changing this forces a new Backup Policy Disk to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the Backup Policy Disk should exist. Changing this forces a new Backup Policy Disk to be created.

* `backup_repeating_time_intervals` - (Required) Specifies a list of repeating time interval. It should follow `ISO 8601` repeating time interval . Changing this forces a new Backup Policy Disk to be created.

* `default_retention_duration` - (Required) The duration of default retention rule. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Disk to be created.

---

* `retention_rule` - (Optional) One or more `retention_rule` blocks as defined below. Changing this forces a new Backup Policy Disk to be created.

---

A `retention_rule` block supports the following:

* `name` - (Required) The name which should be used for this retention rule. Changing this forces a new Backup Policy Disk to be created.

* `duration` - (Required) Duration of deletion after given timespan. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Disk to be created.

* `criteria` - (Required) A `criteria` block as defined below. Changing this forces a new Backup Policy Disk to be created.

* `priority` - (Required) Retention Tag priority. Changing this forces a new Backup Policy Disk to be created.

---

A `criteria` block supports the following:

* `absolute_criteria` - (Optional) Possible values are `FirstOfDay` and `FirstOfWeek`. Changing this forces a new Backup Policy Disk to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Backup Policy Disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Policy Disk.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Policy Disk.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Policy Disk.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Policy Disk.

## Import

Backup Policy Disks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_disk.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupPolicies/backupPolicy1
```
