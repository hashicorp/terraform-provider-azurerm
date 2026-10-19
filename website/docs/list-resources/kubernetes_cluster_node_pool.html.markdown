---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_node_pool"
description: |-
    Lists Kubernetes Cluster Node Pool resources.
---

# List resource: azurerm_kubernetes_cluster_node_pool

Lists Kubernetes Cluster Node Pool resources.

## Example Usage

### List Kubernetes Cluster Node Pools in a Kubernetes Cluster

```hcl
list "azurerm_kubernetes_cluster_node_pool" "example" {
  provider = azurerm
  config {
    kubernetes_cluster_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster to query.
