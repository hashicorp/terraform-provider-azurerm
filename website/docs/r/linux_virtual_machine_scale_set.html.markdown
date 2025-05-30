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

-> **Note:** As of the **v2.86.0** (November 19, 2021) release of the provider this resource will only create Virtual Machine Scale Sets with the **Uniform** Orchestration Mode. For Virtual Machine Scale Sets with **Flexible** orchestration mode, use [`azurerm_orchestrated_virtual_machine_scale_set`](orchestrated_virtual_machine_scale_set.html). Flexible orchestration mode is recommended for workloads on Azure.

-> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

-> **Note:** Terraform will automatically update & reimage the nodes in the Scale Set (if Required) during an Update - this behaviour can be configured [using the `features` setting within the Provider block](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block).

## Example Usage

This example provisions a basic Linux Virtual Machine Scale Set on an internal network. Additional examples of how to use the `azurerm_linux_virtual_machine_scale_set` resource can be found [in the ./examples/vm-scale-set/linux` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/vm-scale-set/linux).

```hcl
locals {
  first_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
}

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
    public_key = local.first_public_key
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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

* `name` - (Required) The name of the Linux Virtual Machine Scale Set. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Linux Virtual Machine Scale Set should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Linux Virtual Machine Scale Set should be exist. Changing this forces a new resource to be created.

* `admin_username` - (Required) The username of the local administrator on each Virtual Machine Scale Set instance. Changing this forces a new resource to be created.

* `instances` - (Optional) The number of Virtual Machines in the Scale Set. Defaults to `0`.

-> **Note:** If you are using AutoScaling, you may wish to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changess) to ignore changes to this field.

* `sku` - (Required) The Virtual Machine SKU for the Scale Set, such as `Standard_F2`.

* `network_interface` - (Required) One or more `network_interface` blocks as defined below.

* `os_disk` - (Required) An `os_disk` block as defined below.

---

* `additional_capabilities` - (Optional) An `additional_capabilities` block as defined below.

* `admin_password` - (Optional) The Password which should be used for the local-administrator on this Virtual Machine. Changing this forces a new resource to be created.

-> **Note:** When an `admin_password` is specified `disable_password_authentication` must be set to `false`.

-> **Note:** One of either `admin_password` or `admin_ssh_key` must be specified.

* `admin_ssh_key` - (Optional) One or more `admin_ssh_key` blocks as defined below.

-> **Note:** One of either `admin_password` or `admin_ssh_key` must be specified.

* `automatic_os_upgrade_policy` - (Optional) An `automatic_os_upgrade_policy` block as defined below. This can only be specified when `upgrade_mode` is set to either `Automatic` or `Rolling`.

