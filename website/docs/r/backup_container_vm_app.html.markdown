---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_container_vm_app"
description: |-
    Manages a VM application container in an Azure Recovery Vault
---

# azurerm_backup_container_vm_app

Manages registration of a VM application (SAP ASE/HANA Database) with Azure Backup. VM applications must be registered with an Azure Recovery Vault to create a protection container before the workloads within the VM can be backed up. Once the container is created, individual databases within the VM can be backed up using the `azurerm_backup_protected_vm_workload` resource.

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

resource "azurerm_virtual_machine" "vm" {
  name                = "example-vm"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  vm_size             = "Standard_DS2_v2"

  # VM configuration for SAP workload
  storage_image_reference {
    publisher = "SUSE"
    offer     = "SLES-SAP"
    sku       = "12-SP5"
    version   = "latest"
  }

  storage_os_disk {
    name              = "example-vm-os-disk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "example-vm"
    admin_username = "adminuser"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  network_interface_ids = [azurerm_network_interface.example.id]
}

resource "azurerm_backup_container_vm_app" "example" {
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  source_resource_id  = azurerm_virtual_machine.vm.id
  resource_group_name = azurerm_resource_group.example.name
  workload_type       = "SAPHanaDatabase"
}
```

## Argument Reference

The following arguments are supported:

* `recovery_vault_name` - (Required) The name of the Recovery Services Vault where the VM application will be registered. Changing this forces a new resource to be created.

* `source_resource_id` - (Required) The ID of the Virtual Machine which contains the application workloads to be protected. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the Recovery Services Vault is located. Changing this forces a new resource to be created.

* `workload_type` - (Required) The type of workload for the container. Possible values include `SAPAseDatabase` and `SAPHanaDatabase`. Changing this forces a new resource to be created.

-> **Note:** The Virtual Machine must have the appropriate SAP workload (ASE or HANA) installed and configured before registering the container.

-> **Note:** Azure Backup may place additional requirements on the Virtual Machine configuration for backup operations to succeed. Please find Prerequisites for registration here: [ASE](https://learn.microsoft.com/en-us/azure/backup/sap-ase-database-backup#prerequisites), [HANA](https://learn.microsoft.com/en-us/azure/backup/backup-azure-sap-hana-database#prerequisites)

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Backup Container VM App.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 120 minutes) Used when creating the Backup Container VM App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Container VM App.
* `delete` - (Defaults to 80 minutes) Used when deleting the Backup Container VM App.

## Import

Backup Container VM Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_container_vm_app.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupFabrics/Azure/protectionContainers/VMAppContainer;Compute;group1;vm1"
```

Note the ID requires quoting as there are semicolons.
