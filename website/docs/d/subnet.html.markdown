---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet"
description: |-
  Gets information about an existing Subnet located within a Virtual Network.
---

# Data Source: azurerm_subnet

Use this data source to access information about an existing Subnet within a Virtual Network.

## Example Usage

```hcl
data "azurerm_virtual_network" "example" {
  name                = "production"
  resource_group_name = "networking"
}

data "azurerm_subnet" "example" {
  name                 = "backend"
  virtual_network_id = data.azurerm_virtual_network.example.id
}

output "subnet_id" {
  value = data.azurerm_subnet.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Subnet.
* `virtual_network_name` - Specifies the name of the Virtual Network this Subnet is located within.
* `resource_group_name` - Specifies the name of the resource group the Virtual Network is located in.

## Attributes Reference

* `id` - The ID of the Subnet.
* `address_prefix` - (Deprecated) The address prefix used for the subnet.
* `address_prefixes` - The address prefixes for the subnet.
* `enforce_private_link_service_network_policies` - Enable or Disable network policies on private link service in the subnet.
* `network_security_group_id` - The ID of the Network Security Group associated with the subnet.
* `route_table_id` - The ID of the Route Table associated with this subnet.
* `service_endpoints` - A list of Service Endpoints within this subnet.
* `enforce_private_link_endpoint_network_policies` - Enable or Disable network policies for the private link endpoint on the subnet.
* `enforce_private_link_service_network_policies` - Enable or Disable network policies for the private link service on the subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Subnet located within a Virtual Network.
