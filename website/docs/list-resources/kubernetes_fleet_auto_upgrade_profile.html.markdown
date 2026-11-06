---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_fleet_auto_upgrade_profile"
description: |-
  Lists Kubernetes Fleet Auto Upgrade Profile resources.
---

# List resource: azurerm_kubernetes_fleet_auto_upgrade_profile

Lists Kubernetes Fleet Auto Upgrade Profile resources.

## Example Usage

### List Auto Upgrade Profiles in a Kubernetes Fleet Manager

```hcl
list "azurerm_kubernetes_fleet_auto_upgrade_profile" "example" {
  provider = azurerm
  config {
    kubernetes_fleet_manager_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ContainerService/fleets/example-fleet"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `kubernetes_fleet_manager_id` - (Required) The ID of the Kubernetes Fleet Manager to query.
