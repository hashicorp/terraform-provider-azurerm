---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_fleet_update_strategy"
description: |-
  Manages a Kubernetes Fleet Update Strategy.
---

# azurerm_kubernetes_fleet_update_strategy

Manages a Kubernetes Fleet Update Strategy.

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
    name = "example-stage-1"
    group {
      name = "example-group-1"
    }
    after_stage_wait_in_seconds = 21
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Kubernetes Fleet Update Strategy. Changing this forces a new Kubernetes Fleet Update Strategy to be created.

* `kubernetes_fleet_manager_id` - (Required) The ID of the Fleet Manager. Changing this forces a new Kubernetes Fleet Update Strategy to be created.

* `stage` - (Required) One or more `stage` blocks as defined below.

---

A `stage` block supports the following:

* `group` - (Required) One or more `group` blocks as defined below.

* `name` - (Required) The name which should be used for this stage.

* `after_stage_wait_in_seconds` - (Optional) Specifies the time in seconds to wait at the end of this stage before starting the next one.

---

A `group` block supports the following:

* `name` - (Required) The name which should be used for this group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Kubernetes Fleet Update Strategy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Fleet Update Strategy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Fleet Update Strategy.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Fleet Update Strategy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Fleet Update Strategy.

## Import

Kubernetes Fleet Update Strategies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_fleet_update_strategy.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/fleets/fleet1/updateStrategies/updateStrategy1
```
