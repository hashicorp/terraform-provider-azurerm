---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine"
description: |-
  Manages a Virtual Machine.
---

# azurerm_virtual_machine

Manages a Virtual Machine.

## Disclaimers

-> **Note:** The `azurerm_virtual_machine` resource has been superseded by the [`azurerm_linux_virtual_machine`](linux_virtual_machine.html) and [`azurerm_windows_virtual_machine`](windows_virtual_machine.html) resources. The existing `azurerm_virtual_machine` resource will continue to be available throughout the 2.x releases however is in a feature-frozen state to maintain compatibility - new functionality will instead be added to the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` resources.

~> **Note:** Data Disks can be attached either directly on the `azurerm_virtual_machine` resource, or using the `azurerm_virtual_machine_data_disk_attachment` resource - but the two cannot be used together. If both are used against the same Virtual Machine, spurious changes will occur.

## Example Usage (from an Azure Platform Image)

This example provisions a Virtual Machine with Managed Disks. Other examples of the `azurerm_virtual_machine` resource can be found in [the `./examples/virtual-machines` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/virtual-machines)

```hcl
variable "prefix" {
  default = "tfvmex"
}

resource "azurerm_resource_group" "main" {
  name     = "${var.prefix}-resources"
  location = "West US 2"
}

resource "azurerm_virtual_network" "main" {
  name                = "${var.prefix}-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "main" {
  name                = "${var.prefix}-nic"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.internal.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "main" {
  name                  = "${var.prefix}-vm"
  location              = azurerm_resource_group.main.location
  resource_group_name   = azurerm_resource_group.main.name
  network_interface_ids = [azurerm_network_interface.main.id]
  vm_size               = "Standard_DS1_v2"

  # Uncomment this line to delete the OS disk automatically when deleting the VM
  # delete_os_disk_on_termination = true

  # Uncomment this line to delete the data disks automatically when deleting the VM
  # delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }
  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }
  os_profile_linux_config {
    disable_password_authentication = false
  }
  tags = {
    environment = "staging"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Virtual Machine. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Virtual Machine should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Virtual Machine exists. Changing this forces a new resource to be created.

* `network_interface_ids` - (Required) A list of Network Interface ID's which should be associated with the Virtual Machine.

* `os_profile_linux_config` - (Required, when a Linux machine) A `os_profile_linux_config` block.

* `os_profile_windows_config` - (Required, when a Windows machine) A `os_profile_windows_config` block.

* `vm_size` - (Required) Specifies the [size of the Virtual Machine](https://docs.microsoft.com/azure/virtual-machines/sizes-general). See also [Azure VM Naming Conventions](https://docs.microsoft.com/azure/virtual-machines/vm-naming-conventions).

---

* `availability_set_id` - (Optional) The ID of the Availability Set in which the Virtual Machine should exist. Changing this forces a new resource to be created.

* `boot_diagnostics` - (Optional) A `boot_diagnostics` block.

* `additional_capabilities` - (Optional) A `additional_capabilities` block.

* `delete_os_disk_on_termination` - (Optional) Should the OS Disk (either the Managed Disk / VHD Blob) be deleted when the Virtual Machine is destroyed? Defaults to `false`.

~> **Note:** This setting works when instance is deleted via Terraform only and don't forget to delete disks manually if you deleted VM manually. It can increase spending.

* `delete_data_disks_on_termination` - (Optional) Should the Data Disks (either the Managed Disks / VHD Blobs) be deleted when the Virtual Machine is destroyed? Defaults to `false`.

-> **Note:** This setting works when instance is deleted via Terraform only and don't forget to delete disks manually if you deleted VM manually. It can increase spending.

* `identity` - (Optional) A `identity` block.

* `license_type` - (Optional) Specifies the BYOL Type for this Virtual Machine. This is only applicable to Windows Virtual Machines. Possible values are `Windows_Client` and `Windows_Server`.

* `os_profile` - (Optional) An `os_profile` block. Required when `create_option` in the `storage_os_disk` block is set to `FromImage`.

* `os_profile_secrets` - (Optional) One or more `os_profile_secrets` blocks.

* `plan` - (Optional) A `plan` block.

* `primary_network_interface_id` - (Optional) The ID of the Network Interface (which must be attached to the Virtual Machine) which should be the Primary Network Interface for this Virtual Machine.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group to which this Virtual Machine should be assigned. Changing this forces a new resource to be created

* `storage_data_disk` - (Optional) One or more `storage_data_disk` blocks.

~> **Please Note:** Data Disks can also be attached either using this block or [the `azurerm_virtual_machine_data_disk_attachment` resource](virtual_machine_data_disk_attachment.html) - but not both.

* `storage_image_reference` - (Optional) A `storage_image_reference` block.

* `storage_os_disk` - (Required) A `storage_os_disk` block.

* `tags` - (Optional) A mapping of tags to assign to the Virtual Machine.

* `zones` - (Optional) A list of a single item of the Availability Zone which the Virtual Machine should be allocated in.

-> **Please Note**: Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview).

For more information on the different example configurations, please check out the [Azure documentation](https://docs.microsoft.com/en-gb/rest/api/compute/virtualmachines/createorupdate#examples)

---

A `additional_unattend_config` block supports the following:

* `pass` - (Required) Specifies the name of the pass that the content applies to. The only allowable value is `oobeSystem`.

* `component` - (Required) Specifies the name of the component to configure with the added content. The only allowable value is `Microsoft-Windows-Shell-Setup`.

* `setting_name` - (Required) Specifies the name of the setting to which the content applies. Possible values are: `FirstLogonCommands` and `AutoLogon`.

* `content` - (Optional) Specifies the base-64 encoded XML formatted content that is added to the unattend.xml file for the specified path and component.

---

A `boot_diagnostics` block supports the following:

* `enabled` - (Required) Should Boot Diagnostics be enabled for this Virtual Machine?

* `storage_uri` - (Required) The Storage Account's Blob Endpoint which should hold the virtual machine's diagnostic files.

~> **NOTE:** This needs to be the root of a Storage Account and not a Storage Container.

---

A `additional_capabilities` block supports the following:

* `ultra_ssd_enabled` - (Required) Should Ultra SSD disk be enabled for this Virtual Machine?

-> **Note:** Azure Ultra Disk Storage is only available in a region that support availability zones and can only enabled on the following VM series: `ESv3`, `DSv3`, `FSv3`, `LSv2`, `M` and `Mv2`. For more information see the `Azure Ultra Disk Storage` [product documentation](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/disks-enable-ultra-ssd).

---

A `identity` block supports the following:

* `type` - (Required) The Managed Service Identity Type of this Virtual Machine. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), `UserAssigned` (where you can specify the Service Principal ID's) to be used by this Virtual Machine using the `identity_ids` field, and `SystemAssigned, UserAssigned` which assigns both a system managed identity as well as the specified user assigned identities.

-> **NOTE:** Managed Service Identity previously required the installation of a VM Extension, but this information [is now available via the Azure Instance Metadata Service](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview#how-does-it-work).

~> **NOTE:** When `type` is set to `SystemAssigned`, identity the Principal ID can be retrieved after the virtual machine has been created. More details are available below. See [documentation](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview) for additional information.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned to the VM. Required if `type` is `UserAssigned`.

---

A `os_profile` block supports the following:

* `computer_name` - (Required) Specifies the name of the Virtual Machine.

* `admin_username` - (Required) Specifies the name of the local administrator account.

* `admin_password` - (Required for Windows, Optional for Linux) The password associated with the local administrator account.

-> **NOTE:** If using Linux, it may be preferable to use SSH Key authentication (available in the `os_profile_linux_config` block) instead of password authentication.

~> **NOTE:** `admin_password` must be between 6-72 characters long and must satisfy at least 3 of password complexity requirements from the following:
1. Contains an uppercase character
2. Contains a lowercase character
3. Contains a numeric digit
4. Contains a special character

* `custom_data` - (Optional) Specifies custom data to supply to the machine. On Linux-based systems, this can be used as a cloud-init script. On other systems, this will be copied as a file on disk. Internally, Terraform will base64 encode this value before sending it to the API. The maximum length of the binary array is 65535 bytes.

---

A `os_profile_linux_config` block supports the following:

* `disable_password_authentication` - (Required) Specifies whether password authentication should be disabled. If set to `false`, an `admin_password` must be specified.

* `ssh_keys` - (Optional) One or more `ssh_keys` blocks. This field is required if `disable_password_authentication` is set to `true`.

---

A `os_profile_secrets` block supports the following:

* `source_vault_id` - (Required) Specifies the ID of the Key Vault to use.

* `vault_certificates` - (Required) One or more `vault_certificates` blocks.

---

A `os_profile_windows_config` block supports the following:

* `provision_vm_agent` - (Optional) Should the Azure Virtual Machine Guest Agent be installed on this Virtual Machine? Defaults to `false`.

-> **NOTE:** This is different from the Default value used for this field within Azure.

* `enable_automatic_upgrades` - (Optional) Are automatic updates enabled on this Virtual Machine? Defaults to `false.`

* `timezone` - (Optional) Specifies the time zone of the virtual machine, [the possible values are defined here](http://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/).

* `winrm` - (Optional) One or more `winrm` block.

* `additional_unattend_config` - (Optional) A `additional_unattend_config` block.

---

A `plan` block supports the following:

* `name` - (Required) Specifies the name of the image from the marketplace.

* `publisher` - (Required) Specifies the publisher of the image.

* `product` - (Required) Specifies the product of the image from the marketplace.

---

A `ssh_keys` block supports the following:

* `key_data` - (Required) The Public SSH Key which should be written to the `path` defined above.

~> **Note:** Azure only supports RSA SSH2 key signatures of at least 2048 bits in length

-> **NOTE:** Rather than defining this in-line you can source this from a local file using [the `file` function](https://www.terraform.io/docs/configuration/functions/file.html) - for example `key_data = file("~/.ssh/id_rsa.pub")`.

* `path` - (Required) The path of the destination file on the virtual machine

-> **NOTE:** Due to a limitation in the Azure VM Agent the only allowed `path` is `/home/{username}/.ssh/authorized_keys`.


---

A `storage_image_reference` block supports the following:

This block provisions the Virtual Machine from one of two sources: an Azure Platform Image (e.g. Ubuntu/Windows Server) or a Custom Image.

To provision from an Azure Platform Image, the following fields are applicable:

* `publisher` - (Required) Specifies the publisher of the image used to create the virtual machine. Changing this forces a new resource to be created.

* `offer` - (Required) Specifies the offer of the image used to create the virtual machine. Changing this forces a new resource to be created.

* `sku` - (Required) Specifies the SKU of the image used to create the virtual machine. Changing this forces a new resource to be created.

* `version` - (Optional) Specifies the version of the image used to create the virtual machine. Changing this forces a new resource to be created.

To provision a Custom Image, the following fields are applicable:

* `id` - (Required) Specifies the ID of the Custom Image which the Virtual Machine should be created from. Changing this forces a new resource to be created.

-> **NOTE:** An example of how to use this is available within [the `./examples/virtual-machines/managed-disks/from-custom-image` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/virtual-machines/managed-disks/from-custom-image)

---

A `storage_data_disk` block supports the following:

~> **NOTE:** Data Disks can also be attached either using this block or [the `azurerm_virtual_machine_data_disk_attachment` resource](virtual_machine_data_disk_attachment.html) - but not both.

* `name` - (Required) The name of the Data Disk.

* `caching` - (Optional) Specifies the caching requirements for the Data Disk. Possible values include `None`, `ReadOnly` and `ReadWrite`.

* `create_option` - (Required) Specifies how the data disk should be created. Possible values are `Attach`, `FromImage` and `Empty`.

~> **NOTE:** If using an image that does not have data to be written to the Data Disk, use `Empty` as the create option in order to create the desired disk without any data.

* `disk_size_gb` - (Optional) Specifies the size of the data disk in gigabytes.

* `lun` - (Required) Specifies the logical unit number of the data disk. This needs to be unique within all the Data Disks on the Virtual Machine.

* `write_accelerator_enabled` - (Optional) Specifies if Write Accelerator is enabled on the disk. This can only be enabled on `Premium_LRS` managed disks with no caching and [M-Series VMs](https://docs.microsoft.com/en-us/azure/virtual-machines/workloads/sap/how-to-enable-write-accelerator). Defaults to `false`.

The following properties apply when using Managed Disks:

* `managed_disk_type` - (Optional) Specifies the type of managed disk to create. Possible values are either `Standard_LRS`, `StandardSSD_LRS`, `Premium_LRS` or `UltraSSD_LRS`.

-> **Note**: `managed_disk_type` of type `UltraSSD_LRS` is currently in preview and are not available to subscriptions that have not [requested](https://aka.ms/UltraSSDPreviewSignUp) onboarding to `Azure Ultra Disk Storage` preview. `Azure Ultra Disk Storage` is only available in `East US 2`, `North Europe`, and `Southeast Asia` regions. For more information see the `Azure Ultra Disk Storage` [product documentation](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/disks-enable-ultra-ssd), [product blog](https://azure.microsoft.com/en-us/blog/announcing-the-general-availability-of-azure-ultra-disk-storage/) and [FAQ](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/faq-for-disks#ultra-disks). You must also set `additional_capabilities.ultra_ssd_enabled` to `true`.

* `managed_disk_id` - (Optional) Specifies the ID of an Existing Managed Disk which should be attached to this Virtual Machine. When this field is set `create_option` must be set to `Attach`.

The following properties apply when using Unmanaged Disks:

* `vhd_uri` - (Optional) Specifies the URI of the VHD file backing this Unmanaged Data Disk. Changing this forces a new resource to be created.

---

A `storage_os_disk` block supports the following:

* `name` - (Required) Specifies the name of the OS Disk.

* `create_option` - (Required) Specifies how the OS Disk should be created. Possible values are `Attach` (managed disks only) and `FromImage`.

* `caching` - (Optional) Specifies the caching requirements for the OS Disk. Possible values include `None`, `ReadOnly` and `ReadWrite`.

* `disk_size_gb` - (Optional) Specifies the size of the OS Disk in gigabytes.

* `image_uri` - (Optional) Specifies the Image URI in the format `publisherName:offer:skus:version`. This field can also specify the [VHD uri](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-linux-cli-deploy-templates/#create-a-custom-vm-image) of a custom VM image to clone. When cloning a Custom (Unmanaged) Disk Image the `os_type` field must be set.

* `os_type` - (Optional) Specifies the Operating System on the OS Disk. Possible values are `Linux` and `Windows`.

* `write_accelerator_enabled` - (Optional) Specifies if Write Accelerator is enabled on the disk. This can only be enabled on `Premium_LRS` managed disks with no caching and [M-Series VMs](https://docs.microsoft.com/en-us/azure/virtual-machines/workloads/sap/how-to-enable-write-accelerator). Defaults to `false`.

The following properties apply when using Managed Disks:

* `managed_disk_id` - (Optional) Specifies the ID of an existing Managed Disk which should be attached as the OS Disk of this Virtual Machine. If this is set then the `create_option` must be set to `Attach`.

* `managed_disk_type` - (Optional) Specifies the type of Managed Disk which should be created. Possible values are `Standard_LRS`, `StandardSSD_LRS` or `Premium_LRS`.

The following properties apply when using Unmanaged Disks:

* `vhd_uri` - (Optional) Specifies the URI of the VHD file backing this Unmanaged OS Disk. Changing this forces a new resource to be created.

---

A `vault_certificates` block supports the following:

* `certificate_url` - (Required) The ID of the Key Vault Secret. Stored secret is the Base64 encoding of a JSON Object that which is encoded in UTF-8 of which the contents need to be:

```json
{
  "data":"<Base64-encoded-certificate>",
  "dataType":"pfx",
  "password":"<pfx-file-password>"
}
```

-> **NOTE:** If your certificate is stored in Azure Key Vault - this can be sourced from the `secret_id` property on the `azurerm_key_vault_certificate` resource.

* `certificate_store` - (Required, on windows machines) Specifies the certificate store on the Virtual Machine where the certificate should be added to, such as `My`.

---

A `winrm` block supports the following:

* `protocol` - (Required) Specifies the protocol of listener. Possible values are `HTTP` or `HTTPS`.

* `certificate_url` - (Optional) The ID of the Key Vault Secret which contains the encrypted Certificate which should be installed on the Virtual Machine. This certificate must also be specified in the `vault_certificates` block within the `os_profile_secrets` block.

-> **NOTE:** This can be sourced from the `secret_id` field on the `azurerm_key_vault_certificate` resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Machine.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Virtual Machine.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Virtual Machine.

-> You can access the Principal ID via `${azurerm_virtual_machine.example.identity.0.principal_id}`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Virtual Machine.
* `update` - (Defaults to 60 minutes) Used when updating the Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine.
* `delete` - (Defaults to 60 minutes) Used when deleting the Virtual Machine.

## Import

Virtual Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/virtualMachines/machine1
```
