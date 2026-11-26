---
subcategory: "Red Hat OpenShift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_hcp_openshift_external_auth"
description: |-
  Lists Azure Red Hat OpenShift with Hosted Control Planes External Auth resources.
---

# List resource: azurerm_redhat_hcp_openshift_external_auth

Lists external authentication providers in an Azure Red Hat OpenShift with Hosted Control Planes cluster.

## Example Usage

```hcl
list "azurerm_redhat_hcp_openshift_external_auth" "example" {
  provider = azurerm
  config {
    hcp_openshift_cluster_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/example-cluster"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `hcp_openshift_cluster_id` - (Required) The ID of the ARO HCP cluster whose external authentication providers should be listed.
