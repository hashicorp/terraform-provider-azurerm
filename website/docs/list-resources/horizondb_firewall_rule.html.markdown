---
subcategory: "HorizonDB"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_horizondb_firewall_rule"
description: |-
  Lists HorizonDB Firewall Rule resources.
---

# List resource: azurerm_horizondb_firewall_rule

Lists HorizonDB Firewall Rule resources.

## Example Usage

```hcl
list "azurerm_horizondb_firewall_rule" "example" {
  provider = azurerm
  config {
    horizondb_pool_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HorizonDb/clusters/cluster1/pools/pool1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `horizondb_pool_id` - (Required) The ID of the HorizonDB Pool.
