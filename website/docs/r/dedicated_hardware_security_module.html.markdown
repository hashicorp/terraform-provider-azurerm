---
subcategory: "Hardware Security Module"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dedicated_hardware_security_module"
description: |-
  Manages a Dedicated Hardware Security Module.
---

# azurerm_dedicated_hardware_security_module

Manages a Dedicated Hardware Security Module.

-> **Note**: Before using this resource, it's required to submit the request of registering the providers and features with Azure CLI `az provider register --namespace Microsoft.HardwareSecurityModules && az feature register --namespace Microsoft.HardwareSecurityModules --name AzureDedicatedHSM && az provider register --namespace Microsoft.Network && az feature register --namespace Microsoft.Network --name AllowBaremetalServers` and ask service team (hsmrequest@microsoft.com) to approve. See more details from https://docs.microsoft.com/en-us/azure/dedicated-hsm/tutorial-deploy-hsm-cli#prerequisites.

-> **Note**: If the quota is not enough in some region, please submit the quota request to service team.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.2.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-compute"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.2.0.0/24"]
}

resource "azurerm_subnet" "example2" {
  name                 = "example-hsmsubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.2.1.0/24"]

  delegation {
    name = "first"

    service_delegation {
      name = "Microsoft.HardwareSecurityModules/dedicatedHSMs"
      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet" "example3" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.2.255.0/26"]
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "example" {
  name                = "example-vnetgateway"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type     = "ExpressRoute"
  vpn_type = "PolicyBased"
  sku      = "Standard"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.example.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.example3.id
  }
}

resource "azurerm_dedicated_hardware_security_module" "example" {
  name                = "example-hsm"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "SafeNet Luna Network HSM A790"

  network_profile {
    network_interface_private_ip_addresses = ["10.2.1.8"]
    subnet_id                              = azurerm_subnet.example2.id
  }

  stamp_id = "stamp2"

  tags = {
    env = "Test"
  }

  depends_on = [azurerm_virtual_network_gateway.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Dedicated Hardware Security Module. Changing this forces a new Dedicated Hardware Security Module to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Dedicated Hardware Security Module should exist. Changing this forces a new Dedicated Hardware Security Module to be created.

* `location` - (Required) The Azure Region where the Dedicated Hardware Security Module should exist. Changing this forces a new Dedicated Hardware Security Module to be created.

* `network_profile` - (Required)  A `network_profile` block as defined below.

* `sku_name` - (Required) The sku name of the dedicated hardware security module. Changing this forces a new Dedicated Hardware Security Module to be created.

* `stamp_id` - (Optional) The ID of the stamp. Possible values are `stamp1` or `stamp2`. Changing this forces a new Dedicated Hardware Security Module to be created.

* `zones` - (Optional) The Dedicated Hardware Security Module zones. Changing this forces a new Dedicated Hardware Security Module to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Dedicated Hardware Security Module.

---

An `network_profile` block exports the following:

* `network_interface_private_ip_addresses` - (Required) The private IPv4 address of the network interface. Changing this forces a new Dedicated Hardware Security Module to be created.

* `subnet_id` - (Required) The ID of the subnet. Changing this forces a new Dedicated Hardware Security Module to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Dedicated Hardware Security Module.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dedicated Hardware Security Module.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dedicated Hardware Security Module.
* `update` - (Defaults to 30 minutes) Used when updating the Dedicated Hardware Security Module.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dedicated Hardware Security Module.

## Import

Dedicated Hardware Security Module can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dedicated_hardware_security_module.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HardwareSecurityModules/dedicatedHSMs/hsm1
```
