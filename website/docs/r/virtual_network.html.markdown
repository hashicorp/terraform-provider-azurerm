---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network"
description: |-
  Manages a virtual network including any configured subnets. Each subnet can optionally be configured with a security group to be associated with the subnet.
---

# azurerm_virtual_network

Manages a virtual network including any configured subnets. Each subnet can
optionally be configured with a security group to be associated with the subnet.

~> **NOTE on Virtual Networks and Subnet's:** Terraform currently
provides both a standalone [Subnet resource](subnet.html), and allows for Subnets to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line Subnets in conjunction with any Subnet resources. Doing so will cause a conflict of Subnet configurations and will overwrite Subnet's.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_group" "example" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_ddos_protection_plan" "example" {
  name                = "ddospplan1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "virtualNetwork1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]

  ddos_protection_plan {
    id     = azurerm_network_ddos_protection_plan.example.id
    enable = true
  }

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }

  subnet {
    name           = "subnet2"
    address_prefix = "10.0.2.0/24"
  }

  subnet {
    name           = "subnet3"
    address_prefix = "10.0.3.0/24"
    security_group = azurerm_network_security_group.example.id
  }

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual network. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the virtual network.

* `address_space` - (Required) The address space that is used the virtual network. You can supply more than one address space.

* `location` - (Required) The location/region where the virtual network is created. Changing this forces a new resource to be created.

* `bgp_community` - (Optional) The BGP community attribute in format `<as-number>:<community-value>`.

-> **NOTE** The `as-number` segment is the Microsoft ASN, which is always `12076` for now.

* `ddos_protection_plan` - (Optional) A `ddos_protection_plan` block as documented below.

* `dns_servers` - (Optional) List of IP addresses of DNS servers

* `subnet` - (Optional) Can be specified multiple times to define multiple subnets. Each `subnet` block supports fields documented below.

-> **NOTE** Since `subnet` can be configured both inline and via the separate `azurerm_subnet` resource, we have to explicitly set it to empty slice (`[]`) to remove it.

* `vm_protection_enabled` - (Optional) Whether to enable VM protection for all the subnets in this Virtual Network. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `ddos_protection_plan` block supports the following:

* `id` - (Required) The ID of DDoS Protection Plan.

* `enable` - (Required) Enable/disable DDoS Protection Plan on Virtual Network.

---

The `subnet` block supports:

* `name` - (Required) The name of the subnet.

* `address_prefix` - (Required) The address prefix to use for the subnet.

* `security_group` - (Optional) The Network Security Group to associate with the subnet. (Referenced by `id`, ie. `azurerm_network_security_group.example.id`)

## Attributes Reference

The following attributes are exported:

* `id` - The virtual NetworkConfiguration ID.

* `name` - The name of the virtual network.

* `resource_group_name` - The name of the resource group in which to create the virtual network.

* `location` - The location/region where the virtual network is created.

* `address_space` - The list of address spaces used by the virtual network.

* `guid` - The GUID of the virtual network.

* `subnet`- One or more `subnet` blocks as defined below.

---

The `subnet` block exports:

* `id` - The ID of this subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Network.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Network.

## Import

Virtual Networks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_network.exampleNetwork /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1
```
