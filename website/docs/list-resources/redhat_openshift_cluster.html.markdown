---
subcategory: "Red Hat OpenShift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_openshift_cluster"
description: |-
  Lists Red Hat OpenShift Cluster resources.
---

# List resource: azurerm_redhat_openshift_cluster

Lists Red Hat OpenShift Cluster resources.

## Example Usage

### List all Red Hat OpenShift Clusters in the subscription

```hcl
list "azurerm_redhat_openshift_cluster" "example" {
  provider = azurerm
  config {}
}
```

### List all Red Hat OpenShift Clusters in a specific resource group

```hcl
list "azurerm_redhat_openshift_cluster" "example" {
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
