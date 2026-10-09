---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_routing_appliance"
description: |-
  Manages a Virtual Network Routing Appliance.
---

# azurerm_virtual_network_routing_appliance

Manages a Virtual Network Routing Appliance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-virtual-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "VirtualNetworkApplianceSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/24"]
}

resource "azurerm_virtual_network_routing_appliance" "example" {
  name                = "example-virtual-network-routing-appliance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  bandwidth_in_gbps   = 10
  subnet_id           = azurerm_subnet.example.id

  tags = {
    environment = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Network Routing Appliance. Must be between 1 and 64 characters, start with a letter or number, end with a letter, number, or underscore, and contain only letters, numbers, underscores, periods, and hyphens. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Network Routing Appliance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Virtual Network Routing Appliance should exist. Changing this forces a new resource to be created.

* `bandwidth_in_gbps` - (Required) The bandwidth of the Virtual Network Routing Appliance in Gbps. Possible values are `10`, `50`, `100`, and `200`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the dedicated subnet named `VirtualNetworkApplianceSubnet` where the Virtual Network Routing Appliance should be deployed. Changing this forces a new resource to be created.

* `private_ip_address_version` - (Optional) The IP address version of the Virtual Network Routing Appliance. Possible values are `IPv4` and `DualStack`. Defaults to `IPv4`. Changing this forces a new resource to be created.

~> **Note:** When `private_ip_address_version` is `DualStack`, the Virtual Network and subnet must have both IPv4 and IPv6 address prefixes.

* `tags` - (Optional) A mapping of tags which should be assigned to the Virtual Network Routing Appliance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Network Routing Appliance.

* `ip_configuration` - A list of `ip_configuration` blocks as defined below.

* `resource_guid` - The resource GUID of the Virtual Network Routing Appliance.

---

An `ip_configuration` block exports the following:

* `name` - The name of the IP configuration.

* `primary` - Whether this is the primary IP configuration.

* `private_ip_address` - The private IP address assigned to the IP configuration.

* `private_ip_address_allocation` - The private IP address allocation method.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Network Routing Appliance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network Routing Appliance.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Network Routing Appliance.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Network Routing Appliance.

## Import

A Virtual Network Routing Appliance can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_network_routing_appliance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/virtualNetworkAppliances/virtualNetworkAppliance1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-07-01
