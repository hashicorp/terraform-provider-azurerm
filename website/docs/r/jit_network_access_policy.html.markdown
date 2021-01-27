---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_jit_network_access_policy"
description: |-
  Manages a security JitNetworkAccessPolicy.
---

# azurerm_jit_network_access_policy

Manages a security JitNetworkAccessPolicy.

## Example Usage

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

resource "azurerm_public_ip" "example" {
  name                = "example_pip"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Dynamic"
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.example.id
  }
}

resource "azurerm_network_security_group" "example" {
  name                = "example-nsg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_interface_security_group_association" "example" {
  network_interface_id      = azurerm_network_interface.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_windows_virtual_machine" "example" {
  name                = "example-machine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  depends_on = [azurerm_network_interface_security_group_association.example]
}

resource "azurerm_jit_network_access_policy" "example" {
  name                = "example-jitnetworkaccesspolicy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  virtual_machine {
    id = azurerm_windows_virtual_machine.example.id
    ports {
      port                            = 42
      protocol                        = "*"
      allowed_source_address_prefixes = ["*"]
      max_request_access_duration     = "PT3H"
    }

    ports {
      port                            = 43
      protocol                        = "TCP"
      allowed_source_address_prefixes = ["192.168.0.3", "192.168.0.0/16"]
      max_request_access_duration     = "PT3H"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Jit Network Access Policy. Changing this forces a new security resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Jit Network Access Policy should exist. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Jit Network Access Policy should exist. Changing this forces a new resource to be created.

* `virtual_machine` - (Required) One or more `virtual_machine` block as defined below.

---

A `virtual_machine` block exports the following:

* `id` - (Required) The resource id of virtual machine.

* `ports` - (Required) One or more `ports` block as defined below.

---

A `ports` block exports the following:

* `allowed_source_address_prefixes` - (Required) A list of source Ip Address allowed for the Jit Network Access Policy. Possible values are IP address, CIDR or `*`.

* `max_request_access_duration` - (Required) Maximum duration requests can be made for the Jit Network Access Policy, represented in ISO 8601 duration format.

* `port` - (Required) The port number the Jit Network Access Policy will expose.

* `protocol` - (Required) The network protocol associated with port. Possible values are `TCP`, `UDP` and `*`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Jit Network Access Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Jit Network Access Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Jit Network Access Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Jit Network Access Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Jit Network Access Policy.

## Import

Jit Network Access Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_jit_network_access_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Security/locations/westeurope/jitNetworkAccessPolicies/policy1
```
