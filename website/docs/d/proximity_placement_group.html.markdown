---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_proximity_placement_group"
sidebar_current: "docs-azurerm-datasource-proximity-placement-group"
description: |-
  Gets information about an existing Proximity Placement Group.
---

# Data Source: azurerm_proximity_placement_group

Use this data source to access information about an existing Proximity Placement Group.

## Example Usage

```hcl
data "azurerm_proximity_placement_group" "example" {
  name                = "tf-appsecuritygroup"
  resource_group_name = "my-resource-group"
}

output "proximity_placement_group_id" {
  value = "${data.azurerm_proximity_placement_group.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Proximity Placement Group.

* `resource_group_name` - The name of the resource group in which the Proximity Placement Group exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Proximity Placement Group.
