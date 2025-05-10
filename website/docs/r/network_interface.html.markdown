---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_interface"
description: |-
  Manages a Network Interface.

---

# azurerm_network_interface

Manages a Network Interface.

## Example Usage

```hcl
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

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip_configuration` - (Required) One or more `ip_configuration` blocks as defined below.

* `location` - (Required) The location where the Network Interface should exist. Changing this forces a new resource to be created.

* `name` - (Required) The name of the Network Interface. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the Resource Group in which to create the Network Interface. Changing this forces a new resource to be created.

---

* `auxiliary_mode` - (Optional) Specifies the auxiliary mode used to enable network high-performance feature on Network Virtual Appliances (NVAs). This feature offers competitive performance in Connections Per Second (CPS) optimization, along with improvements to handling large amounts of simultaneous connections. Possible values are `AcceleratedConnections`, `Floating`, `MaxConnections` and `None`.

-> **Note:** `auxiliary_mode` is in **Preview** and requires that the preview is enabled - [more information can be found in the Azure documentation](https://learn.microsoft.com/azure/networking/nva-accelerated-connections#prerequisites).

* `auxiliary_sku` - (Optional) Specifies the SKU used for the network high-performance feature on Network Virtual Appliances (NVAs). Possible values are `A8`, `A4`, `A1`, `A2` and `None`.

-> **Note:** `auxiliary_sku` is in **Preview** and requires that the preview is enabled - [more information can be found in the Azure documentation](https://learn.microsoft.com/azure/networking/nva-accelerated-connections#prerequisites).

* `dns_servers` - (Optional) A list of IP Addresses defining the DNS Servers which should be used for this Network Interface.

-> **Note:** Configuring DNS Servers on the Network Interface will override the DNS Servers defined on the Virtual Network.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Network Interface should exist. Changing this forces a new Network Interface to be created.

* `ip_forwarding_enabled` - (Optional) Should IP Forwarding be enabled? Defaults to `false`.

* `accelerated_networking_enabled` - (Optional) Should Accelerated Networking be enabled? Defaults to `false`.

-> **Note:** Only certain Virtual Machine sizes are supported for Accelerated Networking - [more information can be found in this document](https://docs.microsoft.com/azure/virtual-network/create-vm-accelerated-networking-cli).

-> **Note:** To use Accelerated Networking in an Availability Set, the Availability Set must be deployed onto an Accelerated Networking enabled cluster.

* `internal_dns_name_label` - (Optional) The (relative) DNS Name used for internal communications between Virtual Machines in the same Virtual Network.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `ip_configuration` block supports the following:

* `name` - (Required) A name used for this IP Configuration.

* `gateway_load_balancer_frontend_ip_configuration_id` - (Optional) The Frontend IP Configuration ID of a Gateway SKU Load Balancer.

* `subnet_id` - (Optional) The ID of the Subnet where this Network Interface should be located in.

-> **Note:** This is required when `private_ip_address_version` is set to `IPv4`.

* `private_ip_address_version` - (Optional) The IP Version to use. Possible values are `IPv4` or `IPv6`. Defaults to `IPv4`.

* `private_ip_address_allocation` - (Required) The allocation method used for the Private IP Address. Possible values are `Dynamic` and `Static`.

~> **Note:** `Dynamic` means "An IP is automatically assigned during creation of this Network Interface"; `Static` means "User supplied IP address will be used"

* `public_ip_address_id` - (Optional) Reference to a Public IP Address to associate with this NIC

* `primary` - (Optional) Is this the Primary IP Configuration? Must be `true` for the first `ip_configuration` when multiple are specified. Defaults to `false`.

When `private_ip_address_allocation` is set to `Static` the following fields can be configured:

* `private_ip_address` - (Optional) The Static IP Address which should be used.

When `private_ip_address_version` is set to `IPv4` the following fields can be configured:

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `applied_dns_servers` - If the Virtual Machine using this Network Interface is part of an Availability Set, then this list will have the union of all DNS servers from all Network Interfaces that are part of the Availability Set.

* `id` - The ID of the Network Interface.

* `internal_domain_name_suffix` - Even if `internal_dns_name_label` is not specified, a DNS entry is created for the primary NIC of the VM. This DNS name can be constructed by concatenating the VM name with the value of `internal_domain_name_suffix`.

* `mac_address` - The Media Access Control (MAC) Address of the Network Interface.

* `private_ip_address` - The first private IP address of the network interface.

~> **Note:** If a `Dynamic` allocation method is used Azure will allocate an IP Address on Network Interface creation.

* `private_ip_addresses` - The private IP addresses of the network interface.

~> **Note:** If a `Dynamic` allocation method is used Azure will allocate an IP Address on Network Interface creation.

* `virtual_machine_id` - The ID of the Virtual Machine which this Network Interface is connected to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Interface.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Interface.
* `update` - (Defaults to 30 minutes) Used when updating the Network Interface.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Interface.

## Import

Network Interfaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_interface.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkInterfaces/nic1
```
