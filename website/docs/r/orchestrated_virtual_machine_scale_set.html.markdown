---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orchestrated_virtual_machine_scale_set"
description: |-
  Manages an Orchestrated Virtual Machine Scale Set.
---

# azurerm_orchestrated_virtual_machine_scale_set

Manages an Orchestrated Virtual Machine Scale Set.

## Disclaimers

~> **NOTE:** As of the **v2.86.0** (November 19, 2021) release of the provider this resource will only create Virtual Machine Scale Sets with the **Flexible** Orchestration Mode.

~> **NOTE:** All arguments including the administrator login and password will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "example" {
  name                = "example-VMSS"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  platform_fault_domain_count = 1

  zones = ["1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Orchestrated Virtual Machine Scale Set. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Orchestrated Virtual Machine Scale Set should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Orchestrated Virtual Machine Scale Set should exist. Changing this forces a new resource to be created.

* `platform_fault_domain_count` - (Required) Specifies the number of fault domains that are used by this Orchestrated Virtual Machine Scale Set. Changing this forces a new resource to be created.

~> **NOTE:** The number of Fault Domains varies depending on which Azure Region you're using - a list can be found [here](https://github.com/MicrosoftDocs/azure-docs/blob/master/includes/managed-disks-common-fault-domain-region-list.md).

* `sku_name` - (Optional) The `name` of the SKU to be used by this Orcestrated Virtual Machine Scale Set. Valid values include: any of the [General purpose](https://docs.microsoft.com/azure/virtual-machines/sizes-general), [Compute optimized](https://docs.microsoft.com/azure/virtual-machines/sizes-compute), [Memory optimized](https://docs.microsoft.com/azure/virtual-machines/sizes-memory), [Storage optimized](https://docs.microsoft.com/azure/virtual-machines/sizes-storage), [GPU optimized](https://docs.microsoft.com/azure/virtual-machines/sizes-gpu), [FPGA optimized](https://docs.microsoft.com/azure/virtual-machines/sizes-field-programmable-gate-arrays), [High performance](https://docs.microsoft.com/azure/virtual-machines/sizes-hpc), or [Previous generation](https://docs.microsoft.com/azure/virtual-machines/sizes-previous-gen) virtual machine SKUs.

* `instances`- (Optional) The number of Virtual Machines in the Orcestrated Virtual Machine Scale Set.

* `network_interface` - (Optional) One or more `network_interface` blocks as defined below.

* `os_profile` - (Optional) An `os_profile` block as defined below.

* `os_disk` - (Optional) An `os_disk` block as defined below.

* `boot_diagnostics` - (Optional) A `boot_diagnostics` block as defined below.

* `computer_name_prefix` - (Optional) The prefix which should be used for the name of the Virtual Machines in this Scale Set. If unspecified this defaults to the value for the name field. If the value of the name field is not a valid `computer_name_prefix`, then you must specify `computer_name_prefix`.

* `data_disk` - (Optional) One or more `data_disk` blocks as defined below.

* `extension` - (Optional) One or more `extension` blocks as defined below

* `extensions_time_budget` - (Optional) Specifies the time alloted for all extensions to start. The time duration should be between 15 minutes and 120 minutes (inclusive) and should be specified in ISO 8601 format. The default value is 90 minutes (PT1H30M).

* `eviction_policy` - (Optional) The Policy which should be used Virtual Machines are Evicted from the Scale Set. Changing this forces a new resource to be created.

* `identity` - (Optional) A `identity` block as defined below.

* `license_type` - (Optional) Specifies the type of on-premise license (also known as Azure Hybrid Use Benefit) which should be used for this Orchestrated Virtual Machine Scale Set. Possible values are `None`, `Windows_Client` and `Windows_Server`. Changing this forces a new resource to be created.

* `max_bid_price` - (Optional) The maximum price you're willing to pay for each Orchestrated Virtual Machine in this Scale Set, in US Dollars; which must be greater than the current spot price. If this bid price falls below the current spot price the Virtual Machines in the Scale Set will be evicted using the eviction_policy. Defaults to `-1`, which means that each Virtual Machine in the Orchestrated Scale Set should not be evicted for price reasons.

* `plan` - (Optional) A `plan` block as documented below.

* `priority` - (Optional) The Priority of this Orchestrated Virtual Machine Scale Set. Possible values are `Regular` and `Spot`. Defaults to `Regular`. Changing this value forces a new resource.

* `source_image_id` - (Optional) The ID of an Image which each Virtual Machine in this Scale Set should be based on.

* `source_image_reference` - (Optional) A `source_image_reference` block as defined below.

* `termination_notification` - (Optional) A `termination_notification` block as defined below.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group which the Orchestrated Virtual Machine should be assigned to. Changing this forces a new resource to be created.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Orchestrated Virtual Machine should be located. Changing this forces a new Orchestrated Virtual Machine to be created.

~> **NOTE:** Due to a limitation of the Azure API at this time only one Availability Zone can be defined.

* `tags` - (Optional) A mapping of tags which should be assigned to this Orchestrated Virtual Machine Scale Set.

---

A `os_profile` block supports the following:

* `custom_data` - (Optional) The Base64-Encoded Custom Data which should be used for this Orchestrated Virtual Machine Scale Set.

~> **NOTE:** When Custom Data has been configured, it's not possible to remove it without tainting the Orchestrated Virtual Machine Scale Set, due to a limitation of the Azure API.

* `windows_configuration` - (Optional) A `windows_configuration` block as documented below.

* `linux_configuration` - (Optional) A `linux_configuration` block as documented below.

---

A `windows_configuration` block supports the following:

* `admin_username` - (Required) The username of the local administrator on each Orchestrated Virtual Machine Scale Set instance. Changing this forces a new resource to be created.

* `admin_password` - (Required) The Password which should be used for the local-administrator on this Virtual Machine. Changing this forces a new resource to be created.

* `computer_name_prefix` - (Optional) The prefix which should be used for the name of the Virtual Machines in this Scale Set. If unspecified this defaults to the value for the `name` field. If the value of the `name` field is not a valid `computer_name_prefix`, then you must specify `computer_name_prefix`.

* `enable_automatic_updates` - (Optional) Are automatic updates enabled for this Virtual Machine? Defaults to `true`.

* `hotpatching_enabled` - (Optional) Should the VM be patched without requiring a reboot? Possible values are `true` or `false`. Defaults to `false`. For more information about hot patching please see the [product documentation](https://docs.microsoft.com/azure/automanage/automanage-hotpatch).

-> **NOTE:** Hotpatching can only be enabled if the `patch_mode` is set to `AutomaticByPlatform`, the `provision_vm_agent` is set to `true`, your `source_image_reference` references a hotpatching enabled image, the VM's `sku_name` is set to a [Azure generation 2](https://docs.microsoft.com/azure/virtual-machines/generation-2#generation-2-vm-sizes) VM SKU and the `extension` contains an application health extension. An example of how to correctly configure a Orchestrated Virtual Machine Scale Set to provision a Windows Virtual Machine with hotpatching enabled can be found in the [`./examples/orchestrated-vm-scale-set/hotpatching-enabled`](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/orchestrated-vm-scale-set/hotpatching-enabled) directory within the GitHub Repository.

* `patch_mode` - (Optional) Specifies the mode of in-guest patching of this Windows Virtual Machine. Possible values are `Manual`, `AutomaticByOS` and `AutomaticByPlatform`. Defaults to `AutomaticByOS`. For more information on patch modes please see the [product documentation](https://docs.microsoft.com/azure/virtual-machines/automatic-vm-guest-patching#patch-orchestration-modes).

-> **NOTE:** If `patch_mode` is set to `AutomaticByPlatform` the `provision_vm_agent` must be set to `true` and the `extension` must contain at least one application health extension.

* `provision_vm_agent` - (Optional) Should the Azure VM Agent be provisioned on each Virtual Machine in the Scale Set? Defaults to `true`. Changing this value forces a new resource to be created.

* `secret` - (Optional) One or more `secret` blocks as defined below.

* `timezone` - (Optional) Specifies the time zone of the virtual machine, the possible values are defined [here](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/).

* `winrm_listener` - (Optional) One or more `winrm_listener` blocks as defined below.

---

A `linux_configuration` block supports the following:

* `admin_username` - (Required) The username of the local administrator on each Orchestrated Virtual Machine Scale Set instance. Changing this forces a new resource to be created.

* `admin_password` - (Optional) The Password which should be used for the local-administrator on this Virtual Machine. Changing this forces a new resource to be created.

* `admin_ssh_key` - (Optional) A `admin_ssh_key` block as documented below.

* `disable_password_authentication` - (Optional) When an `admin_password` is specified `disable_password_authentication` must be set to `false`. Defaults to `true`.

~> **NOTE:** Either `admin_password` or `admin_ssh_key` must be specified.

* `patch_mode` - (Optional) Specifies the mode of in-guest patching of this Windows Virtual Machine. Possible values are `ImageDefault` or `AutomaticByPlatform`. Defaults to `ImageDefault`. For more information on patch modes please see the [product documentation](https://docs.microsoft.com/azure/virtual-machines/automatic-vm-guest-patching#patch-orchestration-modes).

-> **NOTE:** If `patch_mode` is set to `AutomaticByPlatform` the `provision_vm_agent` must be set to `true` and the `extension` must contain at least one application health extension.  An example of how to correctly configure a Orchestrated Virtual Machine Scale Set to provision a Linux Virtual Machine with Automatic VM Guest Patching enabled can be found in the [`./examples/orchestrated-vm-scale-set/automatic-vm-guest-patching`](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/orchestrated-vm-scale-set/automatic-vm-guest-patching) directory within the GitHub Repository.

* `provision_vm_agent` - (Optional) Should the Azure VM Agent be provisioned on each Virtual Machine in the Scale Set? Defaults to `true`. Changing this value forces a new resource to be created.

* `secret` - (Optional) One or more `secret` blocks as defined below.

---

A `secret` block supports the following:

* `key_vault_id` - (Required) The ID of the Key Vault from which all Secrets should be sourced.

* `certificate` - (Required) One or more `certificate` blocks as defined below.

~> **NOTE:** The schema of the `certificate` block is slightly different depending on if you are provisioning a `windows_configuration` or a `linux_configuration`.

---

A (Windows) `certificate` block supports the following:

* `store` - (Required) The certificate store on the Virtual Machine where the certificate should be added.

* `url` - (Required) The Secret URL of a Key Vault Certificate.

---

A (Linux) `certificate` block supports the following:

* `url` - (Required) The Secret URL of a Key Vault Certificate.

---

A `admin_ssh_key` block supports the following:

* `public_key` - (Required) The Public Key which should be used for authentication, which needs to be at least 2048-bit and in ssh-rsa format. Changing this forces a new resource to be created.

* `username` - (Required) The Username for which this Public SSH Key should be configured. Changing this forces a new resource to be created.

~> **NOTE:** The Azure VM Agent only allows creating SSH Keys at the path `/home/{username}/.ssh/authorized_keys` - as such this public key will be written to the authorized keys file.

---

A `winrm_listener` block supports the following:

* `certificate_url` - (Optional) The Secret URL of a Key Vault Certificate, which must be specified when protocol is set to `Https`.

~> **NOTE:** This can be sourced from the `secret_id` field within the `azurerm_key_vault_certificate` Resource.

---

A `automatic_instance_repair` block supports the following:

* `enabled` - (Required) Should the automatic instance repair be enabled on this Orchestrated Virtual Machine Scale Set? Possible values are `true` and `false`. Defaults to `false`.

* `grace_period` - (Optional) Amount of time for which automatic repairs will be delayed. The grace period starts right after the VM is found unhealthy. Possible values are between `30` and `90` minutes. Defaults to `30` minutes. The time duration should be specified in `ISO 8601` format (e.g. `PT30M` to `PT90M`).

---

A `boot_diagnostics` block supports the following:

* `storage_account_uri` - (Optional) The Primary/Secondary Endpoint for the Azure Storage Account which should be used to store Boot Diagnostics, including Console Output and Screenshots from the Hypervisor. By including a `boot_diagnostics` block without passing the `storage_account_uri` field will cause the API to utilize a Managed Storage Account to store the Boot Diagnostics output.

---

A `certificate` block supports the following:

* `store` - (Required) The certificate store on the Virtual Machine where the certificate should be added.

* `url` - (Required) The Secret URL of a Key Vault Certificate.

~> **NOTE:** This can be sourced from the `secret_id` field within the `azurerm_key_vault_certificate` Resource.

---

A `diff_disk_settings` block supports the following:

* `option` - (Required) Specifies the Ephemeral Disk Settings for the OS Disk. At this time the only possible value is `Local`. Changing this forces a new resource to be created.

* `placement` - (Optional) Specifies where to store the Ephemeral Disk. Possible values are `CacheDisk` and `ResourceDisk`. Defaults to `CacheDisk`. Changing this forces a new resource to be created.

---

A `data_disk` block supports the following:

* `caching` - (Required) The type of Caching which should be used for this Data Disk. Possible values are None, ReadOnly and ReadWrite.

* `create_option` - (Optional) The create option which should be used for this Data Disk. Possible values are Empty and FromImage. Defaults to Empty. (FromImage should only be used if the source image includes data disks).

* `disk_size_gb` - (Required) The size of the Data Disk which should be created.

* `lun` - (Required) The Logical Unit Number of the Data Disk, which must be unique within the Virtual Machine.

* `storage_account_type` - (Required) The Type of Storage Account which should back this Data Disk. Possible values include Standard_LRS, StandardSSD_LRS, Premium_LRS and UltraSSD_LRS.

---

An `extension` block supports the following:

* `name` - (Required) The name for the Virtual Machine Scale Set Extension.

* `publisher` - (Required) Specifies the Publisher of the Extension.

* `type` - (Required) Specifies the Type of the Extension.

* `type_handler_version` - (Required) Specifies the version of the extension to use, available versions can be found using the Azure CLI.

* `auto_upgrade_minor_version_enabled` - (Optional) Should the latest version of the Extension be used at Deployment Time, if one is available? This won't auto-update the extension on existing installation. Defaults to true.

* `extensions_to_provision_after_vm_creation` - (Optional) An ordered list of Extension names which Orchestrated Virtual Machine Scale Set should provision after VM creation.

* `force_extension_execution_on_change` - (Optional) A value which, when different to the previous value can be used to force-run the Extension even if the Extension Configuration hasn't changed.

* `protected_settings` - (Optional) A JSON String which specifies Sensitive Settings (such as Passwords) for the Extension.

~> **NOTE:** Keys within the `protected_settings` block are notoriously case-sensitive, where the casing required (e.g. `TitleCase` vs `snakeCase`) depends on the Extension being used. Please refer to the documentation for the specific Orchestrated Virtual Machine Extension you're looking to use for more information.

---

A `ip_configuration` block supports the following:

* `name` - (Required) The Name which should be used for this IP Configuration.

* `application_gateway_backend_address_pool_ids` - (Optional) A list of Backend Address Pools IDs from a Application Gateway which this Orchestrated Virtual Machine Scale Set should be connected to.

* `application_security_group_ids` - (Optional) A list of Application Security Group IDs which this Orchestrated Virtual Machine Scale Set should be connected to.

* `load_balancer_backend_address_pool_ids` - (Optional) A list of Backend Address Pools IDs from a Load Balancer which this Orchestrated Virtual Machine Scale Set should be connected to.

~> **NOTE:** When using this field you'll also need to configure a Rule for the Load Balancer, and use a depends_on between this resource and the Load Balancer Rule.

* `primary` - (Optional) Is this the Primary IP Configuration for this Network Interface? Possible values are `true` and `false`. Defaults to `false`.

~> **NOTE:** One `ip_configuration` block must be marked as Primary for each Network Interface.

* `public_ip_address` - (Optional) A `public_ip_address` block as defined below.

* `subnet_id` - (Optional) The ID of the Subnet which this IP Configuration should be connected to.

~> **NOTE:** `subnet_id` is required if version is set to `IPv4`.

* `version` - (Optional) The Internet Protocol Version which should be used for this IP Configuration. Possible values are `IPv4` and `IPv6`. Defaults to `IPv4`.

---

A `ip_tag` block supports the following:

* `tag` - The IP Tag associated with the Public IP, such as `SQL` or `Storage`.

* `type` - The Type of IP Tag, such as `FirstPartyUsage`.

---

A `network_interface` block supports the following:

* `name` - (Required) The Name which should be used for this Network Interface. Changing this forces a new resource to be created.

* `ip_configuration` - (Required) One or more `ip_configuration` blocks as defined above.

* `dns_servers` - (Optional) A list of IP Addresses of DNS Servers which should be assigned to the Network Interface.

* `enable_accelerated_networking` - (Optional) Does this Network Interface support Accelerated Networking? Possible values are `true` and `false`. Defaults to `false`.

* `enable_ip_forwarding` - (Optional) Does this Network Interface support IP Forwarding? Possible values are `true` and `false`. Defaults to `false`.

* `network_security_group_id` - (Optional) The ID of a Network Security Group which should be assigned to this Network Interface.

* `primary` - (Optional) Is this the Primary IP Configuration? Possible values are `true` and `false`. Defaults to `false`.

~> **NOTE:** If multiple `network_interface` blocks are specified, one must be set to `primary`.

---

A `os_disk` block supports the following:

* `caching` - (Required) The Type of Caching which should be used for the Internal OS Disk. Possible values are `None`, `ReadOnly` and `ReadWrite`.

* `storage_account_type` - (Required) The Type of Storage Account which should back this the Internal OS Disk. Possible values include `Standard_LRS`, `StandardSSD_LRS` and `Premium_LRS`.

* `diff_disk_settings` - (Optional) A `diff_disk_settings` block as defined above.

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to encrypt this OS Disk.

~> **NOTE:** Disk Encryption Sets are in Public Preview in a limited set of regions

* `disk_size_gb` - (Optional) The Size of the Internal OS Disk in GB, if you wish to vary from the size used in the image this Virtual Machine Scale Set is sourced from.

---

A `plan` block supports the following:

* `name` - (Required) Specifies the name of the image from the marketplace. Changing this forces a new resource to be created.

* `publisher` - (Required) Specifies the publisher of the image. Changing this forces a new resource to be created.

* `product` - (Required) Specifies the product of the image from the marketplace. Changing this forces a new resource to be created.

---

A `identity` block supports the following:

* `type` - (Required) The type of Managed Identity that should be configured on this Orchestrated Windows Virtual Machine Scale Set. Only possible value is `UserAssigned`.

* `identity_ids` - (Required) Specifies a list of User Managed Identity IDs to be assigned to this Orchestrated Windows Virtual Machine Scale Set.

---

A `public_ip_address` block supports the following:

* `name` - (Required) The Name of the Public IP Address Configuration.

* `domain_name_label` - (Optional) The Prefix which should be used for the Domain Name Label for each Virtual Machine Instance. Azure concatenates the Domain Name Label and Virtual Machine Index to create a unique Domain Name Label for each Virtual Machine. Valid values must be between `1` and `26` characters long, start with a lower case letter, end with a lower case letter or number and contains only `a-z`, `0-9` and `hyphens`.

* `idle_timeout_in_minutes` - (Optional) The Idle Timeout in Minutes for the Public IP Address. Possible values are in the range `4` to `32`.

* `ip_tag` - (Optional) One or more `ip_tag` blocks as defined above.

* `public_ip_prefix_id` - (Optional) The ID of the Public IP Address Prefix from where Public IP Addresses should be allocated. Changing this forces a new resource to be created.

---

A `termination_notification` block supports the following:

* `enabled` - (Required) Should the termination notification be enabled on this Virtual Machine Scale Set? Possible values `true` or `false` Defaults to `false`.

* `timeout` - (Optional) Length of time (in minutes, between `5` and `15`) a notification to be sent to the VM on the instance metadata server till the VM gets deleted. The time duration should be specified in `ISO 8601` format. Defaults to `PT5M`.

---

A `source_image_reference` block supports the following:

* `publisher` - (Optional) Specifies the publisher of the image used to create the virtual machines.

* `offer` - (Optional) Specifies the offer of the image used to create the virtual machines.

* `sku` - (Optional) Specifies the SKU of the image used to create the virtual machines.

* `version` - (Optional) Specifies the version of the image used to create the virtual machines.

---

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Orchestrated Virtual Machine Scale Set.

* `unique_id` - The Unique ID for the Orchestrated Virtual Machine Scale Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Orchestrated Virtual Machine Scale Set.
* `update` - (Defaults to 30 minutes) Used when updating the Orchestrated Virtual Machine Scale Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Orchestrated Virtual Machine Scale Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the Orchestrated Virtual Machine Scale Set.

## Import

An Orchestrated Virtual Machine Scale Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_orchestrated_virtual_machine_scale_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1
```
