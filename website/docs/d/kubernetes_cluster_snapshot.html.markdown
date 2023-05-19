---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_snapshot"
description: |-
  Gets information about an existing Kubernetes Cluster Snapshot
---

# Data Source: azurerm_kubernetes_cluster_snapshot

Use this data source to access information about an existing Kubernetes Cluster Snapshot.

## Example Usage

```hcl
data "azurerm_kubernetes_node_pool_snapshot" "example" {
  name                = "example"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Kubernetes Cluster Snapshot.

* `resource_group_name` - The name of the Resource Group in which the Kubernetes Cluster Snapshot exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kubernetes Node Pool Snapshot.

* `source_cluster_id` - The ID of the source Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster Snapshot.
