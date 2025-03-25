---
subcategory: "Compute Fleet"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_compute_fleet"
description: |-
  Manages a Compute Fleet in Flexible Orchestration Mode.
---

# azurerm_compute_fleet

Manages a Compute Fleet in Flexible Orchestration Mode.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "examplepip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "example" {
  name                = "example-lb"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_lb_backend_address_pool" "example" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.example.id
}

resource "azurerm_compute_fleet" "example" {
  name                = "example-fleet"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = "testadmin1234"
        password_authentication_enabled = false
        admin_ssh_keys                  = [file("~/.ssh/id_rsa.pub")]
      }
    }

    network_interface {
      name = "networkProTest"
      ip_configuration {
        name                                   = "ipConfig"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.example.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.example.id
      }
      primary_network_interface_enabled = true
    }
  }
}
```

## Arguments Reference

The following arguments are supported:
* `name` - (Required) The name of the Compute Fleet. Changing this forces a new resource to be created.

* `location` - (Required) The azure region where the Compute Fleet should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Compute Fleet should exist. Changing this forces a new resource to be created.

* `virtual_machine_profile` - (Required) A `virtual_machine_profile` block as defined below. Changing this forces a new resource to be created.

* `additional_capabilities_hibernation_enabled` - (Optional) Whether to enable the hibernation capability on the Compute Fleet.  Defaults to `false`. Changing this forces a new resource to be created.

* `additional_capabilities_ultra_ssd_enabled` - (Optional) Whether to enable Data Disks of the `UltraSSD_LRS` storage account type on this Compute Fleet. Defaults to `false`. Changing this forces a new resource to be created.

* `additional_location_profile` - (Optional) One or more `additional_location_profile` blocks as defined below. Changing this forces a new resource to be created.

* `plan` - (Optional) A `plan` block as defined below. Changing this forces a new resource to be created.

* `platform_fault_domain_count` - (Optional)  Specifies the number of fault domains that are used by the Compute Fleet. Defaults to `1`. Changing this forces a new resource to be created.

* `vm_attributes` - (Optional) A `vm_attributes` block as defined below. Changing this forces a new resource to be created.

* `zones` - (Optional) Specifies a list of availability zones in which the Compute Fleet is available. Changing this forces a new resource to be created.

* `compute_api_version` - (Optional) Specifies the `Microsoft.Compute` API version to use when creating the Compute Fleet. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `regular_priority_profile` - (Optional) A `regular_priority_profile` block as defined below.

* `spot_priority_profile` - (Optional) A `spot_priority_profile` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Compute Fleet.

* `vm_sizes_profile` - (Optional) One or more `vm_sizes_profile` blocks as defined below.

-> **Note:** If `spot_priority_profile` is specified, `regular_priority_profile` is not specified and `spot_priority_profile.0.maintain_enabled` is specified as to `false`, changing `vm_sizes_profile` forces a new resource to be created.

---

A `virtual_machine_profile` block supports the following:

* `network_api_version` - (Required) Specifies the Microsoft.Network API version used when creating networking resources in the network interface configurations for the Compute Fleet. Changing this forces a new resource to be created.

* `network_interface` - (Required) One or more `network_interface` blocks as defined above. Changing this forces a new resource to be created.

* `os_profile` - (Required) A `os_profile` block as defined above. Changing this forces a new resource to be created.

* `boot_diagnostic_enabled` - (Optional) Whether to enable the boot diagnostics on the virtual machine. Defaults to `false`. Changing this forces a new resource to be created.

* `boot_diagnostic_storage_account_endpoint` - (Optional) Specifies endpoint of the storage account to use for placing the console output and screenshot. Changing this forces a new resource to be created.

* `capacity_reservation_group_id` - (Optional) Specifies the ID of the capacity reservation group which the Compute Fleet should be allocated to. Changing this forces a new resource to be created.

* `data_disk` - (Optional) One or more `data_disk` blocks as defined above. Changing this forces a new resource to be created.

* `encryption_at_host_enabled` - (Optional) Whether to enable encryption at host on disks attached to the Compute Fleet. Defaults to `false`. Changing this forces a new resource to be created.

* `extension` - (Optional) One or more `extension` blocks as defined above. Changing this forces a new resource to be created.

* `extension_operations_enabled` - (Optional) Whether to enable extension operations on the Compute Fleet.  Defaults to `true`. Changing this forces a new resource to be created.

* `extensions_time_budget` - (Optional) Specifies the time alloted for all extensions to start. Changing this forces a new resource to be created.

* `gallery_application` - (Optional) One or more `gallery_application` blocks as defined above. Changing this forces a new resource to be created.

* `license_type` - (Optional) Specifies the type of on-premise license (also known as Azure Hybrid Use Benefit) which should be used for the Compute Fleet. Possible values are `RHEL_BYOS`, `SLES_BYOS`, `Windows_Client` and `Windows_Server`. Changing this forces a new resource to be created.

* `os_disk` - (Optional) A `os_disk` block as defined above. Changing this forces a new resource to be created.

* `scheduled_event_os_image_timeout` - (Optional) Specifies the length of time a virtual machine being deleted will have to potentially approve the terminate scheduled event before the event is auto approved (timed out). The configuration must be specified in ISO 8601 format. The only possible value is `PT15M`. Changing this forces a new resource to be created.

* `scheduled_event_termination_timeout` - (Optional) Specifies the length of time a virtual machine being reimaged or having its OS upgraded will have to potentially approve the OS image scheduled event before the event is auto approved (timed out). The configuration is specified in ISO 8601 format. Possible values are `PT5M` and `PT15M`. Changing this forces a new resource to be created.

* `secure_boot_enabled` - (Optional) Whether to enable the secure boot on the virtual machine. Defaults to `false`. Changing this forces a new resource to be created.

* `source_image_id` - (Optional) The ID of an image which each virtual machine in the Compute Fleet should be based on. Possible Image ID types include `Image ID`s, `Shared Image ID`s, `Shared Image Version ID`s, `Community Gallery Image ID`s, `Community Gallery Image Version ID`s, `Shared Gallery Image ID`s and `Shared Gallery Image Version ID`s. Changing this forces a new resource to be created.

* `source_image_reference` - (Optional) A `source_image_reference` block as defined above. Changing this forces a new resource to be created.

* `user_data_base64` - (Optional) The base64-encoded User Data which should be used for the Compute Fleet. Changing this forces a new resource to be created.

* `vtpm_enabled` - (Optional) Whether to enable the vTPM on the virtual machine. Defaults to `false`. Changing this forces a new resource to be created.

---

A `network_interface` block supports the following:

* `name` - (Required)  The name which should be used for the network interface. Changing this forces a new resource to be created.

* `ip_configuration` - (Required) One or more `ip_configuration` blocks as defined above. Changing this forces a new resource to be created.

* `accelerated_networking_enabled` - (Optional) Whether to enable accelerated networking on the network interface. Defaults to `false`. Changing this forces a new resource to be created.

* `auxiliary_mode` - (Optional) Specifies the auxiliary mode for the network interface. Possible values are `AcceleratedConnections` and `Floating`. Changing this forces a new resource to be created.

* `auxiliary_sku` - (Optional) Specifies the auxiliary sku for the network interface. Possible values are `A8`, `A4`, `A1` and `A2`. Changing this forces a new resource to be created.

* `delete_option` - (Optional) Specify what happens to the network interface when the virtual machine is deleted. Possible values are `Delete` and `Detach`. Changing this forces a new resource to be created.

* `dns_servers` - (Optional) Specifies a list of IP addresses of DNS servers which should be assigned to the network interface. Changing this forces a new resource to be created.

* `ip_forwarding_enabled` - (Optional) Whether to enable IP forwarding on the network interface. Defaults to `false`. Changing this forces a new resource to be created.

* `network_security_group_id` - (Optional) The ID of the network security group which should be assigned to the network interface. Changing this forces a new resource to be created.

* `primary_network_interface_enabled` - (Optional) Whether to set it as the primary network interface. Defaults to `false`. Changing this forces a new resource to be created.

---

A `os_profile` block supports the following:

* `custom_data_base64` - (Optional) The base64-encoded Custom Data which should be used for the Compute Fleet. Changing this forces a new resource to be created.

* `linux_configuration` - (Optional) A `linux_configuration` block as defined above. Changing this forces a new resource to be created.

* `windows_configuration` - (Optional) A `windows_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `accelerator_count` block supports the following:

