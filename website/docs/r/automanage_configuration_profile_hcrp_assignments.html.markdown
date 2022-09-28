---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automanage_configuration_profile_hcrp_assignment"
description: |-
  Manages a automanage ConfigurationProfileHCRPAssignment.
---

# azurerm_automanage_configuration_profile_hcrp_assignment

Manages a automanage ConfigurationProfileHCRPAssignment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-automanage"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-VN"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-sub"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "example" {
  name                = "examplepublicip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "examplesource" {
  name                = "acctnicsource"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "exampleconfigurationsource"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.example.id
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesads"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "acctexample-sc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "example" {
  name                  = "example-vm"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2012-Datacenter"
    version   = "laexample"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.example.primary_blob_endpoint}${azurerm_storage_container.example.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "winhost01"
    admin_username = "exampleadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    timezone = "Pacific Standard Time"
  }
}

resource "azurerm_automanage_configuration_profile_assignment" "example" {
  name = "example-configurationprofileassignment"
  resource_group_name = azurerm_resource_group.example.name
  vm_name = azurerm_virtual_machine.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this automanage ConfigurationProfileHCRPAssignment. Changing this forces a new automanage ConfigurationProfileHCRPAssignment to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfileHCRPAssignment should exist. Changing this forces a new automanage ConfigurationProfileHCRPAssignment to be created.

* `machine_name` - (Required) The name of the Arc machine. Changing this forces a new automanage ConfigurationProfileHCRPAssignment to be created.

---

* `configuration_profile` - (Optional) The Automanage configurationProfile ARM Resource URI.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfileHCRPAssignment.

* `managed_by` - Azure resource id. Indicates if this resource is managed by another Azure resource.

* `target_id` - The ID of the target.

* `type` - The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts".

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the automanage ConfigurationProfileHCRPAssignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfileHCRPAssignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the automanage ConfigurationProfileHCRPAssignment.

## Import

automanage ConfigurationProfileHCRPAssignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automanage_configuration_profile_hcrp_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1
```