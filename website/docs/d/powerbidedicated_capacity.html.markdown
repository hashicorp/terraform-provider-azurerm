---
subcategory: "PowerBIDedicated"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_powerbidedicated_capacity"
sidebar_current: "docs-azurerm-datasource-powerbidedicated-capacity"
description: |-
  Gets information about an existing PowerBIDedicated Capacity
---

# Data Source: azurerm_powerbidedicated_capacity

Use this data source to access information about an existing PowerBIDedicated Capacity.

## Example Usage

```hcl
data "azurerm_powerbidedicated_capacity" "example" {
  resource_group_name = "acctestRG"
  name                = "example-capacity"
}

output "powerbidedicated_capacity_id" {
  value = "${data.azurerm_powerbidedicated_capacity.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PowerBIDedicated Capacity.

* `resource_group_name` - (Required) The Name of the Resource Group where the PowerBIDedicated Capacity exists.

## Attributes Reference

The following attributes are exported:

* `location` - The supported Azure location where the resource exists.

* `sku` - The SKU of the PowerBIDedicated Capacity.

* `administrators` - A set of administrator user identities.

* `tags` - A mapping of tags to assign to the resource.