* `max` - (Optional) The maximum value of accelerator count.

* `min` - (Optional) The minimum value of accelerator count.

---

A `additional_location_profile` block supports the following:

* `location` - (Required) The Azure Region where the Compute Fleet should exist. Changing this forces a new resource to be created.

* `virtual_machine_profile_override` - (Required) The definition of the `virtual_machine_profile_override` block is the same as the `virtual_machine_profile` block. A `virtual_machine_profile` block as defined below. Changing this forces a new resource to be created.

---

A `additional_unattend_content` block supports the following:

* `content` - (Required) Specifies the XML formatted content that is added to the unattend.xml file for the specified path and component. Changing this forces a new resource to be created.

* `setting` - (Required) Specifies the name of the setting to which the content applies. Possible values are `AutoLogon` and `FirstLogonCommands`. Changing this forces a new resource to be created.

---

A (Windows) `certificate` block supports the following:

* `url` - (Required) The secret URL of a key vault certificate. Changing this forces a new resource to be created.

* `store` - (Optional) The certificate store on the virtual machine where the certificate should be added. Changing this forces a new resource to be created.

---

A (Linux) `certificate` block supports the following:

* `url` - (Required) The secret URL of a key vault certificate. Changing this forces a new resource to be created.

---

A `data_disk` block supports the following:

