---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_route_table"
description: |-
  Gets information about an existing Virtual Hub Route Table
---

# Data Source: azurerm_virtual_hub_route_table

Uses this data source to access information about an existing Virtual Hub Route Table.

## Virtual Hub Route Table Usage

```hcl
data "azurerm_virtual_hub_route_table" "example" {
  name                = "example-hub-route-table"
  resource_group_name = "example-resources"
  virtual_hub_name    = "example-hub-name"
}

output "virtual_hub_route_table_id" {
  value = data.azurerm_virtual_hub_route_table.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Virtual Hub Route Table.

* `resource_group_name` - The Name of the Resource Group where the Virtual Hub Route Table exists.
*  `virtual_hub_name` - The name which should be used for Virtual Hub Route Table.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub Route Table.

* `virtual_hub_id` - The ID of the Virtual Hub within which this route table is created

* `labels` - List of labels associated with this route table.

* `route` - A `route` block as defined below.

---

An `route` block exports the following:

* `name` - The name which is used for this route.

* `destinations` - A list of destination addresses for this route.

* `destinations_type` - The type of destinations.

* `next_hop` - The next hop's resource ID.

* `next_hop_type` - The type of next hop.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub.
