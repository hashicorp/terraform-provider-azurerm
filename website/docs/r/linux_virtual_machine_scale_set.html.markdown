---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_linux_virtual_machine_scale_set"
description: |-
  Manages a Linux Virtual Machine Scale Set.
---

# azurerm_linux_virtual_machine_scale_set

Manages a Linux Virtual Machine Scale Set.

## Disclaimers

~> **Note**: All arguments including the administrator login and password will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

-> **Note** Terraform will automatically update & reimage the nodes in the Scale Set (if Required) during an Update - this behaviour can be configured [using the `features` setting within the Provider block](https://www.terraform.io/docs/providers/azurerm/index.html#features).

~> **Note:** This resource does not support Unmanaged Disks. If you need to use Unmanaged Disks you can continue to use [the `azurerm_virtual_machine_scale_set` resource](virtual_machine_scale_set.html) instead

## Example Usage

This example provisions a basic Linux Virtual Machine Scale Set on an internal network. Additional examples of how to use the `azurerm_linux_virtual_machine_scale_set` resource can be found [in the ./examples/vm-scale-set/linux` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/vm-scale-set/linux).

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_linux_virtual_machine_scale_set" "example" {
  name                = "example-vmss"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = file("~/.ssh/id_rsa.pub")
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.internal.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Linux Virtual Machine Scale Set. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Linux Virtual Machine Scale Set should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Linux Virtual Machine Scale Set should be exist. Changing this forces a new resource to be created.

* `admin_username` - (Required) The username of the local administrator on each Virtual Machine Scale Set instance. Changing this forces a new resource to be created.

* `instances` - (Required) The number of Virtual Machines in the Scale Set.

-> **NOTE:** If you're using AutoScaling, you may wish to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/configuration/resources.html#ignore_changes) to ignore changes to this field.

* `sku` - (Required) The Virtual Machine SKU for the Scale Set, such as `Standard_F2`.

* `network_interface` - (Required) One or more `network_interface` blocks as defined below.

* `os_disk` - (Required) An `os_disk` block as defined below.

---

* `additional_capabilities` - (Optional) A `additional_capabilities` block as defined below.

* `admin_password` - (Optional) The Password which should be used for the local-administrator on this Virtual Machine. Changing this forces a new resource to be created.

-> **NOTE:** When an `admin_password` is specified `disable_password_authentication` must be set to `false`.

~> **NOTE:** One of either `admin_password` or `admin_ssh_key` must be specified.

* `admin_ssh_key` - (Optional) One or more `admin_ssh_key` blocks as defined below.

~> **NOTE:** One of either `admin_password` or `admin_ssh_key` must be specified.

* `automatic_os_upgrade_policy` - (Optional) A `automatic_os_upgrade_policy` block as defined below. This can only be specified when `upgrade_mode` is set to `Automatic`.

* `automatic_instance_repair` - (Optional) A `automatic_instance_repair` block as defined below. To enable the automatic instance repair, this Virtual Machine Scale Set must have a valid `health_probe_id` or an [Application Health Extension](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-health-extension).

~> **NOTE:** For more information about Automatic Instance Repair, please refer to [this doc](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-automatic-instance-repairs).

* `boot_diagnostics` - (Optional) A `boot_diagnostics` block as defined below.

* `computer_name_prefix` - (Optional) The prefix which should be used for the name of the Virtual Machines in this Scale Set. If unspecified this defaults to the value for the `name` field. If the value of the `name` field is not a valid `computer_name_prefix`, then you must specify `computer_name_prefix`.

* `custom_data` - (Optional) The Base64-Encoded Custom Data which should be used for this Virtual Machine Scale Set.

-> **NOTE:** When Custom Data has been configured, it's not possible to remove it without tainting the Virtual Machine Scale Set, due to a limitation of the Azure API.

* `data_disk` - (Optional) One or more `data_disk` blocks as defined below.

* `disable_password_authentication` - Should Password Authentication be disabled on this Virtual Machine Scale Set? Defaults to `true`.

-> In general we'd recommend using SSH Keys for authentication rather than Passwords - but there's tradeoff's to each - please [see this thread for more information](https://security.stackexchange.com/questions/69407/why-is-using-an-ssh-key-more-secure-than-using-passwords).

-> **NOTE:** When an `admin_password` is specified `disable_password_authentication` must be set to `false`.

* `do_not_run_extensions_on_overprovisioned_machines` - (Optional) Should Virtual Machine Extensions be run on Overprovisioned Virtual Machines in the Scale Set? Defaults to `false`.

* `encryption_at_host_enabled` - (Optional) Should all of the disks (including the temp disk) attached to this Virtual Machine be encrypted by enabling Encryption at Host?

* `extension` - (Optional) One or more `extension` blocks as defined below

!> **NOTE:** This block is only available in the Opt-In beta and requires that the Environment Variable `ARM_PROVIDER_VMSS_EXTENSIONS_BETA` is set to `true` to be used.

* `extensions_time_budget` - (Optional) Specifies the duration allocated for all extensions to start. The time duration should be between `15` minutes and `120` minutes (inclusive) and should be specified in ISO 8601 format. Defaults to `90` minutes (`PT1H30M`).

* `eviction_policy` - (Optional) The Policy which should be used Virtual Machines are Evicted from the Scale Set. Changing this forces a new resource to be created.

-> **NOTE:** This can only be configured when `priority` is set to `Spot`.

* `health_probe_id` - (Optional) The ID of a Load Balancer Probe which should be used to determine the health of an instance. Changing this forces a new resource to be created. This is Required and can only be specified when `upgrade_mode` is set to `Automatic` or `Rolling`.

* `identity` - (Optional) A `identity` block as defined below.

* `max_bid_price` - (Optional) The maximum price you're willing to pay for each Virtual Machine in this Scale Set, in US Dollars; which must be greater than the current spot price. If this bid price falls below the current spot price the Virtual Machines in the Scale Set will be evicted using the `eviction_policy`. Defaults to `-1`, which means that each Virtual Machine in this Scale Set should not be evicted for price reasons.

-> **NOTE:** This can only be configured when `priority` is set to `Spot`.

* `overprovision` - (Optional) Should Azure over-provision Virtual Machines in this Scale Set? This means that multiple Virtual Machines will be provisioned and Azure will keep the instances which become available first - which improves provisioning success rates and improves deployment time. You're not billed for these over-provisioned VM's and they don't count towards the Subscription Quota. Defaults to `true`.

* `plan` - (Optional) A `plan` block as documented below.

-> **NOTE:** When using an image from Azure Marketplace a `plan` must be specified.

* `platform_fault_domain_count` - (Optional) Specifies the number of fault domains that are used by this Linux Virtual Machine Scale Set. Changing this forces a new resource to be created.

* `priority` - (Optional) The Priority of this Virtual Machine Scale Set. Possible values are `Regular` and `Spot`. Defaults to `Regular`. Changing this value forces a new resource.

-> **NOTE:** When `priority` is set to `Spot` an `eviction_policy` must be specified.

* `provision_vm_agent` - (Optional) Should the Azure VM Agent be provisioned on each Virtual Machine in the Scale Set? Defaults to `true`. Changing this value forces a new resource to be created.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group in which the Virtual Machine Scale Set should be assigned to. Changing this forces a new resource to be created.

* `rolling_upgrade_policy` - (Optional) A `rolling_upgrade_policy` block as defined below. This is Required and can only be specified when `upgrade_mode` is set to `Automatic` or `Rolling`.

* `scale_in_policy` - (Optional) The scale-in policy rule that decides which virtual machines are chosen for removal when a Virtual Machine Scale Set is scaled in. Possible values for the scale-in policy rules are `Default`, `NewestVM` and `OldestVM`, defaults to `Default`. For more information about scale in policy, please [refer to this doc](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-scale-in-policy).

* `secret` - (Optional) One or more `secret` blocks as defined below.

* `single_placement_group` - (Optional) Should this Virtual Machine Scale Set be limited to a Single Placement Group, which means the number of instances will be capped at 100 Virtual Machines. Defaults to `true`.

* `source_image_id` - (Optional) The ID of an Image which each Virtual Machine in this Scale Set should be based on.

-> **NOTE:** One of either `source_image_id` or `source_image_reference` must be set.

* `source_image_reference` - (Optional) A `source_image_reference` block as defined below.

-> **NOTE:** One of either `source_image_id` or `source_image_reference` must be set.

* `tags` - (Optional) A mapping of tags which should be assigned to this Virtual Machine Scale Set.

* `terminate_notification` - (Optional) A `terminate_notification` block as defined below.

* `upgrade_mode` - (Optional) Specifies how Upgrades (e.g. changing the Image/SKU) should be performed to Virtual Machine Instances. Possible values are `Automatic`, `Manual` and `Rolling`. Defaults to `Manual`.

* `zone_balance` - (Optional) Should the Virtual Machines in this Scale Set be strictly evenly distributed across Availability Zones? Defaults to `false`. Changing this forces a new resource to be created.

-> **NOTE:** This can only be set to `true` when one or more `zones` are configured.

* `zones` - (Optional) A list of Availability Zones in which the Virtual Machines in this Scale Set should be created in. Changing this forces a new resource to be created.

---

A `additional_capabilities` block supports the following:

* `ultra_ssd_enabled` - (Optional) Should the capacity to enable Data Disks of the `UltraSSD_LRS` storage account type be supported on this Virtual Machine Scale Set? Defaults to `false`. Changing this forces a new resource to be created.

---

A `admin_ssh_key` block supports the following:

* `public_key` - (Required) The Public Key which should be used for authentication, which needs to be at least 2048-bit and in `ssh-rsa` format.

* `username` - (Required) The Username for which this Public SSH Key should be configured.

-> **NOTE:** The Azure VM Agent only allows creating SSH Keys at the path `/home/{username}/.ssh/authorized_keys` - as such this public key will be added/appended to the authorized keys file.

---

A `automatic_os_upgrade_policy` block supports the following:

* `disable_automatic_rollback` - (Required) Should automatic rollbacks be disabled? Changing this forces a new resource to be created.

* `enable_automatic_os_upgrade` - (Required) Should OS Upgrades automatically be applied to Scale Set instances in a rolling fashion when a newer version of the OS Image becomes available? Changing this forces a new resource to be created.

---

A `automatic_instance_repair` block supports the following:

* `enabled` - (Required) Should the automatic instance repair be enabled on this Virtual Machine Scale Set?

* `grace_period` - (Optional) Amount of time (in minutes, between 30 and 90, defaults to 30 minutes) for which automatic repairs will be delayed. The grace period starts right after the VM is found unhealthy. The time duration should be specified in ISO 8601 format.

---

A `boot_diagnostics` block supports the following:

* `storage_account_uri` - (Optional) The Primary/Secondary Endpoint for the Azure Storage Account which should be used to store Boot Diagnostics, including Console Output and Screenshots from the Hypervisor. Passing a null value will utilize a Managed Storage Account to store Boot Diagnostics.

---

A `certificate` block supports the following:

* `url` - (Required) The Secret URL of a Key Vault Certificate.

-> **NOTE:** This can be sourced from the `secret_id` field within the `azurerm_key_vault_certificate` Resource.
~> **NOTE:** The certificate must have been uploaded/created in PFX format, PEM certificates are not currently supported by Azure.


---

A `data_disk` block supports the following:

* `caching` - (Required) The type of Caching which should be used for this Data Disk. Possible values are `None`, `ReadOnly` and `ReadWrite`.

* `create_option` - (Optional) The create option which should be used for this Data Disk. Possible values are `Empty` and `FromImage`. Defaults to `Empty`. (`FromImage` should only be used if the source image includes data disks).

* `disk_size_gb` - (Required) The size of the Data Disk which should be created.

* `lun` - (Required) The Logical Unit Number of the Data Disk, which must be unique within the Virtual Machine.

* `storage_account_type` - (Required) The Type of Storage Account which should back this Data Disk. Possible values include `Standard_LRS`, `StandardSSD_LRS`, `Premium_LRS` and `UltraSSD_LRS`.

-> **NOTE:** `UltraSSD_LRS` is only supported when `ultra_ssd_enabled` within the `additional_capabilities` block is enabled.

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to encrypt this Data Disk.

-> **NOTE:** The Disk Encryption Set must have the `Reader` Role Assignment scoped on the Key Vault - in addition to an Access Policy to the Key Vault

~> **NOTE:** Disk Encryption Sets are in Public Preview in a limited set of regions

* `disk_iops_read_write` - (Optional) Specifies the Read-Write IOPS for this Data Disk. Only settable for UltraSSD disks.

* `disk_mbps_read_write` - (Optional) Specifies the bandwidth in MB per second for this Data Disk. Only settable for UltraSSD disks.

* `write_accelerator_enabled` - (Optional) Should Write Accelerator be enabled for this Data Disk? Defaults to `false`.

-> **NOTE:** This requires that the `storage_account_type` is set to `Premium_LRS` and that `caching` is set to `None`.

---

A `diff_disk_settings` block supports the following:

`option` - (Required) Specifies the Ephemeral Disk Settings for the OS Disk. At this time the only possible value is `Local`. Changing this forces a new resource to be created.

---

An `extension` block supports the following:

!> **NOTE:** This block is only available in the Opt-In beta and requires that the Environment Variable `ARM_PROVIDER_VMSS_EXTENSIONS_BETA` is set to `true` to be used.

* `name` - (Required) The name for the Virtual Machine Scale Set Extension.

* `publisher` - (Required) Specifies the Publisher of the Extension.

* `type` - (Required) Specifies the Type of the Extension.

* `type_handler_version` - (Required) Specifies the version of the extension to use, available versions can be found using the Azure CLI.

* `auto_upgrade_minor_version` - (Optional) Should the latest version of the Extension be used at Deployment Time, if one is available? This won't auto-update the extension on existing installation. Defaults to `true`.

* `force_update_tag` - (Optional) A value which, when different to the previous value can be used to force-run the Extension even if the Extension Configuration hasn't changed.

* `protected_settings` - (Optional) A JSON String which specifies Sensitive Settings (such as Passwords) for the Extension.

~> **NOTE:** Keys within the `protected_settings` block are notoriously case-sensitive, where the casing required (e.g. TitleCase vs snakeCase) depends on the Extension being used. Please refer to the documentation for the specific Virtual Machine Extension you're looking to use for more information.

-> **Note:** Rather than defining JSON inline [you can use the `jsonencode` interpolation function](https://www.terraform.io/docs/configuration/functions/jsonencode.html) to define this in a cleaner way.

* `provision_after_extensions` - (Optional) An ordered list of Extension names which this should be provisioned after.

* `settings` - (Optional) A JSON String which specifies Settings for the Extension.

~> **NOTE:** Keys within the `settings` block are notoriously case-sensitive, where the casing required (e.g. TitleCase vs snakeCase) depends on the Extension being used. Please refer to the documentation for the specific Virtual Machine Extension you're looking to use for more information.

-> **Note:** Rather than defining JSON inline [you can use the `jsonencode` interpolation function](https://www.terraform.io/docs/configuration/functions/jsonencode.html) to define this in a cleaner way.

---

A `identity` block supports the following:

* `type` - (Required) The type of Managed Identity which should be assigned to the Linux Virtual Machine Scale Set. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) A list of User Managed Identity ID's which should be assigned to the Linux Virtual Machine Scale Set.

~> **NOTE:** This is required when `type` is set to `UserAssigned`.

---

A `ip_configuration` block supports the following:

* `name` - (Required) The Name which should be used for this IP Configuration.

* `application_gateway_backend_address_pool_ids` - (Optional) A list of Backend Address Pools ID's from a Application Gateway which this Virtual Machine Scale Set should be connected to.

* `application_security_group_ids` - (Optional) A list of Application Security Group ID's which this Virtual Machine Scale Set should be connected to.

* `load_balancer_backend_address_pool_ids` - (Optional) A list of Backend Address Pools ID's from a Load Balancer which this Virtual Machine Scale Set should be connected to.

-> **NOTE:** When using this field you'll also need to configure a Rule for the Load Balancer, and use a `depends_on` between this resource and the Load Balancer Rule.

* `load_balancer_inbound_nat_rules_ids` - (Optional) A list of NAT Rule ID's from a Load Balancer which this Virtual Machine Scale Set should be connected to.

-> **NOTE:** When using this field you'll also need to configure a Rule for the Load Balancer, and use a `depends_on` between this resource and the Load Balancer Rule.

* `primary` - (Optional) Is this the Primary IP Configuration for this Network Interface? Defaults to `false`.

-> **NOTE:** One `ip_configuration` block must be marked as Primary for each Network Interface.

* `public_ip_address` - (Optional) A `public_ip_address` block as defined below.

* `subnet_id` - (Optional) The ID of the Subnet which this IP Configuration should be connected to.

~> `subnet_id` is required if `version` is set to `IPv4`.

* `version` - (Optional) The Internet Protocol Version which should be used for this IP Configuration. Possible values are `IPv4` and `IPv6`. Defaults to `IPv4`.

---

A `ip_tag` block supports the following:

* `tag` - The IP Tag associated with the Public IP, such as `SQL` or `Storage`.

* `type` - The Type of IP Tag, such as `FirstPartyUsage`.

---

A `network_interface` block supports the following:

* `name` - (Required) The Name which should be used for this Network Interface. Changing this forces a new resource to be created.

* `ip_configuration` - (Required) One or more `ip_configuration` blocks as defined above.

* `dns_servers` - (Optional) A list of IP Addresses of DNS Servers which should be assigned to the Network Interface.

* `enable_accelerated_networking` - (Optional) Does this Network Interface support Accelerated Networking? Defaults to `false`.

* `enable_ip_forwarding` - (Optional) Does this Network Interface support IP Forwarding? Defaults to `false`.

* `network_security_group_id` - (Optional) The ID of a Network Security Group which should be assigned to this Network Interface.

* `primary` - (Optional) Is this the Primary IP Configuration?

-> **NOTE:** If multiple `network_interface` blocks are specified, one must be set to `primary`.

---

A `os_disk` block supports the following:

* `caching` - (Required) The Type of Caching which should be used for the Internal OS Disk. Possible values are `None`, `ReadOnly` and `ReadWrite`.

* `storage_account_type` - (Required) The Type of Storage Account which should back this the Internal OS Disk. Possible values include `Standard_LRS`, `StandardSSD_LRS` and `Premium_LRS`.

* `diff_disk_settings` - (Optional) A `diff_disk_settings` block as defined above. Changing this forces a new resource to be created.

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to encrypt this OS Disk.

-> **NOTE:** The Disk Encryption Set must have the `Reader` Role Assignment scoped on the Key Vault - in addition to an Access Policy to the Key Vault

~> **NOTE:** Disk Encryption Sets are in Public Preview in a limited set of regions

* `disk_size_gb` - (Optional) The Size of the Internal OS Disk in GB, if you wish to vary from the size used in the image this Virtual Machine Scale Set is sourced from.

-> **NOTE:** If specified this must be equal to or larger than the size of the Image the VM Scale Set is based on. When creating a larger disk than exists in the image you'll need to repartition the disk to use the remaining space.

* `write_accelerator_enabled` - (Optional) Should Write Accelerator be Enabled for this OS Disk? Defaults to `false`.

-> **NOTE:** This requires that the `storage_account_type` is set to `Premium_LRS` and that `caching` is set to `None`.

---

A `plan` block supports the following:

* `name` - (Required) Specifies the name of the image from the marketplace. Changing this forces a new resource to be created.

* `publisher` - (Required) Specifies the publisher of the image. Changing this forces a new resource to be created.

* `product` - (Required) Specifies the product of the image from the marketplace. Changing this forces a new resource to be created.

---

A `public_ip_address` block supports the following:

* `name` - (Required) The Name of the Public IP Address Configuration.

* `domain_name_label` - (Optional) The Prefix which should be used for the Domain Name Label for each Virtual Machine Instance. Azure concatenates the Domain Name Label and Virtual Machine Index to create a unique Domain Name Label for each Virtual Machine.

* `idle_timeout_in_minutes` - (Optional) The Idle Timeout in Minutes for the Public IP Address. Possible values are in the range `4` to `32`.

* `ip_tag` - (Optional) One or more `ip_tag` blocks as defined above.

* `public_ip_prefix_id` - (Optional) The ID of the Public IP Address Prefix from where Public IP Addresses should be allocated. Changing this forces a new resource to be created.

~> **NOTE:** This functionality is in Preview and must be opted into via `az feature register --namespace Microsoft.Network --name AllowBringYourOwnPublicIpAddress` and then `az provider register -n Microsoft.Network`.

---

A `rolling_upgrade_policy` block supports the following:

* `max_batch_instance_percent` - (Required) The maximum percent of total virtual machine instances that will be upgraded simultaneously by the rolling upgrade in one batch. As this is a maximum, unhealthy instances in previous or future batches can cause the percentage of instances in a batch to decrease to ensure higher reliability. Changing this forces a new resource to be created.

* `max_unhealthy_instance_percent` - (Required) The maximum percentage of the total virtual machine instances in the scale set that can be simultaneously unhealthy, either as a result of being upgraded, or by being found in an unhealthy state by the virtual machine health checks before the rolling upgrade aborts. This constraint will be checked prior to starting any batch. Changing this forces a new resource to be created.

* `max_unhealthy_upgraded_instance_percent` - (Required) The maximum percentage of upgraded virtual machine instances that can be found to be in an unhealthy state. This check will happen after each batch is upgraded. If this percentage is ever exceeded, the rolling update aborts. Changing this forces a new resource to be created.

* `pause_time_between_batches` - (Required) The wait time between completing the update for all virtual machines in one batch and starting the next batch. The time duration should be specified in ISO 8601 format. Changing this forces a new resource to be created.

---

A `secret` block supports the following:

* `certificate` - (Required) One or more `certificate` blocks as defined above.

* `key_vault_id` - (Required) The ID of the Key Vault from which all Secrets should be sourced.

---

A `terminate_notification` block supports the following:

* `enabled` - (Required) Should the terminate notification be enabled on this Virtual Machine Scale Set? Defaults to `false`.

* `timeout` - (Optional) Length of time (in minutes, between 5 and 15) a notification to be sent to the VM on the instance metadata server till the VM gets deleted. The time duration should be specified in ISO 8601 format.

~> For more information about the terminate notification, please [refer to this doc](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-terminate-notification).

---

`source_image_reference` supports the following:

* `publisher` - (Optional) Specifies the publisher of the image used to create the virtual machines.

* `offer` - (Optional) Specifies the offer of the image used to create the virtual machines.

* `sku` - (Optional) Specifies the SKU of the image used to create the virtual machines.

* `version` - (Optional) Specifies the version of the image used to create the virtual machines.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Linux Virtual Machine Scale Set.

* `identity` - An `identity` block as defined below.

* `unique_id` - The Unique ID for this Linux Virtual Machine Scale Set.

---

An `identity` block exports the following:

* `principal_id` - The ID of the System Managed Service Principal.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Linux Virtual Machine Scale Set.
* `update` - (Defaults to 60 minutes) Used when updating (and rolling the instances of) the Linux Virtual Machine Scale Set (e.g. when changing SKU).
* `delete` - (Defaults to 30 minutes) Used when deleting the Linux Virtual Machine Scale Set.

## Import

Linux Virtual Machine Scale Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_linux_virtual_machine_scale_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/Microsoft.Compute/virtualMachineScaleSets/scaleset1
```
