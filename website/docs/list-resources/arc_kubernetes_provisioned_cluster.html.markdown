---
subcategory: "ArcKubernetes"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_kubernetes_provisioned_cluster"
description: |-
    Lists Arc Kubernetes Provisioned Cluster resources.
---

# List resource: azurerm_arc_kubernetes_provisioned_cluster

Lists Arc Kubernetes Provisioned Cluster resources.

## Example Usage

### List all Arc Kubernetes Provisioned Clusters in the subscription

```hcl
list "azurerm_arc_kubernetes_provisioned_cluster" "example" {
  provider = azurerm
  config {}
}
```

### List all Arc Kubernetes Provisioned Clusters in a specific resource group

```hcl
list "azurerm_arc_kubernetes_provisioned_cluster" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
