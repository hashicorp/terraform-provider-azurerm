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
      purge_soft_delete_on_destroy         = true
      recover_soft_deleted_api_managements = true
    }

    application_insights {
      disable_generated_rule = false
    }

    cognitive_account {
      purge_soft_delete_on_destroy = true
    }

    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }

    log_analytics_workspace {
      permanently_delete_on_destroy = true
    }

    resource_group {
      prevent_deletion_if_contains_resources = true
    }

    template_deployment {
      delete_nested_items_during_deletion = true
    }

    virtual_machine {
      delete_os_disk_on_deletion     = true
      graceful_shutdown              = false
      skip_shutdown_and_force_delete = false
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

* `application_insights` - (Optional) An `application_insights` block as defined below.

* `cognitive_account` - (Optional) A `cognitive_account` block as defined below.

* `key_vault` - (Optional) A `key_vault` block as defined below.

* `log_analytics_workspace` - (Optional) A `log_analytics_workspace` block as defined below.

* `resource_group` - (Optional) A `resource_group` block as defined below.

* `template_deployment` - (Optional) A `template_deployment` block as defined below.

* `virtual_machine` - (Optional) A `virtual_machine` block as defined below.

* `virtual_machine_scale_set` - (Optional) A `virtual_machine_scale_set` block as defined below.

---

The `api_management` block supports the following:

* `purge_soft_delete_on_destroy` - (Optional) Should the `azurerm_api_management` resources be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `recover_soft_deleted_api_managements` - (Optional) Should the `azurerm_api_management` resources recover a Soft-Deleted API Management service? Defaults to `true`

---

The `application_insights` block supports the following:

* `disable_generated_rule` - (Optional) Should the `azurerm_application_insights` resources disable the Azure generated Alert Rule during the create step? Defaults to `false`.

---

The `cognitive_account` block supports the following:

* `purge_soft_delete_on_destroy` - (Optional) Should the `azurerm_cognitive_account` resources be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

---

The `key_vault` block supports the following:

* `purge_soft_delete_on_destroy` - (Optional) Should the `azurerm_key_vault` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

~> **Note:** When purge protection is enabled, a key vault or an object in the deleted state cannot be purged until the retention period (7-90 days) has passed.

* `purge_soft_deleted_certificates_on_destroy` - (Optional) Should the `azurerm_key_vault_certificate` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `purge_soft_deleted_keys_on_destroy` - (Optional) Should the `azurerm_key_vault_key` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `purge_soft_deleted_secrets_on_destroy` - (Optional) Should the `azurerm_key_vault_secret` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `purge_soft_deleted_hardware_security_modules_on_destroy` - (Optional) Should the `azurerm_key_vault_managed_hardware_security_module` resource be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

* `recover_soft_deleted_certificates` - (Optional) Should the `azurerm_key_vault_certificate` resource recover a Soft-Deleted Certificate? Defaults to `true`.

* `recover_soft_deleted_key_vaults` - (Optional) Should the `azurerm_key_vault` resource recover a Soft-Deleted Key Vault? Defaults to `true`.

* `recover_soft_deleted_keys` - (Optional) Should the `azurerm_key_vault_key` resource recover a Soft-Deleted Key? Defaults to `true`.

* `recover_soft_deleted_secrets` - (Optional) Should the `azurerm_key_vault_secret` resource recover a Soft-Deleted Secret? Defaults to `true`.

~> **Note:** When recovering soft-deleted Key Vault items (Keys, Certificates, and Secrets) the Principal used by Terraform needs the `"recover"` permission.

---

The `log_analytics_workspace` block supports the following:

* `permanently_delete_on_destroy` - (Optional) Should the `azurerm_log_analytics_workspace` be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

---

The `resource_group` block supports the following:

* `prevent_deletion_if_contains_resources` - (Optional) Should the `azurerm_resource_group` resource check that there are no Resources within the Resource Group during deletion? This means that all Resources within the Resource Group must be deleted prior to deleting the Resource Group. Defaults to `false`.

-> **Note:** This will be defaulted to `true` in the next major version of the Azure Provider (3.0).

---

The `template_deployment` block supports the following:

* `delete_nested_items_during_deletion` - (Optional) Should the `azurerm_resource_group_template_deployment` resource attempt to delete resources that have been provisioned by the ARM Template, when the Resource Group Template Deployment is deleted? Defaults to `true`.

---

The `virtual_machine` block supports the following:

* `delete_os_disk_on_deletion` - (Optional) Should the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` resources delete the OS Disk attached to the Virtual Machine when the Virtual Machine is destroyed? Defaults to `true`.

~> **Note:** This does not affect the older `azurerm_virtual_machine` resource, which has its own flags for managing this within the resource.

* `graceful_shutdown` - (Optional) Should the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` request a graceful shutdown when the Virtual Machine is destroyed? Defaults to `false`.

~> **Note:** When using a graceful shutdown, Azure gives the Virtual Machine a 5 minutes window in which to complete the shutdown process, at which point the machine will be force powered off - [more information can be found in this blog post](https://azure.microsoft.com/en-us/blog/linux-and-graceful-shutdowns-2/).

* `skip_shutdown_and_force_delete` - Should the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` skip the shutdown command and `Force Delete`, this provides the ability to forcefully and immediately delete the VM and detach all sub-resources associated with the virtual machine. This allows those freed resources to be reattached to another VM instance or deleted. Defaults to `false`.

~> **Note:** Support for Force Delete is in an opt-in Preview.

---

The `virtual_machine_scale_set` block supports the following:

* `force_delete` - Should the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources `Force Delete`, this provides the ability to forcefully and immediately delete the VM and detach all sub-resources associated with the virtual machine. This allows those freed resources to be reattached to another VM instance or deleted. Defaults to `false`.

~> **Note:** Support for Force Delete is in an opt-in Preview.

* `roll_instances_when_required` - (Optional) Should the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources automatically roll the instances in the Scale Set when Required (for example when updating the Sku/Image). Defaults to `true`.

* `scale_to_zero_before_deletion` - (Optional) Should the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources scale to 0 instances before deleting the resource. Defaults to `true`.
