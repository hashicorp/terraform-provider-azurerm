---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_replicated_vm"
description: |-
    Manages a VM protected with Azure Site Recovery on Azure.
---

# azurerm_site_recovery_replicated_vm

Manages a VM replicated using Azure Site Recovery (Azure to Azure only). A replicated VM keeps a copiously updated image of the VM in another region in order to be able to start the VM in that region in case of a disaster.

## Example Usage

```hcl
resource "azurerm_resource_group" "primary" {
  name     = "tfex-replicated-vm-primary"
  location = "West US"
}

resource "azurerm_resource_group" "secondary" {
  name     = "tfex-replicated-vm-secondary"
  location = "East US"
}

resource "azurerm_virtual_machine" "vm" {
  name                  = "vm"
  location              = azurerm_resource_group.primary.location
  resource_group_name   = azurerm_resource_group.primary.name
  vm_size               = "Standard_B1s"
  network_interface_ids = [azurerm_network_interface.vm.id]

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "vm-os-disk"
    os_type           = "Linux"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    admin_username = "test-admin-123"
    admin_password = "test-pwd-123"
    computer_name  = "vm"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
  sku                 = "Standard"
}

resource "azurerm_site_recovery_fabric" "primary" {
  name                = "primary-fabric"
  resource_group_name = azurerm_resource_group.secondary.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  location            = azurerm_resource_group.primary.location
}

resource "azurerm_site_recovery_fabric" "secondary" {
  name                = "secondary-fabric"
  resource_group_name = azurerm_resource_group.secondary.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  location            = azurerm_resource_group.secondary.location
}

resource "azurerm_site_recovery_protection_container" "primary" {
  name                 = "primary-protection-container"
  resource_group_name  = azurerm_resource_group.secondary.name
  recovery_vault_name  = azurerm_recovery_services_vault.vault.name
  recovery_fabric_name = azurerm_site_recovery_fabric.primary.name
}

resource "azurerm_site_recovery_protection_container" "secondary" {
  name                 = "secondary-protection-container"
  resource_group_name  = azurerm_resource_group.secondary.name
  recovery_vault_name  = azurerm_recovery_services_vault.vault.name
  recovery_fabric_name = azurerm_site_recovery_fabric.secondary.name
}

resource "azurerm_site_recovery_replication_policy" "policy" {
  name                                                 = "policy"
  resource_group_name                                  = azurerm_resource_group.secondary.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.vault.name
  recovery_point_retention_in_minutes                  = 24 * 60
  application_consistent_snapshot_frequency_in_minutes = 4 * 60
}

resource "azurerm_site_recovery_protection_container_mapping" "container-mapping" {
  name                                      = "container-mapping"
  resource_group_name                       = azurerm_resource_group.secondary.name
  recovery_vault_name                       = azurerm_recovery_services_vault.vault.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.primary.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.primary.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.secondary.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.policy.id
}

resource "azurerm_site_recovery_network_mapping" "network-mapping" {
  name                        = "network-mapping"
  resource_group_name         = azurerm_resource_group.secondary.name
  recovery_vault_name         = azurerm_recovery_services_vault.vault.name
  source_recovery_fabric_name = azurerm_site_recovery_fabric.primary.name
  target_recovery_fabric_name = azurerm_site_recovery_fabric.secondary.name
  source_network_id           = azurerm_virtual_network.primary.id
  target_network_id           = azurerm_virtual_network.secondary.id
}

resource "azurerm_storage_account" "primary" {
  name                     = "primaryrecoverycache"
  location                 = azurerm_resource_group.primary.location
  resource_group_name      = azurerm_resource_group.primary.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "primary" {
  name                = "network1"
  resource_group_name = azurerm_resource_group.primary.name
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.primary.location
}

resource "azurerm_virtual_network" "secondary" {
  name                = "network2"
  resource_group_name = azurerm_resource_group.secondary.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_resource_group.secondary.location
}

resource "azurerm_subnet" "primary" {
  name                 = "network1-subnet"
  resource_group_name  = azurerm_resource_group.primary.name
  virtual_network_name = azurerm_virtual_network.primary.name
  address_prefixes     = ["192.168.1.0/24"]
}

resource "azurerm_subnet" "secondary" {
  name                 = "network2-subnet"
  resource_group_name  = azurerm_resource_group.secondary.name
  virtual_network_name = azurerm_virtual_network.secondary.name
  address_prefixes     = ["192.168.2.0/24"]
}

resource "azurerm_public_ip" "primary" {
  name                = "vm-public-ip-primary"
  allocation_method   = "Static"
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
  sku                 = "Basic"
}

resource "azurerm_public_ip" "secondary" {
  name                = "vm-public-ip-secondary"
  allocation_method   = "Static"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
  sku                 = "Basic"
}

resource "azurerm_network_interface" "vm" {
  name                = "vm-nic"
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name

  ip_configuration {
    name                          = "vm"
    subnet_id                     = azurerm_subnet.primary.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.primary.id
  }
}

resource "azurerm_site_recovery_replicated_vm" "vm-replication" {
  name                                      = "vm-replication"
  resource_group_name                       = azurerm_resource_group.secondary.name
  recovery_vault_name                       = azurerm_recovery_services_vault.vault.name
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.primary.name
  source_vm_id                              = azurerm_virtual_machine.vm.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.policy.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.primary.name

  target_resource_group_id                = azurerm_resource_group.secondary.id
  target_recovery_fabric_id               = azurerm_site_recovery_fabric.secondary.id
  target_recovery_protection_container_id = azurerm_site_recovery_protection_container.secondary.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.vm.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.primary.id
    target_resource_group_id   = azurerm_resource_group.secondary.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id   = azurerm_network_interface.vm.id
    target_subnet_name            = azurerm_subnet.secondary.name
    recovery_public_ip_address_id = azurerm_public_ip.secondary.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.container-mapping,
    azurerm_site_recovery_network_mapping.network-mapping,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the replication for the replicated VM. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) The name of the vault that should be updated. Changing this forces a new resource to be created.

* `recovery_replication_policy_id` - (Required) Id of the policy to use for this replicated vm. Changing this forces a new resource to be created.

* `source_recovery_fabric_name` - (Required) Name of fabric that should contain this replication. Changing this forces a new resource to be created.

* `source_vm_id` - (Required) Id of the VM to replicate Changing this forces a new resource to be created.

* `source_recovery_protection_container_name` - (Required) Name of the protection container to use. Changing this forces a new resource to be created.

* `target_resource_group_id` - (Required) Id of resource group where the VM should be created when a failover is done. Changing this forces a new resource to be created.

* `target_recovery_fabric_id` - (Required) Id of fabric where the VM replication should be handled when a failover is done. Changing this forces a new resource to be created.

* `target_recovery_protection_container_id` - (Required) Id of protection container where the VM replication should be created when a failover is done. Changing this forces a new resource to be created.

* `target_availability_set_id` - (Optional) Id of availability set that the new VM should belong to when a failover is done.

* `target_zone` - (Optional) Specifies the Availability Zone where the Failover VM should exist. Changing this forces a new resource to be created.

* `managed_disk` - (Optional) One or more `managed_disk` block as defined below. Changing this forces a new resource to be created.

* `unmanaged_disk` - (Optional) One or more `unmanaged_disk` block as defined below. Changing this forces a new resource to be created.
 
* `target_edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Managed Kubernetes Cluster should exist. Changing this forces a new resource to be created.

