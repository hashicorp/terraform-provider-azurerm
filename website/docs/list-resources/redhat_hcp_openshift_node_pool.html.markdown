---
subcategory: "Red Hat OpenShift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_hcp_openshift_node_pool"
description: |-
  Lists Azure Red Hat OpenShift with Hosted Control Planes Node Pool resources.
---

# List resource: azurerm_redhat_hcp_openshift_node_pool

Lists Azure Red Hat OpenShift with Hosted Control Planes Node Pool resources.

## Example Usage

### List all Node Pools in an Azure Red Hat OpenShift with Hosted Control Planes Cluster

```hcl
list "azurerm_redhat_hcp_openshift_node_pool" "example" {
  provider = azurerm
  config {
    hcp_openshift_cluster_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/example-cluster"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `hcp_openshift_cluster_id` - (Required) The ID of the Azure Red Hat OpenShift with Hosted Control Planes Cluster whose Node Pools should be listed.