* `create_option` - (Required) The create option which should be used for the data disk. Possible values are `Empty` and `FromImage`. Changing this forces a new resource to be created.

-> **Note:** `FromImage` should only be used if the source image includes data disks.

* `caching` - (Optional) The type of caching which should be used for the data disk. Possible values are `ReadOnly` and `ReadWrite`. Changing this forces a new resource to be created.

* `delete_option` - (Optional) Specify what happens to the data disk when the virtual machine is deleted. Possible values are `Delete` and `Detach`. Changing this forces a new resource to be created.

* `disk_encryption_set_id` - (Optional) The ID of the disk encryption set which should be used to encrypt the data disk. Changing this forces a new resource to be created.

* `disk_size_in_gb` - (Optional) The size of the data disk which should be created.  Changing this forces a new resource to be created.

-> **Note:** Required if `create_option` is specified as `Empty`.

* `lun` - (Optional) The logical unit number of the data disk, which must be unique within the virtual machine. Changing this forces a new resource to be created.

-> **Note:** Required if `create_option` is specified as `Empty`.

* `storage_account_type` - (Optional)  The type of storage account which should back the data disk. Possible values include `Premium_LRS`, `PremiumV2_LRS`, `Premium_ZRS`, `Standard_LRS`, `StandardSSD_LRS`, `StandardSSD_ZRS` and `UltraSSD_LRS`. Changing this forces a new resource to be created.

* `write_accelerator_enabled` - (Optional) Whether to enable the write accelerator on the data disk. Defaults to `false`. Changing this forces a new resource to be created.

---

A `data_disk_count` block supports the following:

* `max` - (Optional) The maximum value of data disk count.

* `min` - (Optional) The minimum value of data disk count.

---

A `extension` block supports the following:

* `name` - (Required) The name for the Compute Fleet extension. Changing this forces a new resource to be created.

* `publisher` - (Required) Specifies the Publisher of the extension. Changing this forces a new resource to be created.

* `type` - (Required) Specifies the Type of the extension. Changing this forces a new resource to be created.

* `type_handler_version` - (Required) Specifies the version of the extension to use, available versions can be found using the Azure CLI. Changing this forces a new resource to be created.

* `auto_upgrade_minor_version_enabled` - (Optional) Whether to use the latest version of the extension at deployment time if one is available. This won't auto-update the extension on existing installation. Defaults to `false`. Changing this forces a new resource to be created.

* `automatic_upgrade_enabled` - (Optional) Whether to automatically upgrade the extension if there is a newer version of the extension available. Defaults to `false`. Changing this forces a new resource to be created.

* `extensions_to_provision_after_vm_creation` - (Optional) Specifies an ordered list of extension names which Compute Fleet should provision after virtual machine creation. Changing this forces a new resource to be created.

* `failure_suppression_enabled` - (Optional) Whether to suppress the failures from the extension. Defaults to `false`. Changing this forces a new resource to be created.

