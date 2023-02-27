---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_hyperv_replicated_vm"
description: |-
    Manages a HyperV VM protected with Azure Site Recovery on Azure.
---

# azurerm_site_recovery_replicated_vm

Manages a HyperV VM replicated using Azure Site Recovery (HyperV to Azure only). A replicated VM keeps a copiously updated image of the VM in Azure in order to be able to start the VM in Azure in case of a disaster.

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

resource "azurerm_site_recovery_services_vault_hyperv_site" "example" {
  name              = "example-hyperv-site"
  recovery_vault_id = azurerm_recovery_services_vault.example.id
}

resource "azurerm_site_recovery_hyperv_replication_policy" "example" {
  recovery_vault_id                                  = azurerm_recovery_services_vault.example.id
  name                                               = "example-policy"
  recovery_point_retention_in_hours                  = 2
  application_consistent_snapshot_frequency_in_hours = 1
  replication_interval_in_seconds                    = 300
}

resource "azurerm_site_recovery_hyperv_replication_policy_association" "example" {
  name           = "example-association"
  hyperv_site_id = azurerm_site_recovery_services_vault_hyperv_site.example.id
  policy_id      = azurerm_site_recovery_hyperv_replication_policy.example.id
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
  resource_group_name  = azurerm_resource_group.primary.name
  virtual_network_name = azurerm_virtual_network.primary.name
  address_prefixes     = ["192.168.2.0/24"]
}

resource "azurerm_site_recovery_hyperv_replicated_vm" "example" {
  name                      = "example-hyperv-vm"
  hyperv_site_id            = azurerm_site_recovery_services_vault_hyperv_site.example.id
  source_vm_name            = "VM1"
  target_resource_group_id  = azurerm_resource_group.example.id
  target_vm_name            = "target-vm"
  target_storage_account_id = azurerm_storage_account.example.id
  replication_policy_id     = azurerm_site_recovery_hyperv_replication_policy.example.id
  os_type                   = "Windows"
  os_disk_name              = "VM1"
  target_network_id         = azurerm_virtual_network.example.id
  disks_to_include          = ["VM1"]

  depends_on = [azurerm_site_recovery_hyperv_replication_policy_association.example]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the replication for the replicated VM. Changing this forces a new resource to be created.

* `hyperv_site_id` - (Required) ID of the HyperV Site or VMM Cloud in the Recovery Vault. Changing this forces a new resource to be created.

* `source_vm_name` - (Required) The name of VM in HyperV. Changing this forces a new resource to be created.

* `target_resource_group_id` - (Required) ID of resource group where the VM should be created when a failover is done. Changing this forces a new resource to be created.

* `target_vm_name` - (Required) Name of the VM that should be created when a failover is done. Changing this forces a new resource to be created.

* `target_storage_account_id` - (Required) ID of the storage account that should be used to replicate VM disks when a failover is done. Changing this forces a new resource to be created.

* `replication_policy_id` - (Required) ID of the policy to use for this replicated VM. Changing this forces a new resource to be created.

* `os_type` - (Required) OS type of the source VM. Possible values are `Linux` and `Windows`. Changing this forces a new resource to be created.

* `os_disk_name` - (Required) Name of the OS disk of the source VM. Changing this forces a new resource to be created.

* `target_network_id` - (Required) Network to use when a failover is done.

* `network_interface` - (Required) One or more `network_interface` block as defined below.

* `disks_to_include` - (Optional) A list of disk names to include from the source VM in the replication. Changing this forces a new resource to be created.

* `use_managed_disk_enabled` - (Optional) Whether to use managed disks in replication. Changing this forces a new resource to be created. Defaults to `false`.

* `managed_disk` - (Optional) One or more `managed_disk` block as defined below. Changing this forces a new resource to be created.

-> **NOTE:** If `use_managed_disk_enabled` is set to `true`, then `managed_disk` block must be specified. If `use_managed_disk_enabled` is set to `false`, then `disks_to_include` block must be specified.

* `log_storage_account_id` - (Optional) ID of the Storage Account to be used for logging during replication. Changing this forces a new resource to be created.

* `enable_rdp_or_ssh_on_target_option` - (Optional) Options to enable RDP or SSH on replicated VM. Possible values are `Never`, `OnlyOnTestFailover` and `Always`.

* `target_vm_size` - (Optional) Size of the VM that should be created when a failover is done, such as `Standard_F2`. If it's not specified, it will automatically be set by detecting the source VM size.

* `target_availability_zone` - (Optional) The Availability Zone where the Failover VM should exist. 

* `target_availability_set_id` - (Optional) Id of availability set that the new VM should belong to when a failover is done.

* `license_type` - (Optional) The License Type for this VM. Possible values are `NoLicenseType`, `NotSpecified` and `WindowsServer`. 

* `sql_server_license_type` - (Optional)  The SQL Server License Type for this VM. Possible values are `NotSpecified`, `NoLicenseType`, `AHUB` and `PAYG`.

* `target_proximity_placement_group_id` - (Optional) ID of Proximity Placement Group the new VM should belong to when a failover is done.

* `target_vm_tags` - (Optional) A mapping of tags assigned to the replicated VM.

* `target_disk_tags` - (Optional) A mapping of tags assigned to the disks of the replicated VM.

* `target_network_interface_tags` - (Optional) A mapping of tags assigned to the Network Interfaces of the replicated VM.

---

A `managed_disk` block supports the following:

* `disk_name` - (Required) Name of the disks from the source VM. Changing this forces a new resource to be created.

* `target_disk_type` - (Required) What type should the disk be when a failover is done. Possible values are `Standard_LRS`, `Premium_LRS`and `StandardSSD_LRS`. Changing this forces a new resource to be created.

* `target_disk_encryption_set_id` - (Optional) The Disk Encryption Set that the Managed Disk will be associated with.

-> **NOTE:** Creating replicated vm with `target_disk_encryption_set_id` wil take more time (up to 5 hours), please extend the `timeout` for `create`. 

---

A `network_interface` block supports the following:

* `network_name` - (Required) Name of the Network in source VM.

* `target_static_ip` - (Optional) Static IP to assign when a failover is done.

* `target_subnet_name` - (Optional) Name of the subnet to use when a failover is done.

* `is_primary` - (Optional) Whether this `network_interface` is primary for the replicated VM. Defaults to `true`.

-> **NOTE:** Only allowed to set `is_primary` to `true` for one `network_interface`.

* `failover_enabled` - (Optional) Whether this `network_interface` should be created when a failover is done. Defaults to `true`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Replicated VM.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Site Recovery HyperV Replicated VM.
* `update` - (Defaults to 80 minutes) Used when updating the Site Recovery HyperV Replicated VM.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery HyperV Replicated VM.
* `delete` - (Defaults to 80 minutes) Used when deleting the Site Recovery HyperV Replicated VM.

## Import

Site Recovery Replicated VM's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_hyperv_replicated_vm.vmreplication /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric-name/replicationProtectionContainers/protection-container-name/replicationProtectedItems/vm-replication-name
```
