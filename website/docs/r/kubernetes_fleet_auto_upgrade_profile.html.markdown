---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_fleet_auto_upgrade_profile"
description: |-
  Manages a Kubernetes Fleet Auto Upgrade Profile.
---

# azurerm_kubernetes_fleet_auto_upgrade_profile

Manages a Kubernetes Fleet Auto Upgrade Profile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_kubernetes_fleet_manager" "example" {
  name                = "example-fleet"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "example" {
  name                = "default"
  resource_group_name = azurerm_resource_group.example.name
  fleet_name          = azurerm_kubernetes_fleet_manager.example.name
  channel             = "Stable"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Auto Upgrade Profile. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Kubernetes Fleet Auto Upgrade Profile should exist.

* `fleet_name` - (Required) The name of the Kubernetes Fleet Manager where the Auto Upgrade Profile should exist.

* `channel` - (Required) The upgrade channel for the Kubernetes cluster. Possible values are `Stable`, `Rapid`, and `NodeImage`.

* `node_image_upgrade_type` - (Optional) The node image upgrade type. Possible values are `Latest` and `Consistent`. Defaults to `Latest`.

* `update_strategy_id` - (Optional) The ID of the update strategy to use for this auto upgrade profile.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Fleet Auto Upgrade Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Fleet Auto Upgrade Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Fleet Auto Upgrade Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Fleet Auto Upgrade Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Fleet Auto Upgrade Profile.

## Import

Kubernetes Fleet Auto Upgrade Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_fleet_auto_upgrade_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/fleets/fleet1/autoUpgradeProfiles/default
``` 

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ContainerService` - 2025-03-01
