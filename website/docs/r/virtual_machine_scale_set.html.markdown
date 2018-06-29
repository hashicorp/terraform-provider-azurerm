---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_scale_set"
sidebar_current: "docs-azurerm-resource-compute-virtualmachine-scale-set"
description: |-
  Create a Virtual Machine scale set.
---

# azurerm\_virtual\_machine\_scale\_set

Create a virtual machine scale set.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage with Managed Disks (Recommended)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "West US 2"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "test"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
  domain_name_label            = "${azurerm_resource_group.test.name}"

  tags {
    environment = "staging"
  }
}

resource "azurerm_lb" "test" {
  name                = "test"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "bpepool" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "BackEndAddressPool"
}

resource "azurerm_lb_nat_pool" "lbnatpool" {
  count                          = 3
  resource_group_name            = "${azurerm_resource_group.test.name}"
  name                           = "ssh"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  protocol                       = "Tcp"
  frontend_port_start            = 50000
  frontend_port_end              = 50119
  backend_port                   = 22
  frontend_ip_configuration_name = "PublicIPAddress"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "mytestscaleset-1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = 2
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_data_disk {
    lun 		   = 0
    caching        = "ReadWrite"
    create_option  = "Empty"
    disk_size_gb   = 10
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/myadmin/.ssh/authorized_keys"
      key_data = "${file("~/.ssh/demo_key.pub")}"
    }
  }

  network_profile {
    name    = "terraformnetworkprofile"
    primary = true

    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.bpepool.id}"]
      load_balancer_inbound_nat_rules_ids    = ["${element(azurerm_lb_nat_pool.lbnatpool.*.id, count.index)}"]
    }
  }

  tags {
    environment = "staging"
  }
}
```

## Example Usage with Unmanaged Disks

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "West US"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn"
  address_space       = ["10.0.0.0/16"]
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "westus"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "mytestscaleset-1"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/myadmin/.ssh/authorized_keys"
      key_data = "${file("~/.ssh/demo_key.pub")}"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the virtual machine scale set resource. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the virtual machine scale set. Changing this forces a new resource to be created.
* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.
* `sku` - (Required) A sku block as documented below.
* `upgrade_policy_mode` - (Required) Specifies the mode of an upgrade to virtual machines in the scale set. Possible values, `Manual` or `Automatic`.
* `overprovision` - (Optional) Specifies whether the virtual machine scale set should be overprovisioned. Defaults to `true`.
* `single_placement_group` - (Optional) Specifies whether the scale set is limited to a single placement group with a maximum size of 100 virtual machines. If set to false, managed disks must be used. Defaults to `true`. Changing this forces a
    new resource to be created. See [documentation](http://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-placement-groups) for more information.
* `license_type` - (Optional, when a Windows machine) Specifies the Windows OS license type. If supplied, the only allowed values are `Windows_Client` and `Windows_Server`.
* `os_profile` - (Required) A Virtual Machine OS Profile block as documented below.
* `os_profile_secrets` - (Optional) A collection of Secret blocks as documented below.
* `os_profile_windows_config` - (Required, when a windows machine) A Windows config block as documented below.
* `os_profile_linux_config` - (Required, when a linux machine) A Linux config block as documented below.
* `network_profile` - (Required) A collection of network profile block as documented below.
* `storage_profile_os_disk` - (Required) A storage profile os disk block as documented below
* `storage_profile_data_disk` - (Optional) A storage profile data disk block as documented below
* `storage_profile_image_reference` - (Optional) A storage profile image reference block as documented below.
* `extension` - (Optional) Can be specified multiple times to add extension profiles to the scale set. Each `extension` block supports the fields documented below.
* `boot_diagnostics` - (Optional) A boot diagnostics profile block as referenced below.
* `plan` - (Optional) A plan block as documented below.
* `priority` - (Optional) Specifies the priority for the virtual machines in the scale set, defaults to `Regular`. Possible values are `Low` and `Regular`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `zones` - (Optional) A collection of availability zones to spread the Virtual Machines over.

-> **Please Note**: Availability Zones are [in Preview and only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview) - as such you must be opted into the Preview to use this functionality. You can [opt into the Availability Zones Preview in the Azure Portal](http://aka.ms/azenroll).

`sku` supports the following:

* `name` - (Required) Specifies the size of virtual machines in a scale set.
* `tier` - (Optional) Specifies the tier of virtual machines in a scale set. Possible values, `standard` or `basic`.
* `capacity` - (Required) Specifies the number of virtual machines in the scale set.

`identity` supports the following:

* `type` - (Required) Specifies the identity type to be assigned to the scale set. The only allowable value is `SystemAssigned`. To enable Managed Service Identity (MSI) on all machines in the scale set, an extension with the type "ManagedIdentityExtensionForWindows" or "ManagedIdentityExtensionForLinux" must also be added. The scale set's Service Principal ID (SPN) can be retrieved after the scale set has been created.

```hcl
resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "vm-scaleset"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "${var.vm_sku}"
    tier     = "Standard"
    capacity = "${var.instance_count}"
  }

  identity {
    type = "systemAssigned"
  }

  extension {
    name                 = "MSILinuxExtension"
    publisher            = "Microsoft.ManagedIdentity"
    type                 = "ManagedIdentityExtensionForLinux"
    type_handler_version = "1.0"
    settings             = "{\"port\": 50342}"
  }

  output "principal_id" {
    value = "${lookup(azurerm_virtual_machine.test.identity[0], "principal_id")}"
  }
