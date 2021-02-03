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
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7.5"
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
    source_network_interface_id           = azurerm_network_interface.vm.id
    target_subnet_name                    = "network2-subnet"
    recovery_public_ip_address_id         = azurerm_public_ip.secondary.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network mapping.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located.

* `recovery_vault_name` - (Required) The name of the vault that should be updated.

* `source_recovery_fabric_name` - (Required) Name of fabric that should contains this replication.

* `source_vm_id` - (Required) Id of the VM to replicate

* `source_recovery_protection_container_name` - (Required) Name of the protection container to use.

* `target_resource_group_id` - (Required) Id of resource group where the VM should be created when a failover is done.

* `target_recovery_fabric_id` - (Required)  Id of fabric where the VM replication should be handled when a failover is done.

* `target_recovery_protection_container_id` - (Required)  Id of protection container where the VM replication should be created when a failover is done.

* `target_availability_set_id` - (Optional)  Id of availability set that the new VM should belong to when a failover is done.

* `managed_disk` - (Required) One or more `managed_disk` block.

* `target_network_id` - (Optional) Network to use when a failover is done (recommended to set if any network_interface is configured for failover). 

* `network_interface` - (Optional) One or more `network_interface` block.

---

A `managed_disk` block supports the following:

* `disk_id` - (Required) Id of disk that should be replicated.

* `staging_storage_account_id` - (Required) Storage account that should be used for caching.

* `target_resource_group_id` - (Required) Resource group disk should belong to when a failover is done.

* `target_disk_type` - (Required) What type should the disk be when a failover is done.

* `target_replica_disk_type` - (Required) What type should the disk be that holds the replication data.

---

A `network_interface` block supports the following:

* `source_network_interface_id` - (Required if the network_interface block is specified) Id source network interface.

* `target_static_ip` - (Optional) Static IP to assign when a failover is done.

* `target_subnet_name` - (Optional) Name of the subnet to to use when a failover is done.

* `recovery_public_ip_address_id` - (Optional) Id of the public IP object to use when a failover is done.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Replicated VM.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 120 minutes) Used when creating the Site Recovery Replicated VM.
* `update` - (Defaults to 80 minutes) Used when updating the Site Recovery Replicated VM.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Replicated VM.
* `delete` - (Defaults to 80 minutes) Used when deleting the Site Recovery Replicated VM.

## Import

Site Recovery Replicated VM's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_replicated_vm.vmreplication /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric-name/replicationProtectionContainers/protection-container-name/replicationProtectedItems/vm-replication-name
```