* `automatic_instance_repair` - (Optional) An `automatic_instance_repair` block as defined below. To enable the automatic instance repair, this Virtual Machine Scale Set must have a valid `health_probe_id` or an [Application Health Extension](https://docs.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-health-extension).

-> **Note:** For more information about Automatic Instance Repair, please refer to the [product documentation](https://docs.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-automatic-instance-repairs).

* `boot_diagnostics` - (Optional) A `boot_diagnostics` block as defined below.

* `capacity_reservation_group_id` - (Optional) Specifies the ID of the Capacity Reservation Group which the Virtual Machine Scale Set should be allocated to. Changing this forces a new resource to be created.

-> **Note:** `capacity_reservation_group_id` cannot be used with `proximity_placement_group_id`

~> **Note:** `single_placement_group` must be set to `false` when `capacity_reservation_group_id` is specified.

* `computer_name_prefix` - (Optional) The prefix which should be used for the name of the Virtual Machines in this Scale Set. If unspecified this defaults to the value for the `name` field. If the value of the `name` field is not a valid `computer_name_prefix`, then you must specify `computer_name_prefix`. Changing this forces a new resource to be created.

* `custom_data` - (Optional) The Base64-Encoded Custom Data which should be used for this Virtual Machine Scale Set.

-> **Note:** When Custom Data has been configured, it's not possible to remove it without tainting the Virtual Machine Scale Set, due to a limitation of the Azure API.

* `data_disk` - (Optional) One or more `data_disk` blocks as defined below.

* `disable_password_authentication` - (Optional) Should Password Authentication be disabled on this Virtual Machine Scale Set? Defaults to `true`.

-> **Note:** In general we'd recommend using SSH Keys for authentication rather than Passwords - but there's tradeoff's to each - please [see this thread for more information](https://security.stackexchange.com/questions/69407/why-is-using-an-ssh-key-more-secure-than-using-passwords).

-> **Note:** When a `admin_password` is specified `disable_password_authentication` must be set to `false`.

* `do_not_run_extensions_on_overprovisioned_machines` - (Optional) Should Virtual Machine Extensions be run on Overprovisioned Virtual Machines in the Scale Set? Defaults to `false`.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Linux Virtual Machine Scale Set should exist. Changing this forces a new Linux Virtual Machine Scale Set to be created.

* `encryption_at_host_enabled` - (Optional) Should all of the disks (including the temp disk) attached to this Virtual Machine be encrypted by enabling Encryption at Host?

* `extension` - (Optional) One or more `extension` blocks as defined below

* `extension_operations_enabled` - (Optional) Should extension operations be allowed on the Virtual Machine Scale Set? Possible values are `true` or `false`. Defaults to `true`. Changing this forces a new Linux Virtual Machine Scale Set to be created.

-> **Note:** `extension_operations_enabled` may only be set to `false` if there are no extensions defined in the `extension` field.

* `extensions_time_budget` - (Optional) Specifies the duration allocated for all extensions to start. The time duration should be between `15` minutes and `120` minutes (inclusive) and should be specified in ISO 8601 format. Defaults to `PT1H30M`.

* `eviction_policy` - (Optional) Specifies the eviction policy for Virtual Machines in this Scale Set. Possible values are `Deallocate` and `Delete`. Changing this forces a new resource to be created.

-> **Note:** This can only be configured when `priority` is set to `Spot`.

* `gallery_application` - (Optional) One or more `gallery_application` blocks as defined below.

* `health_probe_id` - (Optional) The ID of a Load Balancer Probe which should be used to determine the health of an instance. This is Required and can only be specified when `upgrade_mode` is set to `Automatic` or `Rolling`.

* `host_group_id` - (Optional) Specifies the ID of the dedicated host group that the virtual machine scale set resides in. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `max_bid_price` - (Optional) The maximum price you're willing to pay for each Virtual Machine in this Scale Set, in US Dollars; which must be greater than the current spot price. If this bid price falls below the current spot price the Virtual Machines in the Scale Set will be evicted using the `eviction_policy`. Defaults to `-1`, which means that each Virtual Machine in this Scale Set should not be evicted for price reasons.

-> **Note:** This can only be configured when `priority` is set to `Spot`.

* `overprovision` - (Optional) Should Azure over-provision Virtual Machines in this Scale Set? This means that multiple Virtual Machines will be provisioned and Azure will keep the instances which become available first - which improves provisioning success rates and improves deployment time. You're not billed for these over-provisioned VM's and they don't count towards the Subscription Quota. Defaults to `true`.

* `plan` - (Optional) A `plan` block as defined below. Changing this forces a new resource to be created.

-> **Note:** When using an image from Azure Marketplace a `plan` must be specified.

* `platform_fault_domain_count` - (Optional) Specifies the number of fault domains that are used by this Linux Virtual Machine Scale Set. Changing this forces a new resource to be created.

* `priority` - (Optional) The Priority of this Virtual Machine Scale Set. Possible values are `Regular` and `Spot`. Defaults to `Regular`. Changing this value forces a new resource.

-> **Note:** When `priority` is set to `Spot` an `eviction_policy` must be specified.

* `provision_vm_agent` - (Optional) Should the Azure VM Agent be provisioned on each Virtual Machine in the Scale Set? Defaults to `true`. Changing this value forces a new resource to be created.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group in which the Virtual Machine Scale Set should be assigned to. Changing this forces a new resource to be created.

* `rolling_upgrade_policy` - (Optional) A `rolling_upgrade_policy` block as defined below. This is Required and can only be specified when `upgrade_mode` is set to `Automatic` or `Rolling`. Changing this forces a new resource to be created.

* `scale_in` - (Optional) A `scale_in` block as defined below.

* `secret` - (Optional) One or more `secret` blocks as defined below.

* `secure_boot_enabled` - (Optional) Specifies whether secure boot should be enabled on the virtual machine. Changing this forces a new resource to be created.

* `single_placement_group` - (Optional) Should this Virtual Machine Scale Set be limited to a Single Placement Group, which means the number of instances will be capped at 100 Virtual Machines. Defaults to `true`.

* `source_image_id` - (Optional) The ID of an Image which each Virtual Machine in this Scale Set should be based on. Possible Image ID types include `Image ID`, `Shared Image ID`, `Shared Image Version ID`, `Community Gallery Image ID`, `Community Gallery Image Version ID`, `Shared Gallery Image ID` and `Shared Gallery Image Version ID`.

-> **Note:** One of either `source_image_id` or `source_image_reference` must be set.

* `source_image_reference` - (Optional) A `source_image_reference` block as defined below.

-> **Note:** One of either `source_image_id` or `source_image_reference` must be set.

* `spot_restore` - (Optional) A `spot_restore` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to this Virtual Machine Scale Set.

* `termination_notification` - (Optional) A `termination_notification` block as defined below.

* `upgrade_mode` - (Optional) Specifies how Upgrades (e.g. changing the Image/SKU) should be performed to Virtual Machine Instances. Possible values are `Automatic`, `Manual` and `Rolling`. Defaults to `Manual`. Changing this forces a new resource to be created.

-> **Note:** If rolling upgrades are configured and running on a Linux Virtual Machine Scale Set, they will be cancelled when Terraform tries to destroy the resource.

* `user_data` - (Optional) The Base64-Encoded User Data which should be used for this Virtual Machine Scale Set.

* `vtpm_enabled` - (Optional) Specifies whether vTPM should be enabled on the virtual machine. Changing this forces a new resource to be created.

* `zone_balance` - (Optional) Should the Virtual Machines in this Scale Set be strictly evenly distributed across Availability Zones? Defaults to `false`. Changing this forces a new resource to be created.

-> **Note:** This can only be set to `true` when one or more `zones` are configured.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Linux Virtual Machine Scale Set should be located.

-> **Note:** Updating `zones` to remove an existing zone forces a new Virtual Machine Scale Set to be created.

---

An `additional_capabilities` block supports the following:

* `ultra_ssd_enabled` - (Optional) Should the capacity to enable Data Disks of the `UltraSSD_LRS` storage account type be supported on this Virtual Machine Scale Set? Possible values are `true` or `false`. Defaults to `false`. Changing this forces a new resource to be created.

---

An `admin_ssh_key` block supports the following:

* `public_key` - (Required) The Public Key which should be used for authentication, which needs to be in `ssh-rsa` format with at least 2048-bit or in `ssh-ed25519` format.

* `username` - (Required) The Username for which this Public SSH Key should be configured.

-> **Note:** The Azure VM Agent only allows creating SSH Keys at the path `/home/{username}/.ssh/authorized_keys` - as such this public key will be added/appended to the authorized keys file.

---

An `automatic_os_upgrade_policy` block supports the following:

* `disable_automatic_rollback` - (Required) Should automatic rollbacks be disabled?

* `enable_automatic_os_upgrade` - (Required) Should OS Upgrades automatically be applied to Scale Set instances in a rolling fashion when a newer version of the OS Image becomes available?

---

An `automatic_instance_repair` block supports the following:

* `enabled` - (Required) Should the automatic instance repair be enabled on this Virtual Machine Scale Set?

* `grace_period` - (Optional) Amount of time for which automatic repairs will be delayed. The grace period starts right after the VM is found unhealthy. Possible values are between `10` and `90` minutes. The time duration should be specified in `ISO 8601` format (e.g. `PT10M` to `PT90M`).

-> **Note:** Once the `grace_period` field has been set it will always return the last value it was assigned if it is removed from the configuration file.

* `action` - (Optional) The repair action that will be used for repairing unhealthy virtual machines in the scale set. Possible values include `Replace`, `Restart`, `Reimage`.

-> **Note:** Once the `action` field has been set it will always return the last value it was assigned if it is removed from the configuration file.

-> **Note:** If you wish to update the repair `action` of an existing `automatic_instance_repair` policy, you must first `disable` the `automatic_instance_repair` policy before you can re-enable the `automatic_instance_repair` policy with the new repair `action` defined.

---

A `boot_diagnostics` block supports the following:

* `storage_account_uri` - (Optional) The Primary/Secondary Endpoint for the Azure Storage Account which should be used to store Boot Diagnostics, including Console Output and Screenshots from the Hypervisor.

-> **Note:** Passing a null value will utilize a Managed Storage Account to store Boot Diagnostics.

---

A `certificate` block supports the following:

* `url` - (Required) The Secret URL of a Key Vault Certificate.

-> **Note:** This can be sourced from the `secret_id` field within the `azurerm_key_vault_certificate` Resource.

-> **Note:** The certificate must have been uploaded/created in PFX format, PEM certificates are not currently supported by Azure.

---

A `data_disk` block supports the following:

* `name` - (Optional) The name of the Data Disk.

* `caching` - (Required) The type of Caching which should be used for this Data Disk. Possible values are `None`, `ReadOnly` and `ReadWrite`.

* `create_option` - (Optional) The create option which should be used for this Data Disk. Possible values are `Empty` and `FromImage`. Defaults to `Empty`. (`FromImage` should only be used if the source image includes data disks).

* `disk_size_gb` - (Required) The size of the Data Disk which should be created.

* `lun` - (Required) The Logical Unit Number of the Data Disk, which must be unique within the Virtual Machine.

* `storage_account_type` - (Required) The Type of Storage Account which should back this Data Disk. Possible values include `Standard_LRS`, `StandardSSD_LRS`, `StandardSSD_ZRS`, `Premium_LRS`, `PremiumV2_LRS`, `Premium_ZRS` and `UltraSSD_LRS`.

-> **Note:** `UltraSSD_LRS` is only supported when `ultra_ssd_enabled` within the `additional_capabilities` block is enabled.

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to encrypt this Data Disk. Changing this forces a new resource to be created.

-> **Note:** The Disk Encryption Set must have the `Reader` Role Assignment scoped on the Key Vault - in addition to an Access Policy to the Key Vault

-> **Note:** Disk Encryption Sets are in Public Preview in a limited set of regions

* `ultra_ssd_disk_iops_read_write` - (Optional) Specifies the Read-Write IOPS for this Data Disk. Only settable when `storage_account_type` is `PremiumV2_LRS` or `UltraSSD_LRS`.

* `ultra_ssd_disk_mbps_read_write` - (Optional) Specifies the bandwidth in MB per second for this Data Disk. Only settable when `storage_account_type` is `PremiumV2_LRS` or `UltraSSD_LRS`.

* `write_accelerator_enabled` - (Optional) Should Write Accelerator be enabled for this Data Disk? Defaults to `false`.

-> **Note:** This requires that the `storage_account_type` is set to `Premium_LRS` and that `caching` is set to `None`.

---

A `diff_disk_settings` block supports the following:

* `option` - (Required) Specifies the Ephemeral Disk Settings for the OS Disk. At this time the only possible value is `Local`. Changing this forces a new resource to be created.

* `placement` - (Optional) Specifies where to store the Ephemeral Disk. Possible values are `CacheDisk` and `ResourceDisk`. Defaults to `CacheDisk`. Changing this forces a new resource to be created.

---

An `extension` block supports the following:

* `name` - (Required) The name for the Virtual Machine Scale Set Extension.

* `publisher` - (Required) Specifies the Publisher of the Extension.

* `type` - (Required) Specifies the Type of the Extension.

* `type_handler_version` - (Required) Specifies the version of the extension to use, available versions can be found using the Azure CLI.

* `auto_upgrade_minor_version` - (Optional) Should the latest version of the Extension be used at Deployment Time, if one is available? This won't auto-update the extension on existing installation. Defaults to `true`.

* `automatic_upgrade_enabled` - (Optional) Should the Extension be automatically updated whenever the Publisher releases a new version of this VM Extension? 

* `force_update_tag` - (Optional) A value which, when different to the previous value can be used to force-run the Extension even if the Extension Configuration hasn't changed.

* `protected_settings` - (Optional) A JSON String which specifies Sensitive Settings (such as Passwords) for the Extension.

-> **Note:** Keys within the `protected_settings` block are notoriously case-sensitive, where the casing required (e.g. TitleCase vs snakeCase) depends on the Extension being used. Please refer to the documentation for the specific Virtual Machine Extension you're looking to use for more information.

-> **Note:** Rather than defining JSON inline [you can use the `jsonencode` interpolation function](https://www.terraform.io/docs/configuration/functions/jsonencode.html) to define this in a cleaner way.

* `protected_settings_from_key_vault` - (Optional) A `protected_settings_from_key_vault` block as defined below.

~> **Note:** `protected_settings_from_key_vault` cannot be used with `protected_settings`

* `provision_after_extensions` - (Optional) An ordered list of Extension names which this should be provisioned after.

* `settings` - (Optional) A JSON String which specifies Settings for the Extension.

-> **Note:** Keys within the `settings` block are notoriously case-sensitive, where the casing required (e.g. TitleCase vs snakeCase) depends on the Extension being used. Please refer to the documentation for the specific Virtual Machine Extension you're looking to use for more information.

-> **Note:** Rather than defining JSON inline [you can use the `jsonencode` interpolation function](https://www.terraform.io/docs/configuration/functions/jsonencode.html) to define this in a cleaner way.

---

A `gallery_application` block supports the following:

* `version_id` - (Required) Specifies the Gallery Application Version resource ID. Changing this forces a new resource to be created.

* `configuration_blob_uri` - (Optional) Specifies the URI to an Azure Blob that will replace the default configuration for the package if provided. Changing this forces a new resource to be created.

* `order` - (Optional) Specifies the order in which the packages have to be installed. Possible values are between `0` and `2147483647`. Defaults to `0`. Changing this forces a new resource to be created.

* `tag` - (Optional) Specifies a passthrough value for more generic context. This field can be any valid `string` value. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Linux Virtual Machine Scale Set. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Linux Virtual Machine Scale Set.

-> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `ip_configuration` block supports the following:

* `name` - (Required) The Name which should be used for this IP Configuration.

* `application_gateway_backend_address_pool_ids` - (Optional) A list of Backend Address Pools ID's from a Application Gateway which this Virtual Machine Scale Set should be connected to.

* `application_security_group_ids` - (Optional) A list of Application Security Group ID's which this Virtual Machine Scale Set should be connected to.

* `load_balancer_backend_address_pool_ids` - (Optional) A list of Backend Address Pools ID's from a Load Balancer which this Virtual Machine Scale Set should be connected to.

-> **Note:** When the Virtual Machine Scale Set is configured to have public IPs per instance are created with a load balancer, the SKU of the Virtual Machine instance IPs is determined by the SKU of the Virtual Machine Scale Sets Load Balancer (e.g. `Basic` or `Standard`). Alternatively, you may use the `public_ip_prefix_id` field to generate instance-level IPs in a virtual machine scale set as well. The zonal properties of the prefix will be passed to the Virtual Machine instance IPs, though they will not be shown in the output. To view the public IP addresses assigned to the Virtual Machine Scale Sets Virtual Machine instances use the **az vmss list-instance-public-ips --resource-group `ResourceGroupName` --name `VirtualMachineScaleSetName`** CLI command.

-> **Note:** When using this field you'll also need to configure a Rule for the Load Balancer, and use a `depends_on` between this resource and the Load Balancer Rule.

* `load_balancer_inbound_nat_rules_ids` - (Optional) A list of NAT Rule ID's from a Load Balancer which this Virtual Machine Scale Set should be connected to.

-> **Note:** When using this field you'll also need to configure a Rule for the Load Balancer, and use a `depends_on` between this resource and the Load Balancer Rule.

* `primary` - (Optional) Is this the Primary IP Configuration for this Network Interface? Defaults to `false`.

-> **Note:** One `ip_configuration` block must be marked as Primary for each Network Interface.

* `public_ip_address` - (Optional) A `public_ip_address` block as defined below.

* `subnet_id` - (Optional) The ID of the Subnet which this IP Configuration should be connected to.

-> **Note:** `subnet_id` is required if `version` is set to `IPv4`.

* `version` - (Optional) The Internet Protocol Version which should be used for this IP Configuration. Possible values are `IPv4` and `IPv6`. Defaults to `IPv4`.

---

An `ip_tag` block supports the following:

* `tag` - (Required) The IP Tag associated with the Public IP, such as `SQL` or `Storage`. Changing this forces a new resource to be created.

* `type` - (Required) The Type of IP Tag, such as `FirstPartyUsage`. Changing this forces a new resource to be created.

---

A `network_interface` block supports the following:

* `name` - (Required) The Name which should be used for this Network Interface. Changing this forces a new resource to be created.

* `ip_configuration` - (Required) One or more `ip_configuration` blocks as defined above.

* `dns_servers` - (Optional) A list of IP Addresses of DNS Servers which should be assigned to the Network Interface.

* `enable_accelerated_networking` - (Optional) Does this Network Interface support Accelerated Networking? Defaults to `false`.

* `enable_ip_forwarding` - (Optional) Does this Network Interface support IP Forwarding? Defaults to `false`.

* `network_security_group_id` - (Optional) The ID of a Network Security Group which should be assigned to this Network Interface.

* `primary` - (Optional) Is this the Primary IP Configuration?

-> **Note:** If multiple `network_interface` blocks are specified, one must be set to `primary`.

---

An `os_disk` block supports the following:

* `caching` - (Required) The Type of Caching which should be used for the Internal OS Disk. Possible values are `None`, `ReadOnly` and `ReadWrite`.

* `storage_account_type` - (Required) The Type of Storage Account which should back this the Internal OS Disk. Possible values include `Standard_LRS`, `StandardSSD_LRS`, `StandardSSD_ZRS`, `Premium_LRS` and `Premium_ZRS`. Changing this forces a new resource to be created.

* `diff_disk_settings` - (Optional) A `diff_disk_settings` block as defined above. Changing this forces a new resource to be created.

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to encrypt this OS Disk. Conflicts with `secure_vm_disk_encryption_set_id`. Changing this forces a new resource to be created.

-> **Note:** The Disk Encryption Set must have the `Reader` Role Assignment scoped on the Key Vault - in addition to an Access Policy to the Key Vault

-> **Note:** Disk Encryption Sets are in Public Preview in a limited set of regions

* `disk_size_gb` - (Optional) The Size of the Internal OS Disk in GB, if you wish to vary from the size used in the image this Virtual Machine Scale Set is sourced from.

-> **Note:** If specified this must be equal to or larger than the size of the Image the VM Scale Set is based on. When creating a larger disk than exists in the image you'll need to repartition the disk to use the remaining space.

* `secure_vm_disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to Encrypt the OS Disk when the Virtual Machine Scale Set is Confidential VMSS. Conflicts with `disk_encryption_set_id`. Changing this forces a new resource to be created.

-> **Note:** `secure_vm_disk_encryption_set_id` can only be specified when `security_encryption_type` is set to `DiskWithVMGuestState`.

* `security_encryption_type` - (Optional) Encryption Type when the Virtual Machine Scale Set is Confidential VMSS. Possible values are `VMGuestStateOnly` and `DiskWithVMGuestState`. Changing this forces a new resource to be created.

-> **Note:** `vtpm_enabled` must be set to `true` when `security_encryption_type` is specified.

-> **Note:** `encryption_at_host_enabled` cannot be set to `true` when `security_encryption_type` is set to `DiskWithVMGuestState`.

* `write_accelerator_enabled` - (Optional) Should Write Accelerator be Enabled for this OS Disk? Defaults to `false`.

-> **Note:** This requires that the `storage_account_type` is set to `Premium_LRS` and that `caching` is set to `None`.

---

A `plan` block supports the following:

* `name` - (Required) Specifies the name of the image from the marketplace. Changing this forces a new resource to be created.

* `publisher` - (Required) Specifies the publisher of the image. Changing this forces a new resource to be created.

* `product` - (Required) Specifies the product of the image from the marketplace. Changing this forces a new resource to be created.

---

A `protected_settings_from_key_vault` block supports the following:

* `secret_url` - (Required) The URL to the Key Vault Secret which stores the protected settings.

* `source_vault_id` - (Required) The ID of the source Key Vault.

---

A `scale_in` block supports the following:

* `rule` - (Optional) The scale-in policy rule that decides which virtual machines are chosen for removal when a Virtual Machine Scale Set is scaled in. Possible values for the scale-in policy rules are `Default`, `NewestVM` and `OldestVM`, defaults to `Default`. For more information about scale in policy, please [refer to this doc](https://docs.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-scale-in-policy).

* `force_deletion_enabled` - (Optional) Should the virtual machines chosen for removal be force deleted when the virtual machine scale set is being scaled-in? Possible values are `true` or `false`. Defaults to `false`.

---

A `public_ip_address` block supports the following:

* `name` - (Required) The Name of the Public IP Address Configuration.

* `domain_name_label` - (Optional) The Prefix which should be used for the Domain Name Label for each Virtual Machine Instance. Azure concatenates the Domain Name Label and Virtual Machine Index to create a unique Domain Name Label for each Virtual Machine.

* `idle_timeout_in_minutes` - (Optional) The Idle Timeout in Minutes for the Public IP Address. Possible values are in the range `4` to `32`.

* `ip_tag` - (Optional) One or more `ip_tag` blocks as defined above. Changing this forces a new resource to be created.

* `public_ip_prefix_id` - (Optional) The ID of the Public IP Address Prefix from where Public IP Addresses should be allocated. Changing this forces a new resource to be created.

-> **Note:** This functionality is in Preview and must be opted into via `az feature register --namespace Microsoft.Network --name AllowBringYourOwnPublicIpAddress` and then `az provider register -n Microsoft.Network`.

* `version` - (Optional) The Internet Protocol Version which should be used for this public IP address. Possible values are `IPv4` and `IPv6`. Defaults to `IPv4`. Changing this forces a new resource to be created.

---

A `rolling_upgrade_policy` block supports the following:

* `cross_zone_upgrades_enabled` - (Optional) Should the Virtual Machine Scale Set ignore the Azure Zone boundaries when constructing upgrade batches? Possible values are `true` or `false`.

* `max_batch_instance_percent` - (Required) The maximum percent of total virtual machine instances that will be upgraded simultaneously by the rolling upgrade in one batch. As this is a maximum, unhealthy instances in previous or future batches can cause the percentage of instances in a batch to decrease to ensure higher reliability.

* `max_unhealthy_instance_percent` - (Required) The maximum percentage of the total virtual machine instances in the scale set that can be simultaneously unhealthy, either as a result of being upgraded, or by being found in an unhealthy state by the virtual machine health checks before the rolling upgrade aborts. This constraint will be checked prior to starting any batch.

* `max_unhealthy_upgraded_instance_percent` - (Required) The maximum percentage of upgraded virtual machine instances that can be found to be in an unhealthy state. This check will happen after each batch is upgraded. If this percentage is ever exceeded, the rolling update aborts.

* `pause_time_between_batches` - (Required) The wait time between completing the update for all virtual machines in one batch and starting the next batch. The time duration should be specified in ISO 8601 format.

* `prioritize_unhealthy_instances_enabled` - (Optional) Upgrade all unhealthy instances in a scale set before any healthy instances. Possible values are `true` or `false`.

* `maximum_surge_instances_enabled` - (Optional) Create new virtual machines to upgrade the scale set, rather than updating the existing virtual machines. Existing virtual machines will be deleted once the new virtual machines are created for each batch. Possible values are `true` or `false`.

-> **Note:** `overprovision` must be set to `false` when `maximum_surge_instances_enabled` is specified.

---

A `secret` block supports the following:

* `certificate` - (Required) One or more `certificate` blocks as defined above.

* `key_vault_id` - (Required) The ID of the Key Vault from which all Secrets should be sourced.

---

A `termination_notification` block supports the following:

* `enabled` - (Required) Should the termination notification be enabled on this Virtual Machine Scale Set? 

* `timeout` - (Optional) Length of time (in minutes, between 5 and 15) a notification to be sent to the VM on the instance metadata server till the VM gets deleted. The time duration should be specified in ISO 8601 format. Defaults to `PT5M`.

-> **Note:** For more information about the termination notification, please [refer to this doc](https://docs.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-terminate-notification).

---

A `source_image_reference` block supports the following:

* `publisher` - (Required) Specifies the publisher of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `offer` - (Required) Specifies the offer of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `sku` - (Required) Specifies the SKU of the image used to create the virtual machines.

* `version` - (Required) Specifies the version of the image used to create the virtual machines.

---

A `spot_restore` block supports the following:

* `enabled` - (Optional) Should the Spot-Try-Restore feature be enabled? The Spot-Try-Restore feature will attempt to automatically restore the evicted Spot Virtual Machine Scale Set VM instances opportunistically based on capacity availability and pricing constraints. Possible values are `true` or `false`. Defaults to `false`. Changing this forces a new resource to be created.

* `timeout` - (Optional) The length of time that the Virtual Machine Scale Set should attempt to restore the Spot VM instances which have been evicted. The time duration should be between `15` minutes and `120` minutes (inclusive). The time duration should be specified in the ISO 8601 format. Defaults to `PT1H`. Changing this forces a new resource to be created.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Linux Virtual Machine Scale Set.

* `identity` - A `identity` block as defined below.

* `unique_id` - The Unique ID for this Linux Virtual Machine Scale Set.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Linux Virtual Machine Scale Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Linux Virtual Machine Scale Set.
* `update` - (Defaults to 1 hour) Used when updating the Linux Virtual Machine Scale Set.
* `delete` - (Defaults to 1 hour) Used when deleting the Linux Virtual Machine Scale Set.

## Import

Linux Virtual Machine Scale Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_linux_virtual_machine_scale_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1
```
