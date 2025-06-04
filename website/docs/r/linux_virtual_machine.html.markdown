---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_linux_virtual_machine"
description: |-
  Manages a Linux Virtual Machine.
---

# azurerm_linux_virtual_machine

Manages a Linux Virtual Machine.

## Disclaimers

-> **Note:** Terraform will automatically remove the OS Disk by default - this behaviour can be configured using the [features](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block) setting within the Provider block.

-> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text. Read more about [sensitive data](/docs/state/sensitive-data.html) in state.

-> **Note:** This resource does not support Unmanaged Disks. If you need to use Unmanaged Disks you can continue to use the [azurerm_virtual_machine](virtual_machine.html) resource instead.

-> **Note:** This resource does not support attaching existing OS Disks. You can instead capture an [image](image.html) of the OS Disk or continue to use the [azurerm_virtual_machine](virtual_machine.html) resource instead.

-> **Note:** In this release there's a known issue where the `public_ip_address` and `public_ip_addresses` fields may not be fully populated for Dynamic Public IP's.

!> **Note:** Due to a breaking change in the Azure API the `vm_agent_platform_updates_enabled` field is now a `Read-Only` field that is controlled by the platform. Its value cannot be set, modified, or updated.

## Example Usage

This example provisions a basic Linux Virtual Machine on an internal network. Additional examples of how to use the `azurerm_linux_virtual_machine` resource can be found in the [./examples/virtual-machines/linux](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/virtual-machines/linux) directory within the GitHub Repository.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                = "example-machine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = file("~/.ssh/id_rsa.pub")
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
```

## Argument Reference

The following arguments are supported:

* `admin_username` - (Required) The username of the local administrator used for the Virtual Machine. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Linux Virtual Machine should exist. Changing this forces a new resource to be created.

* `license_type` - (Optional) Specifies the License Type for this Virtual Machine. Possible values are `RHEL_BYOS`, `RHEL_BASE`, `RHEL_EUS`, `RHEL_SAPAPPS`, `RHEL_SAPHA`, `RHEL_BASESAPAPPS`, `RHEL_BASESAPHA`, `SLES_BYOS`, `SLES_SAP`, `SLES_HPC`, `UBUNTU_PRO`.

* `name` - (Required) The name of the Linux Virtual Machine. Changing this forces a new resource to be created.

* `network_interface_ids` - (Required). A list of Network Interface IDs which should be attached to this Virtual Machine. The first Network Interface ID in this list will be the Primary Network Interface on the Virtual Machine.

* `os_disk` - (Required) A `os_disk` block as defined below.

* `resource_group_name` - (Required) The name of the Resource Group in which the Linux Virtual Machine should be exist. Changing this forces a new resource to be created.

* `size` - (Required) The SKU which should be used for this Virtual Machine, such as `Standard_F2`.

---

* `additional_capabilities` - (Optional) A `additional_capabilities` block as defined below.

* `admin_password` - (Optional) The Password which should be used for the local-administrator on this Virtual Machine. Changing this forces a new resource to be created.

-> **Note:** When an `admin_password` is specified `disable_password_authentication` must be set to `false`.
~> **Note:** One of either `admin_password` or `admin_ssh_key` must be specified.

* `admin_ssh_key` - (Optional) One or more `admin_ssh_key` blocks as defined below. Changing this forces a new resource to be created.

~> **Note:** One of either `admin_password` or `admin_ssh_key` must be specified.

* `allow_extension_operations` - (Optional) Should Extension Operations be allowed on this Virtual Machine? Defaults to `true`.

* `availability_set_id` - (Optional) Specifies the ID of the Availability Set in which the Virtual Machine should exist. Changing this forces a new resource to be created.

* `boot_diagnostics` - (Optional) A `boot_diagnostics` block as defined below.

* `bypass_platform_safety_checks_on_user_schedule_enabled` - (Optional) Specifies whether to skip platform scheduled patching when a user schedule is associated with the VM. Defaults to `false`.

~> **Note:** `bypass_platform_safety_checks_on_user_schedule_enabled` can only be set to `true` when `patch_mode` is set to `AutomaticByPlatform`.

* `capacity_reservation_group_id` - (Optional) Specifies the ID of the Capacity Reservation Group which the Virtual Machine should be allocated to.

~> **Note:** `capacity_reservation_group_id` cannot be used with `availability_set_id` or `proximity_placement_group_id`

* `computer_name` - (Optional) Specifies the Hostname which should be used for this Virtual Machine. If unspecified this defaults to the value for the `name` field. If the value of the `name` field is not a valid `computer_name`, then you must specify `computer_name`. Changing this forces a new resource to be created.

* `custom_data` - (Optional) The Base64-Encoded Custom Data which should be used for this Virtual Machine. Changing this forces a new resource to be created.

* `dedicated_host_id` - (Optional) The ID of a Dedicated Host where this machine should be run on. Conflicts with `dedicated_host_group_id`.

* `dedicated_host_group_id` - (Optional) The ID of a Dedicated Host Group that this Linux Virtual Machine should be run within. Conflicts with `dedicated_host_id`.

* `disable_password_authentication` - (Optional) Should Password Authentication be disabled on this Virtual Machine? Defaults to `true`. Changing this forces a new resource to be created.

-> **Note:** In general we'd recommend using SSH Keys for authentication rather than Passwords - but there's tradeoff's to each - please [see this thread for more information](https://security.stackexchange.com/questions/69407/why-is-using-an-ssh-key-more-secure-than-using-passwords).

-> **Note:** When an `admin_password` is specified `disable_password_authentication` must be set to `false`.

* `disk_controller_type` - (Optional) Specifies the Disk Controller Type used for this Virtual Machine. Possible values are `SCSI` and `NVMe`.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Linux Virtual Machine should exist. Changing this forces a new Linux Virtual Machine to be created.

* `encryption_at_host_enabled` - (Optional) Should all of the disks (including the temp disk) attached to this Virtual Machine be encrypted by enabling Encryption at Host?

* `eviction_policy` - (Optional) Specifies what should happen when the Virtual Machine is evicted for price reasons when using a Spot instance. Possible values are `Deallocate` and `Delete`. Changing this forces a new resource to be created.

-> **Note:** This can only be configured when `priority` is set to `Spot`.

* `extensions_time_budget` - (Optional) Specifies the duration allocated for all extensions to start. The time duration should be between 15 minutes and 120 minutes (inclusive) and should be specified in ISO 8601 format. Defaults to `PT1H30M`.

* `gallery_application` - (Optional) One or more `gallery_application` blocks as defined below.

~> **Note:** Gallery Application Assignments can be defined either directly on `azurerm_linux_virtual_machine` resource, or using the `azurerm_virtual_machine_gallery_application_assignment` resource - but the two approaches cannot be used together. If both are used with the same Virtual Machine, spurious changes will occur. If `azurerm_virtual_machine_gallery_application_assignment` is used, it's recommended to use `ignore_changes` for the `gallery_application` block on the corresponding `azurerm_linux_virtual_machine` resource, to avoid a persistent diff when using this resource.

* `identity` - (Optional) An `identity` block as defined below.

* `patch_assessment_mode` - (Optional) Specifies the mode of VM Guest Patching for the Virtual Machine. Possible values are `AutomaticByPlatform` or `ImageDefault`. Defaults to `ImageDefault`.

-> **Note:** If the `patch_assessment_mode` is set to `AutomaticByPlatform` then the `provision_vm_agent` field must be set to `true`.

* `patch_mode` - (Optional) Specifies the mode of in-guest patching to this Linux Virtual Machine. Possible values are `AutomaticByPlatform` and `ImageDefault`. Defaults to `ImageDefault`. For more information on patch modes please see the [product documentation](https://docs.microsoft.com/azure/virtual-machines/automatic-vm-guest-patching#patch-orchestration-modes).

-> **Note:** If `patch_mode` is set to `AutomaticByPlatform` then `provision_vm_agent` must also be set to `true`.

* `max_bid_price` - (Optional) The maximum price you're willing to pay for this Virtual Machine, in US Dollars; which must be greater than the current spot price. If this bid price falls below the current spot price the Virtual Machine will be evicted using the `eviction_policy`. Defaults to `-1`, which means that the Virtual Machine should not be evicted for price reasons.

-> **Note:** This can only be configured when `priority` is set to `Spot`.

* `plan` - (Optional) A `plan` block as defined below. Changing this forces a new resource to be created.

* `platform_fault_domain` - (Optional) Specifies the Platform Fault Domain in which this Linux Virtual Machine should be created. Defaults to `-1`, which means this will be automatically assigned to a fault domain that best maintains balance across the available fault domains. Changing this forces a new Linux Virtual Machine to be created.

* `priority` - (Optional) Specifies the priority of this Virtual Machine. Possible values are `Regular` and `Spot`. Defaults to `Regular`. Changing this forces a new resource to be created.

* `provision_vm_agent` - (Optional) Should the Azure VM Agent be provisioned on this Virtual Machine? Defaults to `true`. Changing this forces a new resource to be created.

~> **Note:** If `provision_vm_agent` is set to `false` then `allow_extension_operations` must also be set to `false`.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group which the Virtual Machine should be assigned to.

* `reboot_setting` - (Optional) Specifies the reboot setting for platform scheduled patching. Possible values are `Always`, `IfRequired` and `Never`.

~> **Note:** `reboot_setting` can only be set when `patch_mode` is set to `AutomaticByPlatform`.

* `secret` - (Optional) One or more `secret` blocks as defined below.

* `secure_boot_enabled` - (Optional) Specifies whether secure boot should be enabled on the virtual machine. Changing this forces a new resource to be created.

* `source_image_id` - (Optional) The ID of the Image which this Virtual Machine should be created from. Changing this forces a new resource to be created. Possible Image ID types include `Image ID`s, `Shared Image ID`s, `Shared Image Version ID`s, `Community Gallery Image ID`s, `Community Gallery Image Version ID`s, `Shared Gallery Image ID`s and `Shared Gallery Image Version ID`s.

-> **Note:** One of either `source_image_id` or `source_image_reference` must be set.

* `source_image_reference` - (Optional) A `source_image_reference` block as defined below. Changing this forces a new resource to be created.

-> **Note:** One of either `source_image_id` or `source_image_reference` must be set.

* `tags` - (Optional) A mapping of tags which should be assigned to this Virtual Machine.

* `os_image_notification` - (Optional) A `os_image_notification` block as defined below.

* `termination_notification` - (Optional) A `termination_notification` block as defined below.

* `user_data` - (Optional) The Base64-Encoded User Data which should be used for this Virtual Machine.

* `vtpm_enabled` - (Optional) Specifies whether vTPM should be enabled on the virtual machine. Changing this forces a new resource to be created.

* `virtual_machine_scale_set_id` - (Optional) Specifies the Orchestrated Virtual Machine Scale Set that this Virtual Machine should be created within.

-> **Note:** To update `virtual_machine_scale_set_id` the Preview Feature `Microsoft.Compute/SingleFDAttachDetachVMToVmss` needs to be enabled, see [the documentation](https://review.learn.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-attach-detach-vm#enroll-in-the-preview) for more information.

~> **Note:** Orchestrated Virtual Machine Scale Sets can be provisioned using [the `azurerm_orchestrated_virtual_machine_scale_set` resource](/docs/providers/azurerm/r/orchestrated_virtual_machine_scale_set.html).

~> **Note:** To attach an existing VM to a Virtual Machine Scale Set, the scale set must have `single_placement_group` set to `false`, see [the documentation](https://learn.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-attach-detach-vm?tabs=portal-1%2Cportal-2%2Cportal-3#limitations-for-attaching-an-existing-vm-to-a-scale-set) for more information.

* `zone` - (Optional) Specifies the Availability Zones in which this Linux Virtual Machine should be located. Changing this forces a new Linux Virtual Machine to be created.

---

A `additional_capabilities` block supports the following:

* `ultra_ssd_enabled` - (Optional) Should the capacity to enable Data Disks of the `UltraSSD_LRS` storage account type be supported on this Virtual Machine? Defaults to `false`.

* `hibernation_enabled` - (Optional) Whether to enable the hibernation capability or not.

---

A `admin_ssh_key` block supports the following:

* `public_key` - (Required) The Public Key which should be used for authentication, which needs to be in `ssh-rsa` format with at least 2048-bit or in `ssh-ed25519` format. Changing this forces a new resource to be created.

* `username` - (Required) The Username for which this Public SSH Key should be configured. Changing this forces a new resource to be created.

-> **Note:** The Azure VM Agent only allows creating SSH Keys at the path `/home/{username}/.ssh/authorized_keys` - as such this public key will be written to the authorized keys file.

---

A `boot_diagnostics` block supports the following:

* `storage_account_uri` - (Optional) The Primary/Secondary Endpoint for the Azure Storage Account which should be used to store Boot Diagnostics, including Console Output and Screenshots from the Hypervisor.

-> **Note:** Passing a null value will utilize a Managed Storage Account to store Boot Diagnostics

---

A `certificate` block supports the following:

* `url` - (Required) The Secret URL of a Key Vault Certificate.

-> **Note:** This can be sourced from the `secret_id` field within the `azurerm_key_vault_certificate` Resource.

---

A `diff_disk_settings` block supports the following:

* `option` - (Required) Specifies the Ephemeral Disk Settings for the OS Disk. At this time the only possible value is `Local`. Changing this forces a new resource to be created.

* `placement` - (Optional) Specifies where to store the Ephemeral Disk. Possible values are `CacheDisk` and `ResourceDisk`. Defaults to `CacheDisk`. Changing this forces a new resource to be created.

---

A `gallery_application` block supports the following:

* `version_id` - (Required) Specifies the Gallery Application Version resource ID.

* `automatic_upgrade_enabled` - (Optional) Specifies whether the version will be automatically updated for the VM when a new Gallery Application version is available in PIR/SIG. Defaults to `false`.

* `configuration_blob_uri` - (Optional) Specifies the URI to an Azure Blob that will replace the default configuration for the package if provided.

* `order` - (Optional) Specifies the order in which the packages have to be installed. Possible values are between `0` and `2147483647`. Defaults to `0`.

* `tag` - (Optional) Specifies a passthrough value for more generic context. This field can be any valid `string` value.

* `treat_failure_as_deployment_failure_enabled` - (Optional) Specifies whether any failure for any operation in the VmApplication will fail the deployment of the VM. Defaults to `false`.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Linux Virtual Machine. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Linux Virtual Machine.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `os_disk` block supports the following:

* `caching` - (Required) The Type of Caching which should be used for the Internal OS Disk. Possible values are `None`, `ReadOnly` and `ReadWrite`.

* `storage_account_type` - (Required) The Type of Storage Account which should back this the Internal OS Disk. Possible values are `Standard_LRS`, `StandardSSD_LRS`, `Premium_LRS`, `StandardSSD_ZRS` and `Premium_ZRS`. Changing this forces a new resource to be created.

* `diff_disk_settings` - (Optional) A `diff_disk_settings` block as defined above. Changing this forces a new resource to be created.

-> **Note:** `diff_disk_settings` can only be set when `caching` is set to `ReadOnly`. More information can be found [here](https://docs.microsoft.com/azure/virtual-machines/ephemeral-os-disks-deploy#vm-template-deployment)

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to Encrypt this OS Disk. Conflicts with `secure_vm_disk_encryption_set_id`.

-> **Note:** The Disk Encryption Set must have the `Reader` Role Assignment scoped on the Key Vault - in addition to an Access Policy to the Key Vault

* `disk_size_gb` - (Optional) The Size of the Internal OS Disk in GB, if you wish to vary from the size used in the image this Virtual Machine is sourced from.

-> **Note:** If specified this must be equal to or larger than the size of the Image the Virtual Machine is based on. When creating a larger disk than exists in the image you'll need to repartition the disk to use the remaining space.

* `name` - (Optional) The name which should be used for the Internal OS Disk. Changing this forces a new resource to be created.

* `secure_vm_disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to Encrypt this OS Disk when the Virtual Machine is a Confidential VM. Conflicts with `disk_encryption_set_id`. Changing this forces a new resource to be created.

