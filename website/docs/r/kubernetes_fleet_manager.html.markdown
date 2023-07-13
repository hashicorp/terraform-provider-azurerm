---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_fleet_manager"
description: |-
  Manages a Kubernetes Fleet Manager.
---

<!-- Note: This documentation is generated. Any manual changes will be overwritten -->

# azurerm_kubernetes_fleet_manager

Manages a Kubernetes Fleet Manager

~> **Note:** This Resource is in **Preview** to use this you must be opted into the Preview. You can do this by running `az feature register --namespace Microsoft.ContainerService --name FleetResourcePreview` and then `az provider register -n Microsoft.ContainerService`
.

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

* `hub_profile` - (Optional) A `hub_profile` block as defined below. The FleetHubProfile configures the Fleet's hub. Changing this forces a new Kubernetes Fleet Manager to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Kubernetes Fleet Manager.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Fleet Manager.

---

## Blocks Reference

### `hub_profile` Block


The `hub_profile` block supports the following arguments:

* `dns_prefix` - (Required) 


In addition to the arguments defined above, the `hub_profile` block exports the following attributes:

* `fqdn` - 
* `kubernetes_version` -

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Kubernetes Fleet Manager.
* `delete` - (Defaults to 30 minutes) Used when deleting this Kubernetes Fleet Manager.
* `read` - (Defaults to 5 minutes) Used when retrieving this Kubernetes Fleet Manager.
* `update` - (Defaults to 30 minutes) Used when updating this Kubernetes Fleet Manager.

## Import

An existing Kubernetes Fleet Manager can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_fleet_manager.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/fleets/{fleetName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Kubernetes Fleet Manager exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Kubernetes Fleet Manager exists. For example `example-resource-group`.
* Where `{fleetName}` is the name of the Fleet. For example `fleetValue`.
