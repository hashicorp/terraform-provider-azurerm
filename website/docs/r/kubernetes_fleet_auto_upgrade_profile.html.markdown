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
  name     = "example-rg"
  location = "westeurope"
}

resource "azurerm_kubernetes_fleet_manager" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_kubernetes_fleet_update_strategy" "example" {
  name                        = "example"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.example.id
  stage {
    name = "example-stage"
    group {
      name = "example-group"
    }
  }
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "example" {
  name                        = "example"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.example.id
  channel                     = "Rapid"
  node_image_selection_type   = "Latest"
  update_strategy_id          = azurerm_kubernetes_fleet_update_strategy.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Kubernetes Fleet Auto Upgrade Profile. Changing this forces a new Kubernetes Fleet Auto Upgrade Profile to be created.

* `kubernetes_fleet_manager_id` - (Required) The ID of the Fleet Manager. Changing this forces a new Kubernetes Fleet Auto Upgrade Profile to be created.

* `channel` - (Required) The upgrade channel for the auto upgrade profile. Possible values are `Stable`, `Rapid`, and `NodeImage`.

* `node_image_selection_type` - (Optional) The node image selection type for the auto upgrade profile. Possible values are `Consistent` and `Latest`.

* `update_strategy_id` - (Optional) The ID of the Fleet Update Strategy to use for this auto upgrade profile. Changing this forces a new Kubernetes Fleet Auto Upgrade Profile to be created.

* `disabled` - (Optional) Whether the auto upgrade profile is disabled.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Fleet Auto Upgrade Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Fleet Auto Upgrade Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Fleet Auto Upgrade Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Fleet Auto Upgrade Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Fleet Auto Upgrade Profile.

## Import

Kubernetes Fleet Auto Upgrade Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_fleet_auto_upgrade_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/fleets/fleet1/autoUpgradeProfiles/profile1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ContainerService` - 2025-03-01