~> **Note:** `secure_vm_disk_encryption_set_id` can only be specified when `security_encryption_type` is set to `DiskWithVMGuestState`.

* `security_encryption_type` - (Optional) Encryption Type when the Virtual Machine is a Confidential VM. Possible values are `VMGuestStateOnly` and `DiskWithVMGuestState`. Changing this forces a new resource to be created.

~> **Note:** `vtpm_enabled` must be set to `true` when `security_encryption_type` is specified.

~> **Note:** `encryption_at_host_enabled` cannot be set to `true` when `security_encryption_type` is set to `DiskWithVMGuestState`.

* `write_accelerator_enabled` - (Optional) Should Write Accelerator be Enabled for this OS Disk? Defaults to `false`.

-> **Note:** This requires that the `storage_account_type` is set to `Premium_LRS` and that `caching` is set to `None`.

---

A `plan` block supports the following:

* `name` - (Required) Specifies the Name of the Marketplace Image this Virtual Machine should be created from. Changing this forces a new resource to be created.

* `product` - (Required) Specifies the Product of the Marketplace Image this Virtual Machine should be created from. Changing this forces a new resource to be created.

* `publisher` - (Required) Specifies the Publisher of the Marketplace Image this Virtual Machine should be created from. Changing this forces a new resource to be created.