```

`os_profile` supports the following:

* `computer_name_prefix` - (Required) Specifies the computer name prefix for all of the virtual machines in the scale set. Computer name prefixes must be 1 to 9 characters long for windows images and 1 - 58 for linux. Changing this forces a new resource to be created.
* `admin_username` - (Required) Specifies the administrator account name to use for all the instances of virtual machines in the scale set.
* `admin_password` - (Required) Specifies the administrator password to use for all the instances of virtual machines in a scale set.
* `custom_data` - (Optional) Specifies custom data to supply to the machine. On linux-based systems, this can be used as a cloud-init script. On other systems, this will be copied as a file on disk. Internally, Terraform will base64 encode this value before sending it to the API. The maximum length of the binary array is 65535 bytes.

`os_profile_secrets` supports the following:

* `source_vault_id` - (Required) Specifies the key vault to use.
* `vault_certificates` - (Required, on windows machines) A collection of Vault Certificates as documented below

`vault_certificates` support the following:

* `certificate_url` - (Required) It is the Base64 encoding of a JSON Object that which is encoded in UTF-8 of which the contents need to be `data`, `dataType` and `password`.
* `certificate_store` - (Required, on windows machines) Specifies the certificate store on the Virtual Machine where the certificate should be added to.


`os_profile_windows_config` supports the following:

* `provision_vm_agent` - (Optional) Indicates whether virtual machine agent should be provisioned on the virtual machines in the scale set.
* `enable_automatic_upgrades` - (Optional) Indicates whether virtual machines in the scale set are enabled for automatic updates.
* `winrm` - (Optional) A collection of WinRM configuration blocks as documented below.
* `additional_unattend_config` - (Optional) An Additional Unattended Config block as documented below.

`winrm` supports the following:

* `protocol` - (Required) Specifies the protocol of listener
* `certificate_url` - (Optional) Specifies URL of the certificate with which new Virtual Machines is provisioned.

`additional_unattend_config` supports the following:

* `pass` - (Required) Specifies the name of the pass that the content applies to. The only allowable value is `oobeSystem`.
* `component` - (Required) Specifies the name of the component to configure with the added content. The only allowable value is `Microsoft-Windows-Shell-Setup`.
* `setting_name` - (Required) Specifies the name of the setting to which the content applies. Possible values are: `FirstLogonCommands` and `AutoLogon`.
* `content` - (Optional) Specifies the base-64 encoded XML formatted content that is added to the unattend.xml file for the specified path and component.

`os_profile_linux_config` supports the following:

* `disable_password_authentication` - (Required) Specifies whether password authentication should be disabled. Changing this forces a new resource to be created.
* `ssh_keys` - (Optional) Specifies a collection of `path` and `key_data` to be placed on the virtual machine.

~> _**Note:** Please note that the only allowed `path` is `/home/<username>/.ssh/authorized_keys` due to a limitation of Azure_

`network_profile` supports the following:

* `name` - (Required) Specifies the name of the network interface configuration.
* `primary` - (Required) Indicates whether network interfaces created from the network interface configuration will be the primary NIC of the VM.
* `ip_configuration` - (Required) An ip_configuration block as documented below.
* `accelerated_networking` - (Optional) Specifies whether to enable accelerated networking or not. Defaults to `false`.
* `dns_settings` - (Optional) An dns_settings block as documented below.
* `network_security_group_id` - (Optional) Specifies the identifier for the network security group.

`dns_settings` supports the following:

* `dns_servers` - (Required) Specifies an array of dns servers.

`ip_configuration` supports the following:

* `name` - (Required) Specifies name of the IP configuration.
* `subnet_id` - (Required) Specifies the identifier of the subnet.
* `application_gateway_backend_address_pool_ids` - (Optional) Specifies an array of references to backend address pools of application gateways. A scale set can reference backend address pools of one application gateway. Multiple scale sets cannot use the same application gateway.
* `load_balancer_backend_address_pool_ids` - (Optional) Specifies an array of references to backend address pools of load balancers. A scale set can reference backend address pools of one public and one internal load balancer. Multiple scale sets cannot use the same load balancer.
* `load_balancer_inbound_nat_rules_ids` - (Optional) Specifies an array of references to inbound NAT rules for load balancers.
* `primary` - (Optional) Specifies if this ip_configuration is the primary one.
* `public_ip_address_configuration` - (Optional) describes a virtual machines scale set IP Configuration's
 PublicIPAddress configuration. The public_ip_address_configuration is documented below.

`public_ip_address_configuration` supports the following:

* `name` - (Required) The name of the public ip address configuration
* `idle_timeout` - (Required) The idle timeout in minutes. This value must be between 4 and 32.
* `domain_name_label` - (Required) The domain name label for the dns settings.

`storage_profile_os_disk` supports the following:

* `name` - (Optional) Specifies the disk name. Must be specified when using unmanaged disk ('managed_disk_type' property not set).
* `vhd_containers` - (Optional) Specifies the vhd uri. Cannot be used when `image` or `managed_disk_type` is specified.
* `managed_disk_type` - (Optional) Specifies the type of managed disk to create. Value you must be either `Standard_LRS` or `Premium_LRS`. Cannot be used when `vhd_containers` or `image` is specified.
* `create_option` - (Required) Specifies how the virtual machine should be created. The only possible option is `FromImage`.
* `caching` - (Optional) Specifies the caching requirements. Possible values include: `None` (default), `ReadOnly`, `ReadWrite`.
* `image` - (Optional) Specifies the blob uri for user image. A virtual machine scale set creates an os disk in the same container as the user image.
                       Updating the osDisk image causes the existing disk to be deleted and a new one created with the new image. If the VM scale set is in Manual upgrade mode then the virtual machines are not updated until they have manualUpgrade applied to them.
                       When setting this field `os_type` needs to be specified. Cannot be used when `vhd_containers`, `managed_disk_type` or `storage_profile_image_reference ` are specified.
* `os_type` - (Optional) Specifies the operating system Type, valid values are windows, linux.

`storage_profile_data_disk` supports the following:

* `lun` - (Required) Specifies the Logical Unit Number of the disk in each virtual machine in the scale set.
* `create_option` - (Optional) Specifies how the data disk should be created. The only possible options are `FromImage` and `Empty`.
* `caching` - (Optional) Specifies the caching requirements. Possible values include: `None` (default), `ReadOnly`, `ReadWrite`.
* `disk_size_gb` - (Optional) Specifies the size of the disk in GB. This element is required when creating an empty disk.
* `managed_disk_type` - (Optional) Specifies the type of managed disk to create. Value must be either `Standard_LRS` or `Premium_LRS`.

`storage_profile_image_reference` supports the following:

* `id` - (Optional) Specifies the ID of the (custom) image to use to create the virtual
machine scale set, as in the [example below](#example-of-storage_profile_image_reference-with-id).
* `publisher` - (Optional) Specifies the publisher of the image used to create the virtual machines.
* `offer` - (Optional) Specifies the offer of the image used to create the virtual machines.
* `sku` - (Optional) Specifies the SKU of the image used to create the virtual machines.
* `version` - (Optional) Specifies the version of the image used to create the virtual machines.

`boot_diagnostics` supports the following:

* `enabled`: (Required) Whether to enable boot diagnostics for the virtual machine.
* `storage_uri`: (Required) Blob endpoint for the storage account to hold the virtual machine's diagnostic files. This must be the root of a storage account, and not a storage container.


`extension` supports the following:

* `name` - (Required) Specifies the name of the extension.
* `publisher` - (Required) The publisher of the extension, available publishers can be found by using the Azure CLI.
* `type` - (Required) The type of extension, available types for a publisher can be found using the Azure CLI.
* `type_handler_version` - (Required) Specifies the version of the extension to use, available versions can be found using the Azure CLI.
* `auto_upgrade_minor_version` - (Optional) Specifies whether or not to use the latest minor version available.
* `settings` - (Required) The settings passed to the extension, these are specified as a JSON object in a string.
* `protected_settings` - (Optional) The protected_settings passed to the extension, like settings, these are specified as a JSON object in a string.

`plan` supports the following:

* `name` - (Required) Specifies the name of the image from the marketplace.
* `publisher` - (Required) Specifies the publisher of the image.
* `product` - (Required) Specifies the product of the image from the marketplace.

## Example of storage_profile_image_reference with id

```hcl

resource "azurerm_image" "test" {
	name = "test"
  ...
}

resource "azurerm_virtual_machine_scale_set" "test" {
	name = "test"
  ...

	storage_profile_image_reference {
		id = "${azurerm_image.test.id}"
	}

...
```

## Attributes Reference

The following attributes are exported:

* `id` - The virtual machine scale set ID.

## Import

Virtual Machine Scale Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_scale_set.scaleset1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1
```
