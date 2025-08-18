---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_policy_file_share"
description: |-
  Manages an Azure File Share Backup Policy.
---

# azurerm_backup_policy_file_share

Manages an Azure File Share Backup Policy within a Recovery Services vault.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-recovery_vault"
  location = "West Europe"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "tfex-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_backup_policy_file_share" "policy" {
  name                = "tfex-recovery-vault-policy"
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.example.name

  timezone = "UTC"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }

  retention_weekly {
    count    = 7
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 7
    weekdays = ["Sunday"]
    weeks    = ["Last"]
    months   = ["January"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the policy. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) Specifies the name of the Recovery Services Vault to use. Changing this forces a new resource to be created.

* `backup` - (Required) Configures the Policy backup frequency and times as documented in the `backup` block below.

* `timezone` - (Optional) Specifies the timezone. [the possible values are defined here](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/). Defaults to `UTC`

-> **Note:** The maximum number of snapshots that Azure Files can retain is 200. If your combined snapshot count exceeds 200 based on your retention policies, it will result in an error. See [this](https://docs.microsoft.com/azure/backup/backup-azure-files-faq#what-is-the-maximum-retention-i-can-configure-for-backups) article for more information.

* `retention_daily` - (Required) Configures the policy daily retention as documented in the `retention_daily` block below.

* `retention_weekly` - (Optional) Configures the policy weekly retention as documented in the `retention_weekly` block below.

* `retention_monthly` - (Optional) Configures the policy monthly retention as documented in the `retention_monthly` block below.

* `retention_yearly` - (Optional) Configures the policy yearly retention as documented in the `retention_yearly` block below.

---

The `backup` block supports:

* `frequency` - (Required) Sets the backup frequency. Possible values are `Daily` and `Hourly`. 

-> **Note:** This argument is made available for consistency with VM backup policies and to allow for potential future support of weekly backups

* `time` - (Optional) The time of day to perform the backup in 24-hour format. Times must be either on the hour or half hour (e.g. 12:00, 12:30, 13:00, etc.)

-> **Note:** `time` is required when `frequency` is set to `Daily`.

* `hourly` - (Optional) A `hourly` block defined as below. This is required when `frequency` is set to `Hourly`.

---

The `hourly` block supports:

* `interval` - (Required) Specifies the interval at which backup needs to be triggered. Possible values are `4`, `6`, `8` and `12`.

* `start_time` - (Required) Specifies the start time of the hourly backup. The time format should be in 24-hour format. Times must be either on the hour or half hour (e.g. 12:00, 12:30, 13:00, etc.).

* `window_duration` - (Required) Species the duration of the backup window in hours. Details could be found [here](https://learn.microsoft.com/en-us/azure/backup/backup-azure-files-faq#what-does-the-duration-attribute-in-azure-files-backup-policy-signify-).

---

The `retention_daily` block supports:

* `count` - (Required) The number of daily backups to keep. Must be between `1` and `200` (inclusive)

---

The `retention_weekly` block supports:

* `count` - (Required) The number of daily backups to keep. Must be between `1` and `200` (inclusive)

* `weekdays` - (Required) The weekday backups to retain. Must be one of `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

---

The `retention_monthly` block supports:

* `count` - (Required) The number of monthly backups to keep. Must be between `1` and `120`

* `weekdays` - (Optional) The weekday backups to retain . Must be one of `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

* `weeks` - (Optional) The weeks of the month to retain backups of. Must be one of `First`, `Second`, `Third`, `Fourth`, `Last`.

* `days` - (Optional) The days of the month to retain backups of. Must be between `1` and `31`.

* `include_last_days` - (Optional) Including the last day of the month, default to `false`.

-> **Note:** Either `weekdays` and `weeks` or `days` and `include_last_days` must be specified.

---

The `retention_yearly` block supports:

* `count` - (Required) The number of yearly backups to keep. Must be between `1` and `10`

* `months` - (Required) The months of the year to retain backups of. Must be one of `January`, `February`, `March`, `April`, `May`, `June`, `July`, `Augest`, `September`, `October`, `November` and `December`.

* `weeks` - (Optional) The weeks of the month to retain backups of. Must be one of `First`, `Second`, `Third`, `Fourth`, `Last`.

* `weekdays` - (Optional) The weekday backups to retain . Must be one of `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

* `days` - (Optional) The days of the month to retain backups of. Must be between `1` and `31`.

* `include_last_days` - (Optional) Including the last day of the month, default to `false`.

-> **Note:** Either `weekdays` and `weeks` or `days` and `include_last_days` must be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure File Share Backup Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the File Share Backup Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the File Share Backup Policy.
* `update` - (Defaults to 30 minutes) Used when updating the File Share Backup Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the File Share Backup Policy.

## Import

Azure File Share Backup Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_policy_file_share.policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupPolicies/policy1
```
