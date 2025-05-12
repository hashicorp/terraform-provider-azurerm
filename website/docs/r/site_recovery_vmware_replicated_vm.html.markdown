---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_vmware_replicated_vm"
description: |-
    Manages a VMWare replicated VM protected with Azure Site Recovery on Azure.
---

# azurerm_site_recovery_vmware_replicated_vm

Manages a VMWare replicated VM using Azure Site Recovery (VMWare to Azure only). A replicated VM keeps a copiously updated image of the VM in Azure in order to be able to start the VM in Azure in case of a disaster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_site_recovery_vmware_replication_policy" "example" {
  recovery_vault_id                                    = azurerm_recovery_services_vault.example.id
  name                                                 = "example-policy"
  recovery_point_retention_in_minutes                  = 1440
  application_consistent_snapshot_frequency_in_minutes = 240
}

resource "azurerm_site_recovery_vmware_replication_policy_association" "test" {
  name              = "example-association"
  recovery_vault_id = azurerm_recovery_services_vault.example.id
  policy_id         = azurerm_site_recovery_vmware_replication_policy.example.id
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-net"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_resource_group.example.location
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["192.168.2.0/24"]
}

resource "azurerm_site_recovery_vmware_replicated_vm" "example" {
  name                                       = "example-vmware-vm"
  recovery_vault_id                          = azurerm_recovery_services_vault.example.id
  source_vm_name                             = "example-vm"
  appliance_name                             = "example-appliance"
  recovery_replication_policy_id             = azurerm_site_recovery_vmware_replication_policy_association.example.policy_id
  physical_server_credential_name            = "example-creds"
  license_type                               = "NotSpecified"
  target_boot_diagnostics_storage_account_id = azurerm_storage_account.example.id
  target_vm_name                             = "example_replicated_vm"
  target_resource_group_id                   = azurerm_resource_group.example.id
  default_log_storage_account_id             = azurerm_storage_account.example.id
  default_recovery_disk_type                 = "Standard_LRS"
  target_network_id                          = azurerm_virtual_network.example.id

  network_interface {
    source_mac_address = "00:00:00:00:00:00"
    target_subnet_name = azurerm_subnet.example.name
    is_primary         = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `appliance_name` - (Required) The name of VMWare appliance which handles the replication. Changing this forces a new resource to be created.

* `name` - (Required) The name of the replicated VM. Changing this forces a new resource to be created.

* `physical_server_credential_name` - (Required) The name of the credential to access the source VM. Changing this forces a new resource to be created. More information about the credentials could be found [here](https://learn.microsoft.com/en-us/azure/site-recovery/deploy-vmware-azure-replication-appliance-modernized).

* `recovery_vault_id` - (Required) The ID of the Recovery Services Vault where the replicated VM is created.

* `recovery_replication_policy_id` - (Required) The ID of the policy to use for this replicated VM.

* `source_vm_name` - (Required) The name of the source VM in VMWare. Changing this forces a new resource to be created.

* `target_resource_group_id` - (Required) The ID of resource group where the VM should be created when a failover is done.

* `target_vm_name` - (Required) Name of the VM that should be created when a failover is done. Changing this forces a new resource to be created.

* `target_network_id` - (Optional) The ID of network to use when a failover is done.

~> **Note:** `target_network_id` is required when `network_interface` is specified.

* `default_log_storage_account_id` - (Optional) The ID of the stroage account that should be used for logging during replication. 

~> **Note:** Only standard types of storage accounts are allowed.

~> **Note:** Only one of `default_log_storage_account_id` or `managed_disk` must be specified.

~> **Note:** Changing `default_log_storage_account_id` forces a new resource to be created. But removing it does not.

~> **Note:** When `default_log_storage_account_id` co-exist with `managed_disk`, the value of `default_log_storage_account_id` must be as same as `log_storage_account_id` of every `managed_disk` or it forces a new resource to be created.

* `default_target_disk_encryption_set_id` - (Optional) The ID of the default Disk Encryption Set that should be used for the disks when a failover is done.

~> **Note:** Changing `default_target_disk_encryption_set_id` forces a new resource to be created. But removing it does not.

~> **Note:** When `default_target_disk_encryption_set_id` co-exist with `managed_disk`, the value of `default_target_disk_encryption_set_id` must be as same as `target_disk_encryption_set_id` of every `managed_disk` or it forces a new resource to be created.

* `default_recovery_disk_type` - (Optional) The type of storage account that should be used for recovery disks when a failover is done. Possible values are `Premium_LRS`, `Standard_LRS` and `StandardSSD_LRS`.

~> **Note:** Only one of `default_recovery_disk_type` or `managed_disk` must be specified.

~> **Note:** Changing `default_recovery_disk_type` forces a new resource to be created. But removing it does not.

~> **Note:** When `default_recovery_disk_type` co-exist with `managed_disk`, the value of `default_recovery_disk_type` must be as same as `target_disk_type` of every `managed_disk` or it forces a new resource to be created.

* `license_type` - (Optional) The license type of the VM. Possible values are `NoLicenseType`, `NotSpecified` and `WindowsServer`. Defaults to `NotSpecified`.

* `multi_vm_group_name` - (Optional) Name of group in which all machines will replicate together and have shared crash consistent and app-consistent recovery points when failed over.

* `managed_disk` - (Optional) One or more `managed_disk` block as defined below. It's available only if mobility service is already installed on the source VM.

~> **Note:** A replicated VM could be created without `managed_disk` block, once the block has been specified, changing it expect removing it forces a new resource to be created.

* `network_interface` - (Optional) One or more `network_interface` block as defined below.

* `target_availability_set_id` - (Optional) The ID of availability set that the new VM should belong to when a failover is done.

* `target_boot_diagnostics_storage_account_id` - (Optional) The ID of the storage account that should be used for boot diagnostics when a failover is done.

* `target_proximity_placement_group_id` - (Optional) The ID of Proximity Placement Group the new VM should belong to when a failover is done.

~> **Note:** Only one of `target_availability_set_id` or `target_zone` can be specified.

* `target_zone` - (Optional) Specifies the Availability Zone where the Failover VM should exist.

* `target_vm_size` - (Optional) Size of the VM that should be created when a failover is done, such as `Standard_F2`. If it's not specified, it will automatically be set by detecting the source VM size.

* `test_network_id` - (Optional) The ID of network to use when a test failover is done.
---

A `managed_disk` block supports the following:

* `disk_id` - (Required) The ID of the disk to be replicated.

* `target_disk_type` - (Required) The disk type of the disk to be created when a failover is done. Possible values are `Premium_LRS`, `Standard_LRS` and `StandardSSD_LRS`.

* `log_storage_account_id` - (Optional) The ID of the storage account that should be used for logging during replication.

* `target_disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set that should be used for the disks when a failover is done. 

---

A `network_interface` block supports the following:

* `source_mac_address` - (Required) Mac address of the network interface of source VM. 

* `is_primary` - (Required) Whether this `network_interface` is primary for the replicated VM.

* `target_static_ip` - (Optional) Static IP to assign when a failover is done.

* `target_subnet_name` - (Optional) Name of the subnet to use when a failover is done.

* `test_subnet_name` - (Optional) Name of the subnet to use when a test failover is done.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Replicated VM.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Site Recovery HyperV Replicated VM.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery HyperV Replicated VM.
* `update` - (Defaults to 90 minutes) Used when updating the Site Recovery HyperV Replicated VM.
* `delete` - (Defaults to 90 minutes) Used when deleting the Site Recovery HyperV Replicated VM.

## Import

Site Recovery VMWare Replicated VM's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_vmware_replicated_vm.vmreplication /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric-name/replicationProtectionContainers/protection-container-name/replicationProtectedItems/vm-replication-name
```
