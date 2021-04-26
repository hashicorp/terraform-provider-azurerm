---
subcategory: "VMware (AVS)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vmware_cluster"
description: |-
  Manages a Vmware Cluster.
---

# azurerm_vmware_cluster

Manages a Vmware Cluster.

## Example Usage

```hcl
provider "azurerm" {
  features {}
  disable_correlation_request_id = true
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_vmware_private_cloud" "example" {
  name                = "example-vmware-private-cloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "av36"

  management_cluster {
    size = 3
  }

  network_subnet_cidr         = "192.168.48.0/22"
  internet_connection_enabled = false
  nsxt_password               = "QazWsx13$Edc"
  vcenter_password            = "WsxEdc23$Rfv"
}

resource "azurerm_vmware_cluster" "example" {
  name               = "example-Cluster"
  vmware_cloud_id    = azurerm_vmware_private_cloud.example.id
  cluster_node_count = 3
  sku_name           = "av36"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Vmware Cluster. Changing this forces a new Vmware Cluster to be created.

* `vmware_cloud_id` - (Required) The ID of the Vmware Private Cloud in which to create this Vmware Cluster. Changing this forces a new Vmware Cluster to be created.

* `cluster_node_count` - (Required) The count of the Vmware Cluster nodes.

* `sku_name` - (Required) The cluster sku to use. Possible values are `av20`, `av36`, and `av36t`. Changing this forces a new Vmware Cluster to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Vmware Cluster.

* `cluster_number` - A number that identifies this Vmware Cluster in its Vmware Private Cloud.

* `hosts` - A list of host of the Vmware Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 hours) Used when creating the Vmware Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Vmware Cluster.
* `update` - (Defaults to 5 hours) Used when updating the Vmware Cluster.
* `delete` - (Defaults to 5 hours) Used when deleting the Vmware Cluster.

## Import

Vmware Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vmware_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AVS/privateClouds/privateCloud1/clusters/cluster1
```
