---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_interface"
description: |-
  Manages a Network Interface located in a Virtual Network, usually attached to a Virtual Machine.

---

# azurerm_network_interface

Manages a Network Interface located in a Virtual Network, usually attached to a Virtual Machine.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_virtual_network" "example" {
  name                = "acceptanceTestVirtualNetwork1"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "example" {
  name                = "acceptanceTestNetworkInterface1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    environment = "staging"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network interface. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the network interface. Changing this forces a new resource to be created.

* `location` - (Required) The location/region where the network interface is created. Changing this forces a new resource to be created.

* `network_security_group_id` - (Optional) The ID of the Network Security Group to associate with the network interface.

* `internal_dns_name_label` - (Optional) Relative DNS name for this NIC used for internal communications between VMs in the same VNet

* `enable_ip_forwarding` - (Optional) Enables IP Forwarding on the NIC. Defaults to `false`.

* `enable_accelerated_networking` - (Optional) Enables Azure Accelerated Networking using SR-IOV. Only certain VM instance sizes are supported. Refer to [Create a Virtual Machine with Accelerated Networking](https://docs.microsoft.com/en-us/azure/virtual-network/create-vm-accelerated-networking-cli). Defaults to `false`.

~> **NOTE:** when using Accelerated Networking in an Availability Set - the Availability Set must be deployed on an Accelerated Networking enabled cluster.

* `dns_servers` - (Optional) List of DNS servers IP addresses to use for this NIC, overrides the VNet-level server list

* `ip_configuration` - (Required) One or more `ip_configuration` associated with this NIC as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `ip_configuration` block supports:

* `name` - (Required) User-defined name of the IP.

* `subnet_id` - (Optional) Reference to a subnet in which this NIC has been created. Required when `private_ip_address_version` is IPv4.

* `private_ip_address` - (Optional) Static IP Address.

* `private_ip_address_allocation` - (Required) Defines how a private IP address is assigned. Options are Static or Dynamic.

* `private_ip_address_version` - (Optional) The IP Version to use. Possible values are `IPv4` or `IPv6`. Defaults to `IPv4`.

* `public_ip_address_id` - (Optional) Reference to a Public IP Address to associate with this NIC

* `application_gateway_backend_address_pools_ids` - (Optional / **Deprecated**) List of Application Gateway Backend Address Pool IDs references to which this NIC belongs

-> **NOTE:** At this time Network Interface <-> Application Gateway Backend Address Pool associations need to be configured both using this field (which is now Deprecated) and using the `azurerm_network_interface_application_gateway_backend_address_pool_association` resource. This field is deprecated and will be removed in favour of that resource in the next major version (2.0) of the AzureRM Provider.

* `load_balancer_backend_address_pools_ids` - (Optional / **Deprecated**) List of Load Balancer Backend Address Pool IDs references to which this NIC belongs

-> **NOTE:** At this time Network Interface <-> Load Balancer Backend Address Pool associations need to be configured both using this field (which is now Deprecated) and using the `azurerm_network_interface_backend_address_pool_association` resource. This field is deprecated and will be removed in favour of that resource in the next major version (2.0) of the AzureRM Provider.

* `load_balancer_inbound_nat_rules_ids` - (Optional / **Deprecated**) List of Load Balancer Inbound Nat Rules IDs involving this NIC

-> **NOTE:** At this time Network Interface <-> Load Balancer Inbound NAT Rule associations need to be configured both using this field (which is now Deprecated) and using the `azurerm_network_interface_nat_rule_association` resource. This field is deprecated and will be removed in favour of that resource in the next major version (2.0) of the AzureRM Provider.

* `application_security_group_ids` - (Optional / **Deprecated**) List of Application Security Group IDs which should be attached to this NIC

-> **NOTE:** At this time Network Interface <-> Application Security Group associations need to be configured both using this field (which is now Deprecated) and using the `azurerm_network_interface_application_security_group_association` resource. This field is deprecated and will be removed in favour of that resource in the next major version (2.0) of the AzureRM Provider.

* `primary` - (Optional) Is this the Primary Network Interface? If set to `true` this should be the first `ip_configuration` in the array.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Network Interface.
* `mac_address` - The media access control (MAC) address of the network interface.
* `private_ip_address` - The first private IP address of the network interface.
* `private_ip_addresses` - The private IP addresses of the network interface.
* `virtual_machine_id` - Reference to a VM with which this NIC has been associated.
* `applied_dns_servers` - If the VM that uses this NIC is part of an Availability Set, then this list will have the union of all DNS servers from all NICs that are part of the Availability Set

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Interface.
* `update` - (Defaults to 30 minutes) Used when updating the Network Interface.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Interface.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Interface.

## Import

Network Interfaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_interface.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.network/networkInterfaces/nic1
```
