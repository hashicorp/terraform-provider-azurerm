---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_fleet_manager"
description: |-
  Manages a Kubernetes Fleet Manager.
---

# azurerm_kubernetes_fleet_manager

Manages a Kubernetes Fleet Manager.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_kubernetes_fleet_manager" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Kubernetes Fleet Manager should exist. Changing this forces a new Kubernetes Fleet Manager to be created.

* `name` - (Required) Specifies the name of this Kubernetes Fleet Manager. Changing this forces a new Kubernetes Fleet Manager to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Kubernetes Fleet Manager should exist. Changing this forces a new Kubernetes Fleet Manager to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Kubernetes Fleet Manager.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Fleet Manager.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Fleet Manager.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Fleet Manager.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Fleet Manager.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Fleet Manager.

## Import

An existing Kubernetes Fleet Manager can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_fleet_manager.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/fleets/{fleetName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Kubernetes Fleet Manager exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Kubernetes Fleet Manager exists. For example `example-resource-group`.
* Where `{fleetName}` is the name of the Fleet. For example `fleetValue`.
