---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_avs_cluster"
description: |-
  Manages a avs Cluster.
---

# azurerm_avs_cluster

Manages a avs Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_avs_private_cloud" "example" {
  name = "example-privatecloud"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
  sku {
      name = "example-privatecloud"
  }

  management_cluster {
      cluster_size = 42
  }
  network_block = ""
}

resource "azurerm_avs_cluster" "example" {
  name = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  private_cloud_name = azurerm_avs_private_cloud.example.name
  sku {
      name = "example-cluster"
  }
  cluster_size = 42
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this avs Cluster. Changing this forces a new avs Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the avs Cluster should exist. Changing this forces a new avs Cluster to be created.

* `private_cloud_name` - (Required) The name of the private cloud. Changing this forces a new avs Cluster to be created.

* `sku` - (Required)  A `sku` block as defined below. Changing this forces a new avs Cluster to be created.

* `cluster_size` - (Required) The cluster size.

---

An `sku` block exports the following:

* `name` - (Required) The name which should be used for this sku. Changing this forces a new avs Cluster to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the avs Cluster.

* `cluster_id` - The ID of the cluster.

* `hosts` - The hosts.

* `type` - Resource type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the avs Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the avs Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the avs Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the avs Cluster.

## Import

avs Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_avs_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/privateCloud1/clusters/cluster1
```