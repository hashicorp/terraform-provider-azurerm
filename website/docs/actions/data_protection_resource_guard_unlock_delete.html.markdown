---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_resource_guard_unlock_delete"
description: |-
  Unlocks deletion of a Recovery Services backup item protected by Resource Guard.
---

# Action: azurerm_data_protection_resource_guard_unlock_delete

Unlocks deletion of a backup item protected by Resource Guard. Use this action before deleting the item to keep the vault's Resource Guard association in place. If deletion is not protected by Resource Guard, no unlock is needed and the action completes successfully.

~> **Note:** The identity running Terraform must have permission to read the backup item and its Resource Guard association, and to unlock deletion. See [Multi-user authorization for Recovery Services vaults](https://learn.microsoft.com/azure/backup/multi-user-authorization-concept) for the required permissions.

## Example Usage

This example unlocks a VM backup immediately before Terraform deletes it. `caller.id` supplies the ID of the backup item being deleted. Terraform 1.16 or later is required for `before_destroy`, `caller`, and `on_failure` features.

The provider feature `vm_backup_stop_protection_and_retain_data_on_destroy` must be `false` for deletion. This action does not unlock operations that retain backup data.

```terraform
resource "azurerm_recovery_services_vault_resource_guard_association" "example" {
  vault_id          = azurerm_recovery_services_vault.example.id
  resource_guard_id = azurerm_data_protection_resource_guard.example.id
}

action "azurerm_data_protection_resource_guard_unlock_delete" "example" {
  config {
    protected_item_id = caller.id
  }
}

resource "azurerm_backup_protected_vm" "example" {
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.example.name
  source_vm_id        = azurerm_linux_virtual_machine.example.id
  backup_policy_id    = azurerm_backup_policy_vm.example.id

  lifecycle {
    action_trigger {
      events     = [before_destroy]
      actions    = [action.azurerm_data_protection_resource_guard_unlock_delete.example]
      on_failure = halt
    }
  }

  depends_on = [azurerm_recovery_services_vault_resource_guard_association.example]
}
```

## Argument Reference

This action supports the following arguments:

* `protected_item_id` - (Required) The ID of the Recovery Services protected item to unlock for deletion.
