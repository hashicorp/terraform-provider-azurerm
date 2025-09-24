---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_protected_vm_workload"
description: |-
    Manages an Azure VM workload-specific protected item representing SAP ASE Database.
---

# azurerm_backup_protected_vm_workload_sap_ase_database

Manages an Azure VM workload-specific protected item representing SAP ASE Database. The VM must be registered with a backup container before it can be protected.

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

resource "azurerm_backup_protected_vm_workload_sap_ase_database" "example" {
  database_name       = "SYSTEMDB"
  database_instance_name = "DBInstanceName"
  source_vm_id        = azurerm_virtual_machine.vm.id
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  backup_policy_id    = azurerm_backup_policy_vm_workload.policy.id

  depends_on = [azurerm_backup_container_vm_app.container]
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the Recovery Services Vault is located. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) The name of the Recovery Services Vault where the workload will be protected. Changing this forces a new resource to be created.

* `source_vm_id` - (Required) The ID of the Virtual Machine which contains the workload to be protected. Changing this forces a new resource to be created.

* `backup_policy_id` - (Optional) The ID of the backup policy to be used for this protected workload.

~> **Note:** `backup_policy_id` is required during initial creation of this resource.

~> **Note:** When `protection_state` is set to`ProtectionStopped`, the Azure API may not return `backup_policy_id`. To avoid a perpetual diff, use Terraform's [ignore_changes](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#ignore_changes) argument.

* `protection_state` - (Optional) Specifies Protection state of the backup. Possible values are `Protected` and `ProtectionStopped`.

* `database_name` - (Required) The SAP ASE database name. Changing this forces a new resource to be created.

* `database_instance_name` - (Required) The SAP ASE database instance name. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the SAP ASE (Sybase) protected item in Azure VM.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 120 minutes) Used when creating the Backup Protected VM SAP ASE database workload.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Protected VM SAP ASE database workload.
* `update` - (Defaults to 120 minutes) Used when updating the Backup Protected VM SAP ASE database workload.
* `delete` - (Defaults to 80 minutes) Used when deleting the Backup Protected VM SAP ASE database workload.

## Import

Backup Protected VM SAP ASE database workload can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_protected_vm_workload.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupFabrics/Azure/protectionContainers/VMAppContainer;compute;group1;vm1/protectedItems/SAPAseDatabase;vm1;SYSTEMDB"
```

Note the ID requires quoting as there are semicolons.
