---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_fleet_member"
description: |-
  Manages a Kubernetes Fleet Member.
---

# azurerm_kubernetes_fleet_member

Manages a Kubernetes Fleet Member.

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
  hub_profile {
    dns_prefix = "example-dns-prefix"
  }
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "example"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_fleet_member" "example" {
  name                        = "example"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.example.id
  kubernetes_cluster_id       = azurerm_kubernetes_cluster.example.id
  group                       = "example-group"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Kubernetes Fleet Member. Changing this forces a new Kubernetes Fleet Member to be created.

* `kubernetes_fleet_manager_id` - (Required) The ID of the Fleet Manager. Changing this forces a new Kubernetes Fleet Member to be created.

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster. Changing this forces a new Kubernetes Fleet Member to be created.

* `group` - (Required) The group that this member belongs to for multi-cluster update management.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Fleet Member.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Fleet Member.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Fleet Member.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Fleet Member.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Fleet Member.

## Import

Kubernetes Fleet Member can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_fleet_member.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/fleets/fleet1/members/member1
```
