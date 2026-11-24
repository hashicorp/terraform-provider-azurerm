---
subcategory: "HorizonDB"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_horizondb_cluster"
description: |-
  Lists HorizonDB Cluster resources.
---

# List resource: azurerm_horizondb_cluster

Lists HorizonDB Cluster resources.

## Example Usage

```hcl
list "azurerm_horizondb_cluster" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-resources"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the Resource Group to query. When omitted, Clusters across the subscription are returned.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the subscription configured on the provider.
