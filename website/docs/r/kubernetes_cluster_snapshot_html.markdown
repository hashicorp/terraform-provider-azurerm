---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_snapshot"
description: |-
Manages a Kubernetes Cluster Snapshot.

---

# azurerm_kubernetes_cluster_snapshot

Manages a Kubernetes Cluster Snapshot

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
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

resource "azurerm_kubernetes_cluster_snapshot" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_id          = azurerm_kubernetes_cluster_snapshot.example.id
  tags = {
    environment = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Snapshots resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Snapshots. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `cluster_id` - (Required) This is the ARM ID of the source object to be used to create the target object. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Kubernetes Cluster Snapshot ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Snapshots.
* `update` - (Defaults to 30 minutes) Used when updating the Snapshots.
* `read` - (Defaults to 5 minutes) Used when retrieving the Snapshots.
* `delete` - (Defaults to 30 minutes) Used when deleting the Snapshots.

## Import

Kubernetes Cluster Snapshot can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_snapshot.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedclustersnapshots/managedclustersnapshot1
```