-> **Note:** Operational failures such as not connecting to the virtual machine will not be suppressed regardless of the `failure_suppression_enabled` value.

* `force_extension_execution_on_change` - (Optional) A value which, when different to the previous value can be used to force-run the extension even if the extension configuration hasn't changed. Changing this forces a new resource to be created.

* `protected_settings_from_key_vault` - (Optional) A `protected_settings_from_key_vault` block as defined below. Changing this forces a new resource to be created.

* `protected_settings_json` - (Optional) A JSON string which specifies sensitive settings (such as passwords) for the extension. Changing this forces a new resource to be created.

-> **Note:** Keys within the `protected_settings_json` block are notoriously case-sensitive, where the casing required (e.g. `TitleCase` vs `snakeCase`) depends on the extension being used. Please refer to the documentation for the specific virtual machine extension you're looking to use for more information.

* `settings_json` - (Optional) A JSON string which specifies settings for the extension. Changing this forces a new resource to be created.

---

A `gallery_application` block supports the following:

* `version_id` - (Required) Specifies the gallery application version resource ID. Changing this forces a new resource to be created.

* `automatic_upgrade_enabled` - (Optional) Whether to automatically upgrade the gallery application if there is a newer version of the gallery application available. Defaults to `false`. Changing this forces a new resource to be created.

* `configuration_blob_uri` - (Optional) Specifies the URI to an azure blob that will replace the default configuration for the package if provided. Changing this forces a new resource to be created.

* `order` - (Optional) Specifies the order in which the packages have to be installed. Defaults to `0`. Changing this forces a new resource to be created.

* `tag` - (Optional) Specifies a passthrough value for more generic context. Changing this forces a new resource to be created.

* `treat_failure_as_deployment_failure_enabled` - (Optional) Whether to treat any failure for any operation in the virtual machine application as a deployment failure. Defaults to `false`. Changing this forces a new resource to be created.

---

A `identity` block supports the following:

