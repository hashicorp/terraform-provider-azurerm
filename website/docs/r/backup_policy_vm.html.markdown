---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_policy_vm"
description: |-
  Manages an Azure Backup VM Backup Policy.
---

# azurerm_backup_policy_vm

Manages an Azure Backup VM Backup Policy.

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

resource "azurerm_backup_policy_vm" "example" {
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
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday"]
    weeks    = ["Last"]
    months   = ["January"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Backup Policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the policy. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) Specifies the name of the Recovery Services Vault to use. Changing this forces a new resource to be created.

* `backup` - (Required) Configures the Policy backup frequency, times & days as documented in the `backup` block below.

* `policy_type` - (Optional) Type of the Backup Policy. Possible values are `V1` and `V2` where `V2` stands for the Enhanced Policy. Defaults to `V1`. Changing this forces a new resource to be created.

* `timezone` - (Optional) Specifies the timezone. [the possible values are defined here](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/). Defaults to `UTC`

* `instant_restore_retention_days` - (Optional) Specifies the instant restore retention range in days. Possible values are between `1` and `5` when `policy_type` is `V1`, and `1` to `30` when `policy_type` is `V2`.

~> **NOTE:** `instant_restore_retention_days` **must** be set to `5` if the backup frequency is set to `Weekly`.

* `instant_restore_resource_group` - (Optional) Specifies the instant restore resource group name as documented in the `instant_restore_resource_group` block below.

* `retention_daily` - (Optional) Configures the policy daily retention as documented in the `retention_daily` block below. Required when backup frequency is `Daily`.

* `retention_weekly` - (Optional) Configures the policy weekly retention as documented in the `retention_weekly` block below. Required when backup frequency is `Weekly`.

* `retention_monthly` - (Optional) Configures the policy monthly retention as documented in the `retention_monthly` block below.

* `retention_yearly` - (Optional) Configures the policy yearly retention as documented in the `retention_yearly` block below.

---

The `backup` block supports:

* `frequency` - (Required) Sets the backup frequency. Possible values are `Hourly`, `Daily` and `Weekly`.

* `time` - (Required) The time of day to perform the backup in 24hour format.

* `hour_interval` - (Optional) Interval in hour at which backup is triggered. Possible values are `4`, `6`, `8` and `12`. This is used when `frequency` is `Hourly`.

* `hour_duration` - (Optional) Duration of the backup window in hours. Possible values are between `4` and `24` This is used when `frequency` is `Hourly`.

~> **NOTE:** `hour_duration` must be multiplier of `hour_interval`

* `weekdays` - (Optional) The days of the week to perform backups on. Must be one of `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`. This is used when `frequency` is `Weekly`.

---
The `instant_restore_resource_group` block supports:

* `prefix` - (Required) The prefix for the `instant_restore_resource_group` name.

* `suffix` - (Optional) The suffix for the `instant_restore_resource_group` name.

---

The `retention_daily` block supports:

* `count` - (Required) The number of daily backups to keep. Must be between `7` and `9999`.

~> **Note:** Azure previously allows this field to be set to a minimum of 1 (day) - but for new resources/to update this value on existing Backup Policies - this value must now be at least 7 (days).

---

The `retention_weekly` block supports:

* `count` - (Required) The number of weekly backups to keep. Must be between `1` and `9999`

* `weekdays` - (Required) The weekday backups to retain. Must be one of `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

---

The `retention_monthly` block supports:

* `count` - (Required) The number of monthly backups to keep. Must be between `1` and `9999`

* `weekdays` - (Optional) The weekday backups to retain . Must be one of `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

* `weeks` - (Optional) The weeks of the month to retain backups of. Must be one of `First`, `Second`, `Third`, `Fourth`, `Last`.

* `days` - (Optional) The days of the month to retain backups of. Must be between `1` and `31`.

* `include_last_days` - (Optional) Including the last day of the month, default to `false`.

-> **NOTE:**: Either `weekdays` and `weeks` or `days` and `include_last_days` must be specified.

---

The `retention_yearly` block supports:

* `count` - (Required) The number of yearly backups to keep. Must be between `1` and `9999`

* `months` - (Required) The months of the year to retain backups of. Must be one of `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November` and `December`.

* `weekdays` - (Optional) The weekday backups to retain . Must be one of `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`.

* `weeks` - (Optional) The weeks of the month to retain backups of. Must be one of `First`, `Second`, `Third`, `Fourth`, `Last`.

* `days` - (Optional) The days of the month to retain backups of. Must be between `1` and `31`.

* `include_last_days` - (Optional) Including the last day of the month, default to `false`.

-> **NOTE:**: Either `weekdays` and `weeks` or `days` and `include_last_days` must be specified.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the VM Backup Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the VM Backup Policy.
* `update` - (Defaults to 30 minutes) Used when updating the VM Backup Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the VM Backup Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the VM Backup Policy.

## Import

VM Backup Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_policy_vm.policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupPolicies/policy1
```
