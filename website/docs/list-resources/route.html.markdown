---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_route"
description: |-
    Lists Route resources.
---

# List resource: azurerm_route

Lists Route resources.

## Example Usage

### List Routes in a Route Table

```hcl
list "azurerm_route" "example" {
  provider = azurerm
  config {
    route_table_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/routeTables/mytable1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `route_table_id` - (Required) The ID of the Route Table to query.
