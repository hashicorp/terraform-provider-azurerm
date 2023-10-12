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

~> **NOTE on Virtual Networks and Subnets:** Terraform currently provides both a standalone [Subnet resource](subnet.html), and allows for Subnets to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line Subnets in conjunction with any Subnet resources. Doing so will cause a conflict of Subnet configurations and will overwrite subnets.

~> **NOTE on Virtual Networks and DNS Servers:** Terraform currently provides both a standalone [virtual network DNS Servers resource](virtual_network_dns_servers.html), and allows for DNS servers to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line DNS servers in conjunction with any Virtual Network DNS Servers resources. Doing so will cause a conflict of Virtual Network DNS Servers configurations and will overwrite virtual networks DNS servers.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_group" "example" {
  name                = "example-security-group"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }

  subnet {
    name           = "subnet2"
    address_prefix = "10.0.2.0/24"
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

* `resource_group_name` - (Required) The name of the resource group in which to create the virtual network. Changing this forces a new resource to be created.

* `address_space` - (Required) The address space that is used the virtual network. You can supply more than one address space.

* `location` - (Required) The location/region where the virtual network is created. Changing this forces a new resource to be created. 

---

* `bgp_community` - (Optional) The BGP community attribute in format `<as-number>:<community-value>`.

-> **NOTE** The `as-number` segment is the Microsoft ASN, which is always `12076` for now.

* `ddos_protection_plan` - (Optional) A `ddos_protection_plan` block as documented below.

* `encryption` - (Optional) A `encryption` block as defined below.

* `dns_servers` - (Optional) List of IP addresses of DNS servers

-> **NOTE** Since `dns_servers` can be configured both inline and via the separate `azurerm_virtual_network_dns_servers` resource, we have to explicitly set it to empty slice (`[]`) to remove it.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Virtual Network should exist. Changing this forces a new Virtual Network to be created.

* `flow_timeout_in_minutes` - (Optional) The flow timeout in minutes for the Virtual Network, which is used to enable connection tracking for intra-VM flows. Possible values are between `4` and `30` minutes.

* `subnet` - (Optional) Can be specified multiple times to define multiple subnets. Each `subnet` block supports fields documented below.

-> **NOTE** Since `subnet` can be configured both inline and via the separate `azurerm_subnet` resource, we have to explicitly set it to empty slice (`[]`) to remove it.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `ddos_protection_plan` block supports the following:

* `id` - (Required) The ID of DDoS Protection Plan.

* `enable` - (Required) Enable/disable DDoS Protection Plan on Virtual Network.

---

A `encryption` block supports the following:

* `enforcement` - (Required) Specifies if the encrypted Virtual Network allows VM that does not support encryption. Possible values are `DropUnencrypted` and `AllowUnencrypted`.

---

The `subnet` block supports:

* `name` - (Required) The name of the subnet.

* `address_prefix` - (Required) The address prefix to use for the subnet.

* `security_group` - (Optional) The Network Security Group to associate with the subnet. (Referenced by `id`, ie. `azurerm_network_security_group.example.id`)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The virtual NetworkConfiguration ID.

* `name` - (Required) The name of the virtual network. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the virtual network.

* `location` - (Required) The location/region where the virtual network is created. Changing this forces a new resource to be created.

* `address_space` - (Required) The list of address spaces used by the virtual network.

* `guid` - The GUID of the virtual network.

* `subnet` - (Optional) One or more `subnet` blocks as defined below.

---

The `subnet` block exports:

* `id` - The ID of this subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Network.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Network.

## Import

Virtual Networks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_network.exampleNetwork /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1
```
