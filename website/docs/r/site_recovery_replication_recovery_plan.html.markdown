---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_replication_recovery_plan"
description: |-
    Manages a Site Recovery Replication Recovery Plan within a Recovery Services vault.
---

# azurerm_site_recovery_replication_recovery_plan

Manages a Site Recovery Replication Recovery Plan within a Recovery Services vault. A recovery plan gathers machines into recovery groups for the purpose of failover.

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

resource "azurerm_site_recovery_replication_recovery_plan" "example" {
  name                      = "example-recover-plan"
  recovery_vault_id         = azurerm_recovery_services_vault.vault.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.primary.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.secondary.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.vm-replication.id]
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Replication Plan. The name can contain only letters, numbers, and hyphens. It should start with a letter and end with a letter or a number. Can be a maximum of 63 characters. Changing this forces a new resource to be created.

* `recovery_vault_id` - (Required) The ID of the vault that should be updated. Changing this forces a new resource to be created.

* `source_recovery_fabric_id` - (Required) ID of source fabric to be recovered from. Changing this forces a new Replication Plan to be created.

* `target_recovery_fabric_id` - (Required) ID of target fabric to recover. Changing this forces a new Replication Plan to be created.

* `shutdown_recovery_group` - (Required) One `shutdown_recovery_group` block as defined below.

* `failover_recovery_group` - (Required) One `failover_recovery_group` block as defined below.

* `boot_recovery_group` - (Required) One or more `boot_recovery_group` blocks as defined below.

* `azure_to_azure_settings` - (Optional) An `azure_to_azure_settings` block as defined below.

---

A `shutdown_recovery_group` block supports the following:

* `pre_action` - (Optional) one or more `action` block as defined below. which will be executed before the group recovery.

* `post_action` - (Optional) one or more `action` block as defined below. which will be executed after the group recovery.

---

A `failover_recovery_group` block supports the following:

* `pre_action` - (Optional) one or more `action` block as defined below. which will be executed before the group recovery.

* `post_action` - (Optional) one or more `action` block as defined below. which will be executed after the group recovery.

---

A `boot_recovery_group` block supports the following:

* `replicated_protected_items` - (Optional) One or more protected VM IDs.

* `pre_action` - (Optional) one or more `action` block as defined below. which will be executed before the group recovery.

* `post_action` - (Optional) one or more `action` block as defined below. which will be executed after the group recovery.

---

An `action` block supports the following:

* `name` - (Required) Name of the Action.

* `type` - (Required) Type of the action detail. Possible values are `AutomationRunbookActionDetails`, `ManualActionDetails` and `ScriptActionDetails`.

* `fail_over_directions` - (Required) Directions of fail over. Possible values are `PrimaryToRecovery` and `RecoveryToPrimary`

* `fail_over_types` - (Required) Types of fail over. Possible values are `TestFailover`, `PlannedFailover` and `UnplannedFailover`

* `fabric_location` - (Optional) The fabric location of runbook or script. Possible values are `Primary` and `Recovery`. It must not be specified when `type` is `ManualActionDetails`.

-> **Note:** This is required when `type` is set to `AutomationRunbookActionDetails` or `ScriptActionDetails`.

* `runbook_id` - (Optional) Id of runbook.

-> **Note:** This property is required when `type` is set to `AutomationRunbookActionDetails`.

* `manual_action_instruction` - (Optional) Instructions of manual action.

-> **Note:** This property is required when `type` is set to `ManualActionDetails`.

* `script_path` - (Optional) Path of action script.

-> **Note:** This property is required when `type` is set to `ScriptActionDetails`.

---

An `azure_to_azure_settings` block supports the following:

* `primary_zone` - (Optional) The Availability Zone in which the VM is located. Changing this forces a new Site Recovery Replication Recovery Plan to be created.

* `recovery_zone` - (Optional) The Availability Zone in which the VM is recovered. Changing this forces a new Site Recovery Replication Recovery Plan to be created.

-> **Note:** `primary_zone` and `recovery_zone` must be specified together.

* `primary_edge_zone` - (Optional) The Edge Zone within the Azure Region where the VM exists. Changing this forces a new Site Recovery Replication Recovery Plan to be created.

* `recovery_edge_zone` - (Optional) The Edge Zone within the Azure Region where the VM is recovered. Changing this forces a new Site Recovery Replication Recovery Plan to be created.

-> **Note:** `primary_edge_zone` and `recovery_edge_zone` must be specified together.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Fabric.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Site Recovery Replication Plan.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Replication Plan.
* `update` - (Defaults to 30 minutes) Used when updating the Site Recovery Replication Plan.
* `delete` - (Defaults to 30 minutes) Used when deleting the Site Recovery Replication Plan.

## Import

Site Recovery Fabric can be imported using the `resource id`, e.g.

```shell
terraform import  azurerm_site_recovery_replication_recovery_plan.example /subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/groupName/providers/Microsoft.RecoveryServices/vaults/vaultName/replicationRecoveryPlans/planName
```
