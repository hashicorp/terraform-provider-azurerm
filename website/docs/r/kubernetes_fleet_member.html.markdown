---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_fleet_member"
description: |-
  Manages a Kubernetes Fleet Member.
---

<!-- Note: This documentation is generated. Any manual changes will be overwritten -->

# azurerm_kubernetes_fleet_member

Manages a Kubernetes Fleet Member.

## Example Usage

```hcl
resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "acctestaksexample"
  default_node_pool {
    name       = "example-value"
    node_count = "example-value"
    vm_size    = "example-value"
    upgrade_settings {
      max_surge = "example-value"
    }
  }
  identity {
    type = "example-value"
  }
}
resource "azurerm_kubernetes_fleet_manager" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_kubernetes_fleet_member" "example" {
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  kubernetes_fleet_id   = azurerm_kubernetes_fleet_manager.example.id
  name                  = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `kubernetes_cluster_id` - (Required) The ARM resource ID of the cluster that joins the Fleet. Changing this forces a new Kubernetes Fleet Member to be created.

* `kubernetes_fleet_id` - (Required) Specifies the Kubernetes Fleet Id within which this Kubernetes Fleet Member should exist. Changing this forces a new Kubernetes Fleet Member to be created.

* `name` - (Required) Specifies the name of this Kubernetes Fleet Member. Changing this forces a new Kubernetes Fleet Member to be created.

* `group` - (Optional) The group this member belongs to for multi-cluster update management.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Fleet Member.

---



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Kubernetes Fleet Member.
* `delete` - (Defaults to 30 minutes) Used when deleting this Kubernetes Fleet Member.
* `read` - (Defaults to 5 minutes) Used when retrieving this Kubernetes Fleet Member.
* `update` - (Defaults to 30 minutes) Used when updating this Kubernetes Fleet Member.

## Import

An existing Kubernetes Fleet Member can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_fleet_member.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/fleets/{fleetName}/members/{memberName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Kubernetes Fleet Member exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Kubernetes Fleet Member exists. For example `example-resource-group`.
* Where `{fleetName}` is the name of the Fleet. For example `fleetValue`.
* Where `{memberName}` is the name of the Member. For example `memberValue`.
