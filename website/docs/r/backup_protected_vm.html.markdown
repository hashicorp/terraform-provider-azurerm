---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_protected_vm"
description: |-
  Manages an Azure Backup Protected Virtual Machine.
---

# azurerm_backup_protected_vm

Manages an Azure Backup Protected Virtual Machine.

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

  backup {
    frequency = "Daily"
    time      = "23:00"
  }
  retention_daily {
    count = 10
  }
}

data "azurerm_virtual_machine" "example" {
  name                = "example-vm"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_backup_protected_vm" "vm1" {
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.example.name
  source_vm_id        = data.azurerm_virtual_machine.example.id
  backup_policy_id    = azurerm_backup_policy_vm.example.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) Specifies the name of the Resource Group **associated with** the Recovery Services Vault to use. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) Specifies the name of the Recovery Services Vault to use. Changing this forces a new resource to be created.

* `source_vm_id` - (Optional) Specifies the ID of the virtual machine to back up. Changing this forces a new resource to be created.

~> **Note:** After creation, the `source_vm_id` property can be removed without forcing a new resource to be created; however, setting it to a different ID will create a new resource. This allows the source virtual machine to be deleted without having to remove the backup.

* `backup_policy_id` - (Optional) Specifies the ID of the backup policy to use.

~> **Note:** `backup_policy_id` is required during initial creation of this resource.

~> **Note:** When `protection_state` is set to `BackupsSuspended` or `ProtectionStopped`, the Azure API may not return `backup_policy_id`. To avoid a perpetual diff, use Terraform's [ignore_changes](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#ignore_changes) argument.

* `exclude_disk_luns` - (Optional) A list of Disks' Logical Unit Numbers (LUN) to be excluded for VM Protection.

* `include_disk_luns` - (Optional) A list of Disks' Logical Unit Numbers (LUN) to be included for VM Protection.

* `protection_state` - (Optional) Specifies Protection state of the backup. Possible values are `Protected`, `BackupsSuspended`, and `ProtectionStopped`.

~> **Note:** `protection_state` cannot be set to `BackupsSuspended` unless the `azurerm_recovery_services_vault` has `immutability` set to `Unlocked` or `Locked`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Protected Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Backup Protected Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Protected Virtual Machine.
* `update` - (Defaults to 2 hours) Used when updating the Backup Protected Virtual Machine.
* `delete` - (Defaults to 80 minutes) Used when deleting the Backup Protected Virtual Machine.

## Import

Backup Protected Virtual Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_protected_vm.item1 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1"
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.RecoveryServices`: 2024-01-01
