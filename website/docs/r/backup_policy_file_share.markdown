---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_policy_file_share"
description: |-
  Manages an Azure File Share Backup Policy.
---

# azurerm_backup_policy_file_share

Manages an Azure File Share Backup Policy within a Recovery Services vault.

-> **NOTE:** Azure Backup for Azure File Shares is currently in public preview. During the preview, the service is subject to additional limitations and unsupported backup scenarios. [Read More](https://docs.microsoft.com/en-us/azure/backup/backup-azure-files#limitations-for-azure-file-share-backup-during-preview)

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "tfex-recovery_vault"
  location = "West US"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "tfex-recovery-vault"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  sku                 = "Standard"
}

resource "azurerm_backup_policy_file_share" "policy" {
  name                = "tfex-recovery-vault-policy"
  resource_group_name = azurerm_resource_group.rg.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name

  timezone = "UTC"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the policy. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) Specifies the name of the Recovery Services Vault to use. Changing this forces a new resource to be created.

* `backup` - (Required) Configures the Policy backup frequency and times as documented in the `backup` block below.

* `timezone` - (Optional) Specifies the timezone. [the possible values are defined here](http://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/). Defaults to `UTC`

* `retention_daily` - (Required) Configures the policy daily retention as documented in the `retention_daily` block below.

-> **NOTE:** During the public preview, only daily retentions are supported. This argument is made available in this format for consistency with VM backup policies and to allow for potential future support of additional retention policies

---

The `backup` block supports:

* `frequency` - (Required) Sets the backup frequency. Currently, only `Daily` is supported

-> **NOTE:** During the public preview, only daily backups are supported. This argument is made available for consistency with VM backup policies and to allow for potential future support of weekly backups

* `times` - (Required) The time of day to perform the backup in 24-hour format. Times must be either on the hour or half hour (e.g. 12:00, 12:30, 13:00, etc.)

---

The `retention_daily` block supports:

* `count` - (Required) The number of daily backups to keep. Must be between `1` and `180` (inclusive)

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Azure File Share Backup Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the File Share Backup Policy.
* `update` - (Defaults to 30 minutes) Used when updating the File Share Backup Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the File Share Backup Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the File Share Backup Policy.

## Import

Azure File Share Backup Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_policy_file_share.policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupPolicies/policy1
```