---

A `secret` block supports the following:

* `certificate` - (Required) One or more `certificate` blocks as defined above.

* `key_vault_id` - (Required) The ID of the Key Vault from which all Secrets should be sourced.

---

The `source_image_reference` block supports the following:

* `publisher` - (Required) Specifies the publisher of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `offer` - (Required) Specifies the offer of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `sku` - (Required) Specifies the SKU of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `version` - (Required) Specifies the version of the image used to create the virtual machines. Changing this forces a new resource to be created.

---

A `os_image_notification` block supports the following:

* `timeout` - (Optional) Length of time a notification to be sent to the VM on the instance metadata server till the VM gets OS upgraded. The only possible value is `PT15M`. Defaults to `PT15M`.

---

A `termination_notification` block supports the following:

* `enabled` - (Required) Should the termination notification be enabled on this Virtual Machine? 

* `timeout` - (Optional) Length of time (in minutes, between `5` and `15`) a notification to be sent to the VM on the instance metadata server till the VM gets deleted. The time duration should be specified in ISO 8601 format. Defaults to `PT5M`.

~> **Note:** For more information about the termination notification, please [refer to this doc](https://docs.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-terminate-notification).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Linux Virtual Machine.

* `identity` - An `identity` block as documented below.

* `os_disk` - An `os_disk` block as documented below.

* `private_ip_address` - The Primary Private IP Address assigned to this Virtual Machine.

* `private_ip_addresses` - A list of Private IP Addresses assigned to this Virtual Machine.

* `public_ip_address` - The Primary Public IP Address assigned to this Virtual Machine.

* `public_ip_addresses` - A list of the Public IP Addresses assigned to this Virtual Machine.

* `virtual_machine_id` - A 128-bit identifier which uniquely identifies this Virtual Machine.

* `vm_agent_platform_updates_enabled` - Are Virtual Machine Agent Platform Updates `enabled` on this Virtual Machine?

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

An `os_disk` block exports the following:

* `id` - The ID of the OS disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 45 minutes) Used when creating the Linux Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the Linux Virtual Machine.
* `update` - (Defaults to 45 minutes) Used when updating the Linux Virtual Machine.
* `delete` - (Defaults to 45 minutes) Used when deleting the Linux Virtual Machine.

## Import

Linux Virtual Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_linux_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1
```
