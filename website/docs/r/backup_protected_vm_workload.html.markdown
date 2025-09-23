---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_protected_vm_workload"
description: |-
    Manages a backup protected VM workload item (SAP ASE Database) in an Azure Recovery Vault
---

# azurerm_backup_protected_vm_workload

Manages a backup protected VM workload item (SAP ASE Database) in an Azure Recovery Vault. The VM workload must be registered with a backup container before it can be protected.

~> **Note:** This resource currently supports only SAP ASE Database workloads. The Virtual Machine must have the appropriate SAP workload installed and configured before registering the container and protecting the workload.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-recovery-vault"
  location = "West Europe"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_backup_policy_vm_workload" "policy" {
  name                = "tfex-backup-policy"
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name

  workload_type = "SAPAseDatabase"

  backup {
    frequency     = "Daily"
    time          = "02:00"
    weekdays      = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
  }

  retention_daily {
    count = 30
  }
}

resource "azurerm_virtual_machine" "vm" {
  # VM configuration for SAP workload
  name                = "example-vm"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  # ... other VM configuration
}

resource "azurerm_backup_container_vm_app" "container" {
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  workload_type       = "SAPAseDatabase"
  source_resource_id  = azurerm_virtual_machine.vm.id
}

resource "azurerm_backup_protected_vm_workload" "example" {
  protected_item_name = "SAPAseDatabase;DBInstanceName;SYSTEMDB"
  source_vm_id        = azurerm_virtual_machine.vm.id
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  backup_policy_id    = azurerm_backup_policy_vm_workload.policy.id
  workload_type       = "SAPAseDatabase"

  depends_on = [azurerm_backup_container_vm_app.container]
}
```

## Example with Protection State Management

```hcl
resource "azurerm_backup_protected_vm_workload" "example_with_protection" {
  protected_item_name = "SYSTEMDB"
  source_vm_id        = azurerm_virtual_machine.vm.id
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  backup_policy_id    = azurerm_backup_policy_vm_workload.policy.id
  workload_type       = "SAPAseDatabase"
  protection_state    = "ProtectionStopped"

  depends_on = [azurerm_backup_container_vm_app.container]
}
```

## Argument Reference

The following arguments are supported:

* `protected_item_name` - (Required) The name of the protected item (database/workload) to be protected. Changing this forces a new resource to be created. Format is `workloadType;DatabaseInstanceName;DatabaseName`

* `source_vm_id` - (Required) The ID of the Virtual Machine which contains the workload to be protected. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the Recovery Services Vault is located. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) The name of the Recovery Services Vault where the workload will be protected. Changing this forces a new resource to be created.

* `backup_policy_id` - (Required) The ID of the backup policy to be used for this protected workload.

* `workload_type` - (Required) The type of the workload to protect. Possible values include `SAPAseDatabase`. Changing this forces a new resource to be created.

* `protection_state` - (Optional) Specifies Protection state of the backup. Possible values are `Protected` and `ProtectionStopped`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Backup Protected VM Workload.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 120 minutes) Used when creating the Backup Protected VM Workload.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Protected VM Workload.
* `update` - (Defaults to 120 minutes) Used when updating the Backup Protected VM Workload.
* `delete` - (Defaults to 80 minutes) Used when deleting the Backup Protected VM Workload.

## Import

Backup Protected VM Workloads can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_protected_vm_workload.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupFabrics/Azure/protectionContainers/VMAppContainer;compute;group1;vm1/protectedItems/SAPAseDatabase;vm1;SYSTEMDB"
```

Note the ID requires quoting as there are semicolons.