* `target_proximity_placement_group_id` - (Optional) Id of Proximity Placement Group the new VM should belong to when a failover is done.

* `target_boot_diagnostic_storage_account_id` - (Optional) Id of the storage account which the new VM should used for boot diagnostic when a failover is done.

* `target_capacity_reservation_group_id` - (Optional) Id of the Capacity reservation group where the new VM should belong to when a failover is done.

* `target_virtual_machine_scale_set_id` - (Optional) Id of the Virtual Machine Scale Set which the new Vm should belong to when a failover is done.

* `target_virtual_machine_size` - (Optional) Specifies the size the Virtual Machine should have.

* `target_network_id` - (Optional) Network to use when a failover is done (recommended to set if any network_interface is configured for failover).

* `test_network_id` - (Optional) Network to use when a test failover is done.

* `network_interface` - (Optional) One or more `network_interface` block as defined below.

* `multi_vm_group_name` - (Optional) Name of group in which all machines will replicate together and have shared crash consistent and app-consistent recovery points when failed over.

---

A `managed_disk` block supports the following:

* `disk_id` - (Required) Id of disk that should be replicated. Changing this forces a new resource to be created.

* `staging_storage_account_id` - (Required) Storage account that should be used for caching. Changing this forces a new resource to be created.

* `target_resource_group_id` - (Required) Resource group disk should belong to when a failover is done. Changing this forces a new resource to be created.

