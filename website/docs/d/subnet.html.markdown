---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet"
sidebar_current: "docs-azurerm-datasource-subnet"
description: |-
  Gets information about an existing Subnet located within a Virtual Network.
---

# Data Source: azurerm_subnet

Use this data source to access information about an existing Subnet within a Virtual Network.

## Example Usage

```hcl
data "azurerm_subnet" "example" {
  name                 = "backend"
  virtual_network_name = "production"
  resource_group_name  = "networking"
}

output "subnet_id" {
  value = "${data.azurerm_subnet.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Subnet.
* `virtual_network_name` - (Required) Specifies the name of the Virtual Network this Subnet is located within.
* `resource_group_name` - (Required) Specifies the name of the resource group the Virtual Network is located in.

## Attributes Reference

* `id` - The ID of the Subnet.
* `address_prefix` - The address prefix used for the subnet.
* `enforce_private_link_service_network_policies` - Enable or Disable network policies on private link service in the subnet.
* `network_security_group_id` - The ID of the Network Security Group associated with the subnet.
* `route_table_id` - The ID of the Route Table associated with this subnet.
* `ip_configurations` - The collection of IP Configurations with IPs within this subnet.
* `service_endpoints` - A list of Service Endpoints within this subnet.
* `enforce_private_link_endpoint_network_policies` - Enable or Disable network policies for the private link endpoint on the subnet.
* `enforce_private_link_service_network_policies` - Enable or Disable network policies for the private link service on the subnet.
