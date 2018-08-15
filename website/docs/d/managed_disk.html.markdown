---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk"
sidebar_current: "docs-azurerm-datasource-managed-disk"
description: |-
  Get information about the specified managed disk.
---

# Data Source: azurerm_managed_disk

Use this data source to access the properties of an existing Azure Managed Disk.

## Example Usage

```hcl
data "azurerm_managed_disk" "datasourcemd" {
    name = "testManagedDisk"
    resource_group_name = "acctestRG"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn"
  address_space       = ["10.0.0.0/16"]
  location            = "West US 2"
  resource_group_name = "acctestRG"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = "acctestRG"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni"
  location            = "West US 2"
  resource_group_name = "acctestRG"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm"
  location              = "West US 2"
  resource_group_name   = "acctestRG"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_DS1_v2"

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

  storage_data_disk {
    name              = "datadisk_new"
    managed_disk_type = "Standard_LRS"
    create_option     = "Empty"
    lun               = 0
    disk_size_gb      = "1023"
  }

  storage_data_disk {
    name            = "${data.azurerm_managed_disk.datasourcemd.name}"
    managed_disk_id = "${data.azurerm_managed_disk.datasourcemd.id}"
    create_option   = "Attach"
    lun             = 1
    disk_size_gb    = "${data.azurerm_managed_disk.datasourcemd.disk_size_gb}"
  }

  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags {
    environment = "staging"
  }
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Managed Disk.
* `resource_group_name` - (Required) Specifies the name of the resource group.


## Attributes Reference

* `storage_account_type` - The storage account type for the managed disk.
* `source_uri` - The source URI for the managed disk
* `source_resource_id` - ID of an existing managed disk that the current resource was created from.
* `os_type` - The operating system for managed disk. Valid values are `Linux` or `Windows`
* `disk_size_gb` - The size of the managed disk in gigabytes.
* `tags` - A mapping of tags assigned to the resource.
* `zones` - (Optional) A collection containing the availability zone the managed disk is allocated in.

-> **Please Note**: Availability Zones are [in Preview and only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview) - as such you must be opted into the Preview to use this functionality. You can [opt into the Availability Zones Preview in the Azure Portal](http://aka.ms/azenroll).