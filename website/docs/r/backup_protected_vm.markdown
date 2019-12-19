---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_protected_vm"
sidebar_current: "docs-azurerm-backup-protected-vm"
description: |-
  Manages an Azure Backup Protected VM.
---

# azurerm_backup_protected_vm

Manages Azure Backup for an Azure VM

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-recovery_vault"
  location = "West US"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "tfex-recovery-vault"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Standard"
}

resource "azurerm_backup_policy_vm" "example" {
  name                = "tfex-recovery-vault-policy"
  resource_group_name = "${azurerm_resource_group.example.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.example.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }
}

resource "azurerm_backup_protected_vm" "vm1" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.example.name}"
  source_vm_id        = "${azurerm_virtual_machine.example.id}"
  backup_policy_id    = "${azurerm_backup_policy_vm.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which to create the Recovery Services Protected VM. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) Specifies the name of the Recovery Services Vault to use. Changing this forces a new resource to be created.

* `source_vm_id` - (Required) Specifies the ID of the VM to backup. Changing this forces a new resource to be created.

* `backup_policy_id` - (Required) Specifies the id of the backup policy to use.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Recovery Services Vault.

## Import

Recovery Services Protected VMs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_protected_vm.item1 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1"
```

Note the ID requires quoting as there are semicolons
 