* `type` - (Required)  The type of managed identity that should be configured on the Compute Fleet. The only possible value is `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of user managed identity IDs to be assigned to the Compute Fleet.

---

A `ip_configuration` block supports the following:

* `name` - (Required) The Name which should be used for the IP configuration. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet which the IP configuration should be connected to. Changing this forces a new resource to be created.

* `application_gateway_backend_address_pool_ids` - (Optional) Specifies a list of backend address pools IDs from an application gateway which the Compute Fleet should be connected to. Changing this forces a new resource to be created.

* `application_security_group_ids` - (Optional) Specifies a list of application security group IDs which the Compute Fleet should be connected to. Changing this forces a new resource to be created.

* `load_balancer_backend_address_pool_ids` - (Optional) Specifies a list of backend address pools IDs from a load balancer which the Compute Fleet should be connected to. Changing this forces a new resource to be created.

* `primary_ip_configuration_enabled` - (Optional) Whether to set it as the primary IP configuration. Defaults to `false`. Changing this forces a new resource to be created.

* `public_ip_address` - (Optional) One or more `public_ip_address` blocks as defined below. Changing this forces a new resource to be created.

* `version` - (Optional) The internet protocol version which should be used for the IP configuration. Possible values are `IPv4` and `IPv6`. Defaults to `IPv4`. Changing this forces a new resource to be created.

---

A `ip_tag` block supports the following:

* `tag` - (Required) The IP tag associated with the public IP. Changing this forces a new resource to be created.

* `type` - (Required) The type of IP tag. Changing this forces a new resource to be created.

---

A `linux_configuration` block supports the following:

* `admin_username` - (Required) Specifies the name of the administrator account. Changing this forces a new resource to be created.

* `computer_name_prefix` - (Required) Specifies the computer name prefix for all the linux virtual machines in the Compute Fleet. Changing this forces a new resource to be created.

* `admin_password` - (Optional) Specifies the password of the administrator account. Changing this forces a new resource to be created.

* `admin_ssh_keys` - (Optional) Specifies a list of the public key which should be used for authentication, which needs to be in `ssh-rsa` format with at least 2048-bit or in `ssh-ed25519` format. Changing this forces a new resource to be created.

* `bypass_platform_safety_checks_enabled` - (Optional) Whether to bypass platform safety checks. Defaults to `false`. Changing this forces a new resource to be created.

* `password_authentication_enabled` - (Optional) Whether to enable the password authentication. Defaults to `false`. Changing this forces a new resource to be created.

-> **Note:** When an `admin_password` is specified `password_authentication_enabled` must be set to `true`. 

* `patch_assessment_mode` - (Optional) Specifies the mode of virtual machine Guest Patching for the virtual machines that are associated to the Compute Fleet. The only possible value is `ImageDefault`. Changing this forces a new resource to be created.

* `patch_mode` - (Optional)  Specifies the mode of in-guest patching of the virtual machines. Possible values are `AutomaticByPlatform` and `ImageDefault`. Changing this forces a new resource to be created.

* `provision_vm_agent_enabled` - (Optional) Whether to provision the virtual machine agent on each virtual machine in the Scale Set. Defaults to `true`. Changing this forces a new resource to be created.

* `reboot_setting` - (Optional) Specifies the reboot setting for all `AutomaticByPlatform` patch installation operations. Possible values are `Always`, `IfRequired`, `Never` and `Unknown`. Changing this forces a new resource to be created.

* `secret` - (Optional) One or more `secret` blocks as defined below. Changing this forces a new resource to be created.

* `vm_agent_platform_updates_enabled` - (Optional) Whether to enable the virtual machine agent platform updates for the linux virtual machine in the Compute Fleet. Defaults to `false`. Changing this forces a new resource to be created.

---

A `local_storage_in_gib` block supports the following:

* `max` - (Optional) The maximum value of local storage in GiB.

* `min` - (Optional) The minimum value of local storage in GiB.

---

A `memory_in_gib` block supports the following:

* `max` - (Optional) The maximum value of memory in GiB.

* `min` - (Optional) The minimum value of memory in GiB.

---

A `memory_in_gib_per_vcpu` block supports the following:

* `max` - (Optional) The maximum value of memory per vCPU in GiB.

* `min` - (Optional) The minimum value of memory per vCPU in GiB.

---

A `network_bandwidth_in_mbps` block supports the following:

* `max` - (Optional) The maximum value of network bandwidth in Mbps.

* `min` - (Optional) The minimum value of network bandwidth in Mbps.

---


A `network_interface_count` block supports the following:

* `max` - (Optional) The maximum value of network interface count.

* `min` - (Optional) The minimum value of network interface count.

*

---

A `os_disk` block supports the following:

* `caching` - (Optional) The Type of caching which should be used for the internal OS Disk. Possible values are `ReadOnly` and `ReadWrite`. Changing this forces a new resource to be created.

* `delete_option` - (Optional) Specify what happens to the os disk when the virtual machine is deleted. Possible values are `Delete` and `Detach`. Changing this forces a new resource to be created.

* `diff_disk_option` - (Optional) Specifies the Ephemeral Disk Settings for the OS Disk. The only possible value is `Local`. Changing this forces a new resource to be created.

* `diff_disk_placement` - (Optional) Specifies where to store the Ephemeral Disk. Possible values are `CacheDisk`, `NvmeDisk` and `ResourceDisk`. Changing this forces a new resource to be created.

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to encrypt the OS Disk. Changing this forces a new resource to be created.

* `disk_size_in_gb` - (Optional) The size of the internal OS Disk in GB, if you wish to vary from the size used in the image the Compute Fleet is sourced from. Changing this forces a new resource to be created.

* `security_encryption_type` - (Optional) Specifies the encryption type of the OS Disk. Possible values are `DiskWithVMGuestState`, `NonPersistedTPM` and `VMGuestStateOnly`. Changing this forces a new resource to be created.

* `storage_account_type` - (Optional) The Type of Storage Account which should back the OS Disk. Possible values include `Premium_LRS`, `Premium_ZRS`, `Standard_LRS`, `StandardSSD_LRS` and `StandardSSD_ZRS`. Changing this forces a new resource to be created.

* `write_accelerator_enabled` - (Optional) Whether to enable the write accelerator on the OS Disk. Defaults to `false`. Changing this forces a new resource to be created.

---

A `plan` block supports the following:

* `name` - (Required) Specifies the name of the image from the marketplace. Changing this forces a new resource to be created.

* `product` - (Required) Specifies the product of the image from the marketplace. Changing this forces a new resource to be created.

* `publisher` - (Required) Specifies the publisher of the image. Changing this forces a new resource to be created.

* `promotion_code` - (Optional) Specifies the promotion code of the image from the marketplace.

---

A `protected_settings_from_key_vault` block supports the following:

* `secret_url` - (Required) The URL to the key vault secret which stores the protected settings. Changing this forces a new resource to be created.

* `source_vault_id` - (Required) The ID of the source key vault. Changing this forces a new resource to be created.

---

A `public_ip_address` block supports the following:

* `name` - (Required) The name of the Public IP Address Configuration. Changing this forces a new resource to be created.

* `delete_option` - (Optional) Specify what happens to the public ip address when the virtual machine is deleted. Possible values are `Delete` and `Detach`. Changing this forces a new resource to be created.

* `domain_name_label` - (Optional) The prefix which should be used for the domain name label for each virtual machine. Azure concatenates the domain name label and virtual machine index to create a unique domain name label for each virtual machine. Changing this forces a new resource to be created.

* `domain_name_label_scope` - (Optional) The domain name label scope. Possible values are `NoReuse`, `ResourceGroupReuse`, `SubscriptionReuse` and `TenantReuse`. Changing this forces a new resource to be created.

* `idle_timeout_in_minutes` - (Optional) The idle timeout in minutes for the public IP address. Possible values are in the range `4` to `32`. Changing this forces a new resource to be created.

* `ip_tag` - (Optional) One or more `ip_tag` blocks as defined above. Changing this forces a new resource to be created.

* `public_ip_prefix_id` - (Optional) The ID of the public IP address prefix from where public IP addresses should be allocated. Changing this forces a new resource to be created.

* `sku_name` - (Optional) Specifies what public IP address SKU the public IP address should be provisioned as. Possible values are `Standard_Regional` and `Standard_Global`. Changing this forces a new resource to be created.

* `version` - (Optional) The internet protocol version which should be used for the public IP address. Possible values are `IPv4` and `IPv6`. Defaults to `IPv4`. Changing this forces a new resource to be created.

---

A `rdma_network_interface_count` block supports the following:

* `max` - (Optional) The maximum value of RDMA (Remote Direct Memory Access) network interface count.

* `min` - (Optional) The minimum value of RDMA (Remote Direct Memory Access) network interface count.

---

A `regular_priority_profile` block supports the following:

* `allocation_strategy` - (Optional) Specifies the allocation strategy for the Compute Fleet on which the standard virtual machines will be allocated. Defaults to `LowestPrice`. Possible values are `LowestPrice` and `Prioritized`. Changing this forces a new resource to be created.

* `capacity` - (Optional) The total number of the standard virtual machines in the Compute Fleet.

* `min_capacity` - (Optional) The minimum number of standard virtual machines in the Compute Fleet. Changing this forces a new resource to be created.

---

A `secret` block supports the following:

* `certificate` - (Required) One or more `certificate` blocks as defined above. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of the key vault from which all secrets should be sourced. Changing this forces a new resource to be created.

---

A `source_image_reference` block supports the following:

* `offer` - (Required) Specifies the offer of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `publisher` - (Required) Specifies the publisher of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `sku` - (Required) Specifies the SKU of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `version` - (Required) Specifies the version of the image used to create the virtual machines. Changing this forces a new resource to be created.

---

A `spot_priority_profile` block supports the following:

* `allocation_strategy` - (Optional) Specifies the allocation strategy for the Compute Fleet on which the Azure spot virtual machines will be allocated. Defaults to `PriceCapacityOptimized`. Possible values are `LowestPrice`, `PriceCapacityOptimized`, `CapacityOptimized`. Changing this forces a new resource to be created.

* `eviction_policy` - (Optional) The policy which should be used by spot virtual machines that are evicted from the Compute Fleet. Defaults to `Delete`. Possible values are `Deallocate` and `Delete`. Changing this forces a new resource to be created.

* `maintain_enabled` - (Optional) Whether to enable the continuous goal seeking for the desired capacity and restoration of evicted spot virtual machines. Defaults to `true`. Changing this forces a new resource to be created.

* `max_hourly_price_per_vm` - (Optional) The maximum price per hour of each spot virtual machine. Defaults to `-1`. Changing this forces a new resource to be created.

* `min_capacity` - (Optional) The minimum number of spot virtual machines in the Compute Fleet. Changing this forces a new resource to be created.

* `capacity` - (Optional) The total number of the spot virtual machines in the Compute Fleet.

---

A `vcpu_count` block supports the following:

* `max` - (Optional) The maximum value of vCPU count.

* `min` - (Optional) The minimum value of vCPU count.

---

A `vm_attributes` block supports the following:

* `memory_in_gib` - (Required) A `memory_in_gib` block as defined above.

* `vcpu_count` - (Required) A `vcpu_count` block as defined above.

* `accelerator_count` - (Optional) A `accelerator_count` block as defined above.

-> **Note:** Once the `accelerator_count` has been specified, removing it forces a new resource to be created.

* `accelerator_manufacturers` - (Optional) Specifies a list of the accelerator manufacturers. Possible values are `AMD`, `Nvidia` and `Xilinx`.

-> **Note:** Once the `accelerator_manufacturers` has been specified, removing it forces a new resource to be created.

* `accelerator_support` - (Optional) Specifies whether the virtual machine sizes supporting accelerator be used to build the Compute Fleet. Defaults to `Excluded`. Possible values are `Excluded`, `Included` and `Required`.

-> **Note:** Once the `accelerator_support` has been specified, removing it forces a new resource to be created.

* `accelerator_types` - (Optional) Specifies a list of the accelerator types. Possible values are `FPGA` and `GPU`.

-> **Note:** Once the `accelerator_types` has been specified, removing it forces a new resource to be created.

* `architecture_types` - (Optional) Specifies a list of the architecture types. Possible values are `ARM64` and `X64`.

-> **Note:** Once the `architecture_types` has been specified, removing it forces a new resource to be created.

* `burstable_support` - (Optional) Specifies whether the virtual machine Sizes supporting burstable capability be used to build the Compute Fleet. Defaults to `Excluded`. Possible values are `Excluded`, `Included` and `Required`.

-> **Note:** Once the `burstable_support` has been specified, removing it forces a new resource to be created.

* `cpu_manufacturers` - (Optional) Specifies a list of the virtual machine CPU manufacturers. Possible values are `AMD`, `Ampere`, `Intel` and `Microsoft`.

-> **Note:** Once the `cpu_manufacturers` has been specified, removing it forces a new resource to be created.

* `data_disk_count` - (Optional) A `data_disk_count` block as defined above.

-> **Note:** Once the `data_disk_count` has been specified, removing it forces a new resource to be created.

* `excluded_vm_sizes` - (Optional) Specifies a list of excluded virtual machine sizes. Conflicts with `vm_sizes_profile`.

-> **Note:** Once the `excluded_vm_sizes` has been specified, removing it forces a new resource to be created.

* `local_storage_disk_types` - (Optional) Specifies a list of the local storage disk types supported by virtual machines. Possible values are `HDD` and `SSD`.

-> **Note:** Once the `local_storage_disk_types` has been specified, removing it forces a new resource to be created.

* `local_storage_in_gib` - (Optional) A `local_storage_in_gib` block as defined above.

-> **Note:** Once the `local_storage_in_gib` has been specified, removing it forces a new resource to be created.

* `local_storage_support` - (Optional) Specifies whether the virtual machine Sizes supporting local storage be used to build the Compute Fleet. Defaults to `Included`. Possible values are `Excluded`, `Included` and `Required`.

-> **Note:** Once the `local_storage_support` has been specified, removing it forces a new resource to be created.

* `memory_in_gib_per_vcpu` - (Optional) A `memory_in_gib_per_vcpu` block as defined above.

-> **Note:** Once the `memory_in_gib_per_vcpu` has been specified, removing it forces a new resource to be created.

* `network_bandwidth_in_mbps` - (Optional) A `network_bandwidth_in_mbps` block as defined above.

-> **Note:** Once the `network_bandwidth_in_mbps` has been specified, removing it forces a new resource to be created.

* `network_interface_count` - (Optional) A `network_interface_count` block as defined above.

-> **Note:** Once the `network_interface_count` has been specified, removing it forces a new resource to be created.

* `rdma_network_interface_count` - (Optional) A `rdma_network_interface_count` block as defined above.

-> **Note:** Once the `rdma_network_interface_count` has been specified, removing it forces a new resource to be created.

* `rdma_support` - (Optional) Specifies whether the virtual machine Sizes supporting RDMA (Remote Direct Memory Access) be used to build the Compute Fleet. Defaults to `Excluded`. Possible values are `Excluded`, `Included` and `Required`.

-> **Note:** Once the `rdma_support` has been specified, removing it forces a new resource to be created.

* `vm_categories` - (Optional) Specifies a list of the virtual machine categories. Possible values are `ComputeOptimized`, `FpgaAccelerated`, `GeneralPurpose`, `GpuAccelerated`, `HighPerformanceCompute`, `MemoryOptimized` and `StorageOptimized`.

-> **Note:** Once the `vm_categories` has been specified, removing it forces a new resource to be created.

---

A `vm_sizes_profile` block supports the following:

* `name` - (Required) The name of the virtual machine size.

* `rank` - (Optional) The rank of the virtual machine size.

---

A `windows_configuration` block supports the following:

* `admin_username` - (Required) Specifies the name of the administrator account. Changing this forces a new resource to be created.

* `admin_password` - (Required) Specifies the password of the administrator account. Changing this forces a new resource to be created.

* `computer_name_prefix` - (Required) Specifies the computer name prefix for all the windows virtual machines in the Compute Fleet. Changing this forces a new resource to be created.

* `additional_unattend_content` - (Optional) One or more `additional_unattend_content` blocks as defined above. Changing this forces a new resource to be created.

* `automatic_updates_enabled` - (Optional) Whether to enable the automatic updates of the virtual machines. Defaults to `true`. Changing this forces a new resource to be created.

* `bypass_platform_safety_checks_enabled` - (Optional) Whether to bypass platform safety checks. Defaults to `false`. Changing this forces a new resource to be created.

* `hot_patching_enabled` - (Optional) Whether to enable the customers to patch the virtual machines without requiring a reboot. Defaults to `false`. Changing this forces a new resource to be created.

* `patch_assessment_mode` - (Optional) Specifies the mode of virtual machine Guest Patching for the virtual machines that are associated to the Compute Fleet. The only possible value is `ImageDefault`. Changing this forces a new resource to be created.

* `patch_mode` - (Optional)  Specifies the mode of in-guest patching of the virtual machines. Possible values are `AutomaticByOS`, `AutomaticByPlatform` and `Manual`. Changing this forces a new resource to be created.

* `provision_vm_agent_enabled` - (Optional) Whether to provision the virtual machine agent on each virtual machine in the Scale Set. Defaults to `true`. Changing this forces a new resource to be created.

* `reboot_setting` - (Optional) Specifies the reboot setting for all `AutomaticByPlatform` patch installation operations. Possible values are `Always`, `IfRequired`, `Never` and `Unknown`. Changing this forces a new resource to be created.

* `secret` - (Optional) One or more `secret` blocks as defined above. Changing this forces a new resource to be created.

* `time_zone` - (Optional) Specifies the time zone of the windows virtual machine. Changing this forces a new resource to be created.

* `vm_agent_platform_updates_enabled` - (Optional) Whether to enable the virtual machine agent platform updates for the windows virtual machine in the Compute Fleet. Defaults to `false`. Changing this forces a new resource to be created.

* `winrm_listener` - (Optional) One or more `winrm_listener` blocks as defined below.  Changing this forces a new resource to be created.

---

A `winrm_listener` block supports the following:

* `protocol` - (Required) Specifies the protocol of listener. Possible values are `Http` or `Https`. Changing this forces a new resource to be created.

* `certificate_url` - (Optional) The secret URL of a key vault certificate, which must be specified when protocol is set to `Https`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Compute Fleet.

* `unique_id` - The Unique ID for the Compute Fleet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Compute Fleet.
* `read` - (Defaults to 5 minutes) Used when retrieving the Compute Fleet.
* `update` - (Defaults to 30 minutes) Used when updating the Compute Fleet.
* `delete` - (Defaults to 30 minutes) Used when deleting the Compute Fleet.

## Import

Compute Fleets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_compute_fleet.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureFleet/fleets/fleetName
```
