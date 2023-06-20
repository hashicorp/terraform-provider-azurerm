---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_connection"
description: |-
  Gets information about an existing Virtual Hub Connection
---

# Data Source: azurerm_virtual_hub_connection

Uses this data source to access information about an existing Virtual Hub Connection.

## Virtual Hub Connection Usage

```hcl
data "azurerm_virtual_hub_connection" "example" {
  name                = "example-connection"
  resource_group_name = "example-resources"
  virtual_hub_name    = "example-hub-name"
}

output "virtual_hub_connection_id" {
  value = data.azurerm_virtual_hub_connection.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Connection which should be retrieved.

* `resource_group_name` - The Name of the Resource Group where the Virtual Hub Connection exists.

*  `virtual_hub_name` - The name of the Virtual Hub where this Connection exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub Connection.

* `virtual_hub_id` - The ID of the Virtual Hub within which this connection is created

* `remote_virtual_network_id` - The ID of the Virtual Network which the Virtual Hub is connected

* `internet_security_enabled` - Whether Internet Security is enabled to secure internet traffic on this connection

* `routing` - A `routing` block as defined below.

---

An `routing` block exports the following:

* `associated_route_table_id` - The ID of the route table associated with this Virtual Hub connection.

* `propagated_route_table` - A `propagated_route_table` block as defined below.

* `static_vnet_route` - A `static_vnet_route` block as defined below.

---

A `propagated_route_table` block supports the following:

* `labels` - The list of labels assigned to this route table.

* `route_table_ids` - A list of Route Table IDs associated with this Virtual Hub Connection.

---

A `static_vnet_route` block supports the following:

* `name` - The name which is used for this Static Route.

* `address_prefixes` - A list of CIDR Ranges which is used as Address Prefixes.

* `next_hop_ip_address` - The IP Address which is used for the Next Hop.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub.
