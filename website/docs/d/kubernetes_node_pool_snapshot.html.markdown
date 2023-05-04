---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_node_pool_snapshot"
description: |-
  Gets information about an existing Kubernetes Snapshot
---

# Data Source: azurerm_kubernetes_node_pool_snapshot

Use this data source to access information about an existing Kubernetes Snapshot.

## Example Usage

```hcl
data "azurerm_kubernetes_node_pool_snapshot" "example" {
  name                = "example"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Kubernetes Snapshot.

* `resource_group_name` - The name of the Resource Group in which the Kubernetes Snapshot exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kubernetes Snapshot.

* `source_resource_id` - The ID of the source resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Snapshot.
