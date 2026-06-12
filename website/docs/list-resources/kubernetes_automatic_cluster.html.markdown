---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_automatic_cluster"
description: |-
    Lists Kubernetes Automatic Cluster resources.
---

# List resource: azurerm_kubernetes_automatic_cluster

Lists Kubernetes Automatic Cluster resources.

## Example Usage

### List all Kubernetes Automatic Clusters in the subscription

```hcl
list "azurerm_kubernetes_automatic_cluster" "example" {
  provider = azurerm
  config {}
}
```

### List all Kubernetes Automatic Clusters in a specific resource group

```hcl
list "azurerm_kubernetes_automatic_cluster" "example" {
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

