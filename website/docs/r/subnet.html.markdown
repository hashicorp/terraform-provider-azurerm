---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet"
sidebar_current: "docs-azurerm-resource-network-subnet"
description: |-
  Manages a subnet. Subnets represent network segments within the IP space defined by the virtual network.

---

# azurerm_subnet

Manages a subnet. Subnets represent network segments within the IP space defined by the virtual network.

~> **NOTE on Virtual Networks and Subnet's:** Terraform currently
provides both a standalone [Subnet resource](subnet.html), and allows for Subnets to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line Subnets in conjunction with any Subnet resources. Doing so will cause a conflict of Subnet configurations and will overwrite Subnet's.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_virtual_network" "example" {
  # ...
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.1.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Subnet. Changing this forces a new resource to be created.

-> **NOTE:** This needs to be unique within the Virtual Network.

* `resource_group_name` - (Required) The name of the resource group in which to create the subnet. Changing this forces a new resource to be created.

* `virtual_network_name` - (Required) The name of the virtual network to which to attach the subnet. Changing this forces a new resource to be created.

* `address_prefix` - (Required) The address prefix to use for the subnet.

* `network_security_group_id` - (Optional) The ID of the Network Security Group to associate with the subnet.

* `route_table_id` - (Optional) The ID of the Route Table to associate with the subnet.

* `service_endpoints` - (Optional) The list of Service endpoints to associate with the subnet. Possible values include: `Microsoft.Storage` and `Microsoft.Sql`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Subnet.

* `ip_configurations` - A list of ID's of the IP Configurations attached to this subnet.

## Import

Subnets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subnet.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```