* `target_disk_type` - (Required) What type should the disk be when a failover is done. Possible values are `Standard_LRS`, `Premium_LRS`, `StandardSSD_LRS` and `UltraSSD_LRS`. Changing this forces a new resource to be created.

* `target_replica_disk_type` - (Required) What type should the disk be that holds the replication data. Possible values are `Standard_LRS`, `Premium_LRS`, `StandardSSD_LRS` and `UltraSSD_LRS`. Changing this forces a new resource to be created.

* `target_disk_encryption_set_id` - (Optional) The Disk Encryption Set that the Managed Disk will be associated with. Changing this forces a new resource to be created.

-> **Note:** Creating replicated vm with `target_disk_encryption_set_id` wil take more time (up to 5 hours), please extend the `timeout` for `create`.

* `target_disk_encryption` - (Optional) A `target_disk_encryption` block as defined below.

---

A `unmanaged_disk` block supports the following:

* `disk_uri` - (Required) Id of disk that should be replicated. Changing this forces a new resource to be created.

* `staging_storage_account_id` - (Required) Storage account that should be used for caching. Changing this forces a new resource to be created.

* `target_storage_account_id` - (Required) Storage account disk should belong to when a failover is done. Changing this forces a new resource to be created.

---

A `network_interface` block supports the following:

* `source_network_interface_id` - (Optional) (Required if the network_interface block is specified) Id source network interface.

* `target_static_ip` - (Optional) Static IP to assign when a failover is done.

* `target_subnet_name` - (Optional) Name of the subnet to use when a failover is done.

* `recovery_load_balancer_backend_address_pool_ids` - (Optional) A list of IDs of Load Balancer Backend Address Pools to use when a failover is done.

* `recovery_public_ip_address_id` - (Optional) Id of the public IP object to use when a failover is done.

* `failover_test_static_ip` - (Optional) Static IP to assign when a test failover is done.

* `failover_test_subnet_name` - (Optional) Name of the subnet to use when a test failover is done.

* `failover_test_public_ip_address_id` - (Optional) Id of the public IP object to use when a test failover is done.

---

The `target_disk_encryption` block supports:

* `disk_encryption_key` - (Required) A `disk_encryption_key` block as defined below.

* `key_encryption_key` - (Optional) A `key_encryption_key` block as defined below.

---

The `disk_encryption_key` block supports:

* `secret_url` - (Required) The URL to the Key Vault Secret used as the Disk Encryption Key that the Managed Disk will be associated with. This can be found as `id` on the `azurerm_key_vault_secret` resource. Changing this forces a new resource to be created.

* `vault_id` - (Required) The ID of the Key Vault. This can be found as `id` on the `azurerm_key_vault` resource. Changing this forces a new resource to be created.

---

The `key_encryption_key` block supports:

* `key_url` - (Required) The URL to the Key Vault Key used as the Key Encryption Key that the Managed Disk will be associated with. This can be found as `id` on the `azurerm_key_vault_key` resource. Changing this forces a new resource to be created.

* `vault_id` - (Required) The ID of the Key Vault. This can be found as `id` on the `azurerm_key_vault` resource. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Replicated VM.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Site Recovery Replicated VM.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Replicated VM.
* `update` - (Defaults to 80 minutes) Used when updating the Site Recovery Replicated VM.
* `delete` - (Defaults to 80 minutes) Used when deleting the Site Recovery Replicated VM.

## Import

Site Recovery Replicated VM's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_replicated_vm.vmreplication /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric-name/replicationProtectionContainers/protection-container-name/replicationProtectedItems/vm-replication-name
```
