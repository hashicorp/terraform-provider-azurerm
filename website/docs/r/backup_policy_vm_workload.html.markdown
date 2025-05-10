---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_policy_vm_workload"
description: |-
  Manages an Azure VM Workload Backup Policy.
---

# azurerm_backup_policy_vm_workload

Manages an Azure VM Workload Backup Policy within a Recovery Services vault.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-bpvmw"
  location = "West Europe"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "example-rsv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm_workload" "example" {
  name                = "example-bpvmw"
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.example.name

  workload_type = "SQLDataBase"

  settings {
    time_zone           = "UTC"
    compression_enabled = false
  }

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Daily"
      time      = "15:00"
    }

    retention_daily {
      count = 8
    }
  }

  protection_policy {
    policy_type = "Log"

    backup {
      frequency_in_minutes = 15
    }

    simple_retention {
      count = 8
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VM Workload Backup Policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the VM Workload Backup Policy. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) The name of the Recovery Services Vault to use. Changing this forces a new resource to be created.

* `protection_policy` - (Required) One or more `protection_policy` blocks as defined below.

* `settings` - (Required) A `settings` block as defined below.

* `workload_type` - (Required) The VM Workload type for the Backup Policy. Possible values are `SQLDataBase` and `SAPHanaDatabase`. Changing this forces a new resource to be created.

---

The `protection_policy` block supports the following:

* `policy_type` - (Required) The type of the VM Workload Backup Policy. Possible values are `Differential`, `Full`, `Incremental` and `Log`.

* `backup` - (Required) A `backup` block as defined below.

* `retention_daily` - (Optional) A `retention_daily` block as defined below.

* `retention_weekly` - (Optional) A `retention_weekly` block as defined below.

* `retention_monthly` - (Optional) A `retention_monthly` block as defined below.

* `retention_yearly` - (Optional) A `retention_yearly` block as defined below.

* `simple_retention` - (Optional) A `simple_retention` block as defined below.

---

The `simple_retention` block supports the following:

* `count` - (Required) The count that is used to count retention duration with duration type `Days`. Possible values are between `7` and `35`.

---

The `settings` block supports the following:

* `time_zone` - (Required) The timezone for the VM Workload Backup Policy. [The possible values are defined here](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/).

* `compression_enabled` - (Optional) The compression setting for the VM Workload Backup Policy. Defaults to `false`.

---

The `backup` block supports the following:

* `frequency` - (Optional) The backup frequency for the VM Workload Backup Policy. Possible values are `Daily` and `Weekly`.

* `frequency_in_minutes` - (Optional) The backup frequency in minutes for the VM Workload Backup Policy. Possible values are `15`, `30`, `60`, `120`, `240`, `480`, `720` and `1440`.

* `time` - (Optional) The time of day to perform the backup in 24hour format.

* `weekdays` - (Optional) The days of the week to perform backups on. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`. This is used when `frequency` is `Weekly`.

---

The `retention_daily` block supports the following:

* `count` - (Required) The number of daily backups to keep. Possible values are between `7` and `9999`.

---

The `retention_weekly` block supports the following:

* `count` - (Required) The number of weekly backups to keep. Possible values are between `1` and `5163`.

* `weekdays` - (Required) The weekday backups to retain. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

---

The `retention_monthly` block supports the following:

* `count` - (Required) The number of monthly backups to keep. Must be between `1` and `1188`.

* `format_type` - (Required) The retention schedule format type for monthly retention policy. Possible values are `Daily` and `Weekly`.

* `monthdays` - (Optional) The monthday backups to retain. Possible values are between `0` and `28`.

* `weekdays` - (Optional) The weekday backups to retain. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

* `weeks` - (Optional) The weeks of the month to retain backups of. Possible values are `First`, `Second`, `Third`, `Fourth` and `Last`.

---

The `retention_yearly` block supports the following:

* `count` - (Required) The number of yearly backups to keep. Possible values are between `1` and `99`

* `format_type` - (Required) The retention schedule format type for yearly retention policy. Possible values are `Daily` and `Weekly`.

* `months` - (Required) The months of the year to retain backups of. Possible values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November` and `December`.

* `monthdays` - (Optional) The monthday backups to retain. Possible values are between `0` and `28`.

* `weekdays` - (Optional) The weekday backups to retain. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

* `weeks` - (Optional) The weeks of the month to retain backups of. Possible values are `First`, `Second`, `Third`, `Fourth`, `Last`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure VM Workload Backup Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the VM Workload Backup Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the VM Workload Backup Policy.
* `update` - (Defaults to 30 minutes) Used when updating the VM Workload Backup Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the VM Workload Backup Policy.

## Import

Azure VM Workload Backup Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_policy_vm_workload.policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupPolicies/policy1
```
