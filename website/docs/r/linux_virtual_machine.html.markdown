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

-> **Note** Terraform will automatically remove the OS Disk by default - this behaviour can be configured [using the `features` setting within the Provider block](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#features).

~> **Note** All arguments including the administrator login and password will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

-> **Note** This resource does not support Unmanaged Disks. If you need to use Unmanaged Disks you can continue to use [the `azurerm_virtual_machine` resource](virtual_machine.html) instead.

~> **Note** This resource does not support attaching existing OS Disks. You can instead [capture an image of the OS Disk](image.html) or continue to use [the `azurerm_virtual_machine` resource](virtual_machine.html) instead.

~> In this release there's a known issue where the `public_ip_address` and `public_ip_addresses` fields may not be fully populated for Dynamic Public IP's.

## Example Usage

This example provisions a basic Linux Virtual Machine on an internal network. Additional examples of how to use the `azurerm_linux_virtual_machine` resource can be found [in the ./examples/virtual-machines/linux directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/virtual-machines/linux).

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
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
```

## Argument Reference

The following arguments are supported:

* `admin_username` - (Required) The username of the local administrator used for the Virtual Machine. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Linux Virtual Machine should exist. Changing this forces a new resource to be created.

* `license_type` - (Optional) Specifies the BYOL Type for this Virtual Machine. Possible values are `RHEL_BYOS` and `SLES_BYOS`.

* `name` - (Required) The name of the Linux Virtual Machine. Changing this forces a new resource to be created.

* `network_interface_ids` - (Required). A list of Network Interface IDs which should be attached to this Virtual Machine. The first Network Interface ID in this list will be the Primary Network Interface on the Virtual Machine.

* `os_disk` - (Required) A `os_disk` block as defined below.

* `resource_group_name` - (Required) The name of the Resource Group in which the Linux Virtual Machine should be exist. Changing this forces a new resource to be created.

* `size` - (Required) The SKU which should be used for this Virtual Machine, such as `Standard_F2`.

---

* `additional_capabilities` - (Optional) A `additional_capabilities` block as defined below.

* `admin_password` - (Optional) The Password which should be used for the local-administrator on this Virtual Machine. Changing this forces a new resource to be created.

-> **NOTE:** When an `admin_password` is specified `disable_password_authentication` must be set to `false`.
~> **NOTE:** One of either `admin_password` or `admin_ssh_key` must be specified.

* `admin_ssh_key` - (Optional) One or more `admin_ssh_key` blocks as defined below.

~> **NOTE:** One of either `admin_password` or `admin_ssh_key` must be specified.

* `allow_extension_operations` - (Optional) Should Extension Operations be allowed on this Virtual Machine?

* `availability_set_id` - (Optional) Specifies the ID of the Availability Set in which the Virtual Machine should exist. Changing this forces a new resource to be created.

* `boot_diagnostics` - (Optional) A `boot_diagnostics` block as defined below.

* `capacity_reservation_group_id` - (Optional) Specifies the ID of the Capacity Reservation Group which the Virtual Machine should be allocated to.

~> **NOTE:** `capacity_reservation_group_id` cannot be used with `availability_set_id` or `proximity_placement_group_id`

* `computer_name` - (Optional) Specifies the Hostname which should be used for this Virtual Machine. If unspecified this defaults to the value for the `name` field. If the value of the `name` field is not a valid `computer_name`, then you must specify `computer_name`. Changing this forces a new resource to be created.

* `custom_data` - (Optional) The Base64-Encoded Custom Data which should be used for this Virtual Machine. Changing this forces a new resource to be created.

* `dedicated_host_id` - (Optional) The ID of a Dedicated Host where this machine should be run on. Conflicts with `dedicated_host_group_id`.

* `dedicated_host_group_id` - (Optional) The ID of a Dedicated Host Group that this Linux Virtual Machine should be run within. Conflicts with `dedicated_host_id`.

* `disable_password_authentication` - (Optional) Should Password Authentication be disabled on this Virtual Machine? Defaults to `true`. Changing this forces a new resource to be created.

-> In general we'd recommend using SSH Keys for authentication rather than Passwords - but there's tradeoff's to each - please [see this thread for more information](https://security.stackexchange.com/questions/69407/why-is-using-an-ssh-key-more-secure-than-using-passwords).

-> **NOTE:** When an `admin_password` is specified `disable_password_authentication` must be set to `false`.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Linux Virtual Machine should exist. Changing this forces a new Linux Virtual Machine to be created.

* `encryption_at_host_enabled` - (Optional) Should all of the disks (including the temp disk) attached to this Virtual Machine be encrypted by enabling Encryption at Host?

* `eviction_policy` - (Optional) Specifies what should happen when the Virtual Machine is evicted for price reasons when using a Spot instance. At this time the only supported value is `Deallocate`. Changing this forces a new resource to be created.

-> **NOTE:** This can only be configured when `priority` is set to `Spot`.

* `extensions_time_budget` - (Optional) Specifies the duration allocated for all extensions to start. The time duration should be between 15 minutes and 120 minutes (inclusive) and should be specified in ISO 8601 format. Defaults to 90 minutes (`PT1H30M`).

* `identity` - (Optional) An `identity` block as defined below.

* `patch_mode` - (Optional) Specifies the mode of in-guest patching to this Linux Virtual Machine. Possible values are `AutomaticByPlatform` and `ImageDefault`. Defaults to `ImageDefault`. For more information on patch modes please see the [product documentation](https://docs.microsoft.com/azure/virtual-machines/automatic-vm-guest-patching#patch-orchestration-modes).

-> **NOTE:** If `patch_mode` is set to `AutomaticByPlatform` then `provision_vm_agent` must also be set to `true`.

* `max_bid_price` - (Optional) The maximum price you're willing to pay for this Virtual Machine, in US Dollars; which must be greater than the current spot price. If this bid price falls below the current spot price the Virtual Machine will be evicted using the `eviction_policy`. Defaults to `-1`, which means that the Virtual Machine should not be evicted for price reasons.

-> **NOTE:** This can only be configured when `priority` is set to `Spot`.

* `plan` - (Optional) A `plan` block as defined below. Changing this forces a new resource to be created.

* `platform_fault_domain` - (Optional) Specifies the Platform Fault Domain in which this Linux Virtual Machine should be created. Defaults to `-1`, which means this will be automatically assigned to a fault domain that best maintains balance across the available fault domains. Changing this forces a new Linux Virtual Machine to be created.

* `priority`- (Optional) Specifies the priority of this Virtual Machine. Possible values are `Regular` and `Spot`. Defaults to `Regular`. Changing this forces a new resource to be created.

* `provision_vm_agent` - (Optional) Should the Azure VM Agent be provisioned on this Virtual Machine? Defaults to `true`. Changing this forces a new resource to be created.

~> **NOTE:** If `provision_vm_agent` is set to `false` then `allow_extension_operations` must also be set to `false`.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group which the Virtual Machine should be assigned to.

* `secret` - (Optional) One or more `secret` blocks as defined below.

* `secure_boot_enabled` - (Optional) Specifies whether secure boot should be enabled on the virtual machine. Changing this forces a new resource to be created.

* `source_image_id` - (Optional) The ID of the Image which this Virtual Machine should be created from. Changing this forces a new resource to be created.

-> **NOTE:** One of either `source_image_id` or `source_image_reference` must be set.

* `source_image_reference` - (Optional) A `source_image_reference` block as defined below. Changing this forces a new resource to be created.

-> **NOTE:** One of either `source_image_id` or `source_image_reference` must be set.

* `tags` - (Optional) A mapping of tags which should be assigned to this Virtual Machine.

* `termination_notification` - (Optional) A `termination_notification` block as defined below.

* `user_data` - (Optional) The Base64-Encoded User Data which should be used for this Virtual Machine.

* `vtpm_enabled` - (Optional) Specifies whether vTPM should be enabled on the virtual machine. Changing this forces a new resource to be created.

* `virtual_machine_scale_set_id` - (Optional) Specifies the Orchestrated Virtual Machine Scale Set that this Virtual Machine should be created within. Changing this forces a new resource to be created.

~> **NOTE:** Orchestrated Virtual Machine Scale Sets can be provisioned using [the `azurerm_orchestrated_virtual_machine_scale_set` resource](/docs/providers/azurerm/r/orchestrated_virtual_machine_scale_set.html).

* `zone` - (Optional) Specifies the Availability Zones in which this Linux Virtual Machine should be located. Changing this forces a new Linux Virtual Machine to be created.

---

A `additional_capabilities` block supports the following:

* `ultra_ssd_enabled` - (Optional) Should the capacity to enable Data Disks of the `UltraSSD_LRS` storage account type be supported on this Virtual Machine? Defaults to `false`.

---

A `admin_ssh_key` block supports the following:

* `public_key` - (Required) The Public Key which should be used for authentication, which needs to be at least 2048-bit and in `ssh-rsa` format. Changing this forces a new resource to be created.

* `username` - (Required) The Username for which this Public SSH Key should be configured. Changing this forces a new resource to be created.

-> **NOTE:** The Azure VM Agent only allows creating SSH Keys at the path `/home/{username}/.ssh/authorized_keys` - as such this public key will be written to the authorized keys file.

---

A `boot_diagnostics` block supports the following:

* `storage_account_uri` - (Optional) The Primary/Secondary Endpoint for the Azure Storage Account which should be used to store Boot Diagnostics, including Console Output and Screenshots from the Hypervisor.

-> **NOTE:** Passing a null value will utilize a Managed Storage Account to store Boot Diagnostics

---

A `certificate` block supports the following:

* `url` - (Required) The Secret URL of a Key Vault Certificate.

-> **NOTE:** This can be sourced from the `secret_id` field within the `azurerm_key_vault_certificate` Resource.

---

A `diff_disk_settings` block supports the following:

* `option` - (Required) Specifies the Ephemeral Disk Settings for the OS Disk. At this time the only possible value is `Local`. Changing this forces a new resource to be created.

* `placement` - (Optional) Specifies where to store the Ephemeral Disk. Possible values are `CacheDisk` and `ResourceDisk`. Defaults to `CacheDisk`. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Linux Virtual Machine. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Linux Virtual Machine.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `os_disk` block supports the following:

* `caching` - (Required) The Type of Caching which should be used for the Internal OS Disk. Possible values are `None`, `ReadOnly` and `ReadWrite`.

* `storage_account_type` - (Required) The Type of Storage Account which should back this the Internal OS Disk. Possible values are `Standard_LRS`, `StandardSSD_LRS`, `Premium_LRS`, `StandardSSD_ZRS` and `Premium_ZRS`. Changing this forces a new resource to be created.

* `diff_disk_settings` (Optional) A `diff_disk_settings` block as defined above. Changing this forces a new resource to be created.

-> **NOTE:** `diff_disk_settings` can only be set when `caching` is set to `ReadOnly`. More information can be found [here](https://docs.microsoft.com/azure/virtual-machines/ephemeral-os-disks-deploy#vm-template-deployment)

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to Encrypt this OS Disk. Conflicts with `secure_vm_disk_encryption_set_id`.

-> **NOTE:** The Disk Encryption Set must have the `Reader` Role Assignment scoped on the Key Vault - in addition to an Access Policy to the Key Vault

* `disk_size_gb` - (Optional) The Size of the Internal OS Disk in GB, if you wish to vary from the size used in the image this Virtual Machine is sourced from.

-> **NOTE:** If specified this must be equal to or larger than the size of the Image the Virtual Machine is based on. When creating a larger disk than exists in the image you'll need to repartition the disk to use the remaining space.

* `name` - (Optional) The name which should be used for the Internal OS Disk. Changing this forces a new resource to be created.

* `secure_vm_disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to Encrypt this OS Disk when the Virtual Machine is a Confidential VM. Conflicts with `disk_encryption_set_id`. Changing this forces a new resource to be created.

~> **NOTE:** `secure_vm_disk_encryption_set_id` can only be specified when `security_encryption_type` is set to `DiskWithVMGuestState`.

* `security_encryption_type` - (Optional) Encryption Type when the Virtual Machine is a Confidential VM. Possible values are `VMGuestStateOnly` and `DiskWithVMGuestState`. Changing this forces a new resource to be created.

~> **NOTE:** `secure_boot_enabled` and `vtpm_enabled` must be set to `true` when `security_encryption_type` is specified.

~> **NOTE:** `encryption_at_host_enabled` cannot be set to `true` when `security_encryption_type` is set to `DiskWithVMGuestState`.

* `write_accelerator_enabled` - (Optional) Should Write Accelerator be Enabled for this OS Disk? Defaults to `false`.

-> **NOTE:** This requires that the `storage_account_type` is set to `Premium_LRS` and that `caching` is set to `None`.

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

`source_image_reference` supports the following:

* `publisher` - (Optional) Specifies the publisher of the image used to create the virtual machines.

* `offer` - (Optional) Specifies the offer of the image used to create the virtual machines.

* `sku` - (Optional) Specifies the SKU of the image used to create the virtual machines.

* `version` - (Optional) Specifies the version of the image used to create the virtual machines.

---

A `termination_notification` block supports the following:

* `enabled` - (Required) Should the termination notification be enabled on this Virtual Machine? Defaults to `false`.

* `timeout` - (Optional) Length of time (in minutes, between `5` and `15`) a notification to be sent to the VM on the instance metadata server till the VM gets deleted. The time duration should be specified in ISO 8601 format.

~> **NOTE:** For more information about the termination notification, please [refer to this doc](https://docs.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-terminate-notification).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Linux Virtual Machine.

* `identity` - An `identity` block as documented below.

* `private_ip_address` - The Primary Private IP Address assigned to this Virtual Machine.

* `private_ip_addresses` - A list of Private IP Addresses assigned to this Virtual Machine.

* `public_ip_address` - The Primary Public IP Address assigned to this Virtual Machine.

* `public_ip_addresses` - A list of the Public IP Addresses assigned to this Virtual Machine.

* `virtual_machine_id` - A 128-bit identifier which uniquely identifies this Virtual Machine.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 45 minutes) Used when creating the Linux Virtual Machine.
* `update` - (Defaults to 45 minutes) Used when updating the Linux Virtual Machine.
* `delete` - (Defaults to 45 minutes) Used when deleting the Linux Virtual Machine.

## Import

Linux Virtual Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_linux_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1
```
