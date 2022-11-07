---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automanage_configuration_profile_assignment"
description: |-
  Manages a automanage ConfigurationProfileAssignment.
---

# azurerm_automanage_configuration_profile_assignment

Manages a automanage ConfigurationProfileAssignment.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "example-automanage"
  location = "West Europe"
}

resource "azurerm_virtual_network" "test" {
  name                = "example-VN"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "example-sub"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes       = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "examplepublicip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "examplesads"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "accttest-sc"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "example-vm"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "winhost01"
    admin_username = "testadmin"
    admin_password = "Password1234567!"
  }

  os_profile_windows_config {
    timezone = "Pacific Standard Time"
  }
}

resource "azurerm_automanage_configuration_profile_assignment" "test" {
  name = "default"
  resource_group_name = azurerm_resource_group.test.name
  vm_name = azurerm_virtual_machine.test.name
  configuration_profile = "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this automanage ConfigurationProfileAssignment. Changing this forces a new automanage ConfigurationProfileAssignment to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfileAssignment should exist. Changing this forces a new automanage ConfigurationProfileAssignment to be created.

* `vm_name` - (Required) The name of the virtual machine. Changing this forces a new automanage ConfigurationProfileAssignment to be created.

* `configuration_profile` - (Required) The Automanage configurationProfile ARM Resource URI.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfileAssignment.

* `managed_by` - Azure resource id. Indicates if this resource is managed by another Azure resource.

* `target_id` - The ID of the target.

* `type` - The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts".

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the automanage ConfigurationProfileAssignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfileAssignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the automanage ConfigurationProfileAssignment.

## Import

automanage ConfigurationProfileAssignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automanage_configuration_profile_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1
```