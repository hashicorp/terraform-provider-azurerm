---
layout: "azurerm"
page_title: "Azure Resource Manager: The Features Block"
description: |-
Azure Resource Manager: The Features Block

---

# The Features Block

The Azure Provider allows the behaviour of certain resources to be configured using the `features` block.

This allows different users to select the behaviour they require, for example some users may wish for the OS Disks for a Virtual Machine to be removed automatically when the Virtual Machine is destroyed - whereas other users may wish for these OS Disks to be detached but not deleted.

## Example Usage

If you wish to use the default behaviours of the Azure Provider, then you only need to define an empty `features` block as below:

```hcl
provider "azurerm" {
  features {}
}
```

Each of the blocks defined below can be optionally specified to configure the behaviour as needed - this example shows all the possible behaviours which can be configured:

```hcl
provider "azurerm" {
  features {
    api_management {
      purge_soft_delete_on_destroy = true
      recover_soft_deleted         = true
    }

    app_configuration {
      purge_soft_delete_on_destroy = true
      recover_soft_deleted         = true
    }

    application_insights {
      disable_generated_rule = false
    }

    cognitive_account {
      purge_soft_delete_on_destroy = true
    }

    databricks_workspace {
      force_delete = false
    }

    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }

    log_analytics_workspace {
      permanently_delete_on_destroy = true
    }

    machine_learning {
      purge_soft_deleted_workspace_on_destroy = true
    }

    managed_disk {
      expand_without_downtime = true
    }

    netapp {
      delete_backups_on_backup_vault_destroy = false
      prevent_volume_destruction             = true
    }

    postgresql_flexible_server {
      restart_server_on_configuration_value_change = true
    }

    recovery_service {
      vm_backup_stop_protection_and_retain_data_on_destroy    = true
      vm_backup_suspend_protection_and_retain_data_on_destroy = true
      purge_protected_items_from_vault_on_destroy             = true
    }

    resource_group {
      prevent_deletion_if_contains_resources = true
    }

    recovery_services_vault {
      recover_soft_deleted_backup_protected_vm = true
    }

    subscription {
      prevent_cancellation_on_destroy = false
    }

    template_deployment {
      delete_nested_items_during_deletion = true
    }

    virtual_machine {
      detach_implicit_data_disk_on_deletion = false
      delete_os_disk_on_deletion            = true
      graceful_shutdown                     = false
      skip_shutdown_and_force_delete        = false
    }

    virtual_machine_scale_set {
      force_delete                  = false
      roll_instances_when_required  = true
      scale_to_zero_before_deletion = true
    }
  }
}
```

## Arguments Reference

The `features` block supports the following:

* `api_management` - (Optional) An `api_management` block as defined below.

* `app_configuration` - (Optional) An `app_configuration` block as defined below.

* `application_insights` - (Optional) An `application_insights` block as defined below.

* `cognitive_account` - (Optional) A `cognitive_account` block as defined below.

* `databricks_workspace` - (Optional) A `databricks_workspace` block as defined below.

* `key_vault` - (Optional) A `key_vault` block as defined below.

* `log_analytics_workspace` - (Optional) A `log_analytics_workspace` block as defined below.

* `machine_learning` - (Optional) A `machine_learning` block as defined below.

* `managed_disk` - (Optional) A `managed_disk` block as defined below.

* `netapp` - (Optional) A `netapp` block as defined below.

* `recovery_service` - (Optional) A `recovery_service` block as defined below.

* `resource_group` - (Optional) A `resource_group` block as defined below.

* `recovery_services_vault` - (Optional) A `recovery_services_vault` block as defined below.

* `template_deployment` - (Optional) A `template_deployment` block as defined below.

* `virtual_machine` - (Optional) A `virtual_machine` block as defined below.

* `virtual_machine_scale_set` - (Optional) A `virtual_machine_scale_set` block as defined below.

---

The `api_management` block supports the following:

* `purge_soft_delete_on_destroy` - (Optional) Should the `azurerm_api_management` resources be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `recover_soft_deleted` - (Optional) Should the `azurerm_api_management` resources recover a Soft-Deleted API Management service? Defaults to `true`.

---

The `app_configuration` block supports the following:

* `purge_soft_delete_on_destroy` - (Optional) Should the `azurerm_app_configuration` resources be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `recover_soft_deleted` - (Optional) Should the `azurerm_app_configuration` resources recover a Soft-Deleted App Configuration service? Defaults to `true`.

---

The `application_insights` block supports the following:

* `disable_generated_rule` - (Optional) Should the `azurerm_application_insights` resources disable the Azure generated Alert Rule during the creation step? Defaults to `false`.

---

The `cognitive_account` block supports the following:

* `purge_soft_delete_on_destroy` - (Optional) Should the `azurerm_cognitive_account` resources be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

---

The `databricks_workspace` block supports the following:

* `force_delete` - (Optional) Should the managed resource group that contains the Unity Catalog data be forcibly deleted when the `azurerm_databricks_workspace` is destroyed? Defaults to `false`.

---

The `key_vault` block supports the following:

* `purge_soft_delete_on_destroy` - (Optional) Should the `azurerm_key_vault` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

~> **Note:** When purge protection is enabled, a key vault or an object in the deleted state cannot be purged until the retention period (7-90 days) has passed.

* `purge_soft_deleted_certificates_on_destroy` - (Optional) Should the `azurerm_key_vault_certificate` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `purge_soft_deleted_keys_on_destroy` - (Optional) Should the `azurerm_key_vault_key` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `purge_soft_deleted_secrets_on_destroy` - (Optional) Should the `azurerm_key_vault_secret` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `purge_soft_deleted_hardware_security_modules_on_destroy` - (Optional) Should the `azurerm_key_vault_managed_hardware_security_module` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `purge_soft_deleted_hardware_security_module_keys_on_destroy` - (Optional) Should the `azurerm_key_vault_managed_hardware_security_module_key` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `recover_soft_deleted_certificates` - (Optional) Should the `azurerm_key_vault_certificate` resource recover a Soft-Deleted Certificate? Defaults to `true`.

* `recover_soft_deleted_key_vaults` - (Optional) Should the `azurerm_key_vault` resource recover a Soft-Deleted Key Vault? Defaults to `true`.

* `recover_soft_deleted_keys` - (Optional) Should the `azurerm_key_vault_key` resource recover a Soft-Deleted Key? Defaults to `true`.

* `recover_soft_deleted_secrets` - (Optional) Should the `azurerm_key_vault_secret` resource recover a Soft-Deleted Secret? Defaults to `true`

* `recover_soft_deleted_hardware_security_module_keys` - (Optional) Should the `azurerm_key_vault_managed_hardware_security_module_key` resource recover a Soft-Deleted Key? Defaults to `true`.

~> **Note:** When recovering soft-deleted Key Vault items (Keys, Certificates, and Secrets) the Principal used by Terraform needs the `"recover"` permission.

---

The `log_analytics_workspace` block supports the following:

* `permanently_delete_on_destroy` - (Optional) Should the `azurerm_log_analytics_workspace` be permanently deleted (e.g. purged) when destroyed? Defaults to `false`.

---

The `machine_learning` block supports the following:

* `purge_soft_deleted_workspace_on_destroy` - (Optional) Should the `azurerm_machine_learning_workspace` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `false`.

---

The `managed_disk` block supports the following:

* `expand_without_downtime` - (Optional) Specifies whether Managed Disks which can be Expanded without Downtime (on either [a Linux VM](https://learn.microsoft.com/azure/virtual-machines/linux/expand-disks?tabs=azure-cli%2Cubuntu#expand-without-downtime) [or a Windows VM](https://learn.microsoft.com/azure/virtual-machines/windows/expand-os-disk#expand-without-downtime)) should be expanded without restarting the associated Virtual Machine. Defaults to `true`.

~> **Note:** Expand Without Downtime requires a specific configuration for the Managed Disk and Virtual Machine - Terraform will use Expand Without Downtime when the Managed Disk and Virtual Machine meet these requirements, and shut the Virtual Machine down as needed if this is inapplicable. More information on when Expand Without Downtime is applicable can be found in the [Linux VM](https://learn.microsoft.com/azure/virtual-machines/linux/expand-disks?tabs=azure-cli%2Cubuntu#expand-without-downtime) [or Windows VM](https://learn.microsoft.com/azure/virtual-machines/windows/expand-os-disk#expand-without-downtime) documentation.

---

The `netapp` block supports the following:

* `delete_backups_on_backup_vault_destroy` - (Optional) Should backups be deleted when an `azurerm_netapp_backup_vault` is being deleted? Defaults to `false`.
* `prevent_volume_destruction` - (Optional) Should an `azurerm_netapp_volume` be protected against deletion (intentionally or unintentionally)? Defaults to `true`.

---

The `postgresql_flexible_server` block supports the following:

* `restart_server_on_configuration_value_change` - (Optional) Should the `postgresql_flexible_server` restart after static server parameter change or removal? Defaults to `true`.

---

The `recovery_service` block supports the following:

* `vm_backup_stop_protection_and_retain_data_on_destroy` - (Optional) Should we retain the data and stop protection instead of destroying the backup protected vm? Defaults to `false`.

* `vm_backup_suspend_protection_and_retain_data_on_destroy` - (Optional) Should we retain the data and suspend protection instead of destroying the backup protected vm? Defaults to `false`.

* `purge_protected_items_from_vault_on_destroy` - (Optional) Should we purge all protected items when destroying the vault. Defaults to `false`.

---

The `resource_group` block supports the following:

* `prevent_deletion_if_contains_resources` - (Optional) Should the `azurerm_resource_group` resource check that there are no Resources within the Resource Group during deletion? This means that all Resources within the Resource Group must be deleted prior to deleting the Resource Group. Defaults to `true`.

---

The `recovery_services_vault` block supports the following:

* `recover_soft_deleted_backup_protected_vm` - (Optional) Should the `azurerm_backup_protected_vm` resource recover a Soft-Deleted protected VM? Defaults to `false`.

---

The `subscription` block supports the following:

* `prevent_cancellation_on_destroy` - (Optional) Should the `azurerm_subscription` resource prevent a subscription to be cancelled on destroy? Defaults to `false`.

---

The `template_deployment` block supports the following:

* `delete_nested_items_during_deletion` - (Optional) Should the `azurerm_resource_group_template_deployment` resource attempt to delete resources that have been provisioned by the ARM Template, when the Resource Group Template Deployment is deleted? Defaults to `true`.

---

The `virtual_machine` block supports the following:

* `detach_implicit_data_disk_on_deletion` - (Optional) Should we detach the `azurerm_virtual_machine_implicit_data_disk_from_source` from the virtual machine instead of destroying it? Defaults to `false`.

* `delete_os_disk_on_deletion` - (Optional) Should the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` resources delete the OS Disk attached to the Virtual Machine when the Virtual Machine is destroyed? Defaults to `true`.

~> **Note:** This does not affect the older `azurerm_virtual_machine` resource, which has its own flags for managing this within the resource.

* `graceful_shutdown` - (Optional) Should the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` request a graceful shutdown when the Virtual Machine is destroyed? Defaults to `false`.

!> **Note:** Due to a breaking API change `graceful_shutdown` is no longer effective and has been deprecated. This feature will be removed from v5.0 of the AzureRM provider.

* `skip_shutdown_and_force_delete` - Should the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` skip the shutdown command and `Force Delete`, this provides the ability to forcefully and immediately delete the VM and detach all sub-resources associated with the virtual machine. This allows those freed resources to be reattached to another VM instance or deleted. Defaults to `false`.

---

The `virtual_machine_scale_set` block supports the following:

* `force_delete` - Should the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources `Force Delete`, this provides the ability to forcefully and immediately delete the VM and detach all sub-resources associated with the virtual machine. This allows those freed resources to be reattached to another VM instance or deleted. Defaults to `false`.

~> **Note:** Support for Force Delete is in an opt-in Preview.

* `reimage_on_manual_upgrade` - (Optional) Should the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources automatically reimage during the update the instances in the Scale Set when `upgrade_mode` is `Manual`. Defaults to `true`.

* `roll_instances_when_required` - (Optional) Should the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources automatically roll the instances in the Scale Set when Required (for example when updating the Sku/Image). Defaults to `true`.

* `scale_to_zero_before_deletion` - (Optional) Should the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources scale to 0 instances before deleting the resource. Defaults to `true`.
