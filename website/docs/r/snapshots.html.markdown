---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_snapshots"
description: |-
  Manages a kubernetes Snapshots.

---

# azurerm_snapshots

Manages a kubernetes Snapshots.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "snapshots-rg"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "acctestaks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "acctestaks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "example" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}

resource "azurerm_snapshots" "example" {
  name                = "acctestss"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  source_resource_id  = azurerm_kubernetes_cluster_node_pool.example.id
  snapshot_type       = "NodePool"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Snapshots resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Snapshots. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `source_resource_id` - (Required) This is the ARM ID of the source object to be used to create the target object. Changing this forces a new resource to be created.

* `snapshot_type` - (Optional) Possible values include: 'SnapshotTypeNodePool'.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Snapshots ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Snapshots.
* `update` - (Defaults to 30 minutes) Used when updating the Snapshots.
* `read` - (Defaults to 5 minutes) Used when retrieving the Snapshots.
* `delete` - (Defaults to 30 minutes) Used when deleting the Snapshots.

## Import

Snapshots can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_snapshots.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/snapshots/snapshot1
```
