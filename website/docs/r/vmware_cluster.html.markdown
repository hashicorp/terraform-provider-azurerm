---
subcategory: "Azure VMware Solution"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vmware_cluster"
description: |-
  Manages an Azure VMware Solution Cluster.
---

# azurerm_vmware_cluster

Manages an Azure VMware Solution Cluster.

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

* `name` - (Required) The name which should be used for this Azure VMware Solution Cluster. Changing this forces a new Azure VMware Solution Cluster to be created.

* `vmware_cloud_id` - (Required) The ID of the Azure VMware Solution Private Cloud in which to create this Cluster. Changing this forces a new Azure VMware Solution Cluster to be created.

* `cluster_node_count` - (Required) The count of the Azure VMware Solution Cluster nodes.

* `sku_name` - (Required) The Cluster SKU to use. Possible values are `av20`, `av36`, `av36t`, `av36p`, `av48`, `av48t`, `av36pt`, `av52`, `av52t`, and `av64`. Changing this forces a new Azure VMware Solution Cluster to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure VMware Solution Cluster.

* `cluster_number` - A number that identifies this Cluster in its Azure VMware Solution Private Cloud.

* `hosts` - A list of hosts in the Azure VMware Solution Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 hours) Used when creating the Azure VMware Solution Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure VMware Solution Cluster.
* `update` - (Defaults to 5 hours) Used when updating the Azure VMware Solution Cluster.
* `delete` - (Defaults to 5 hours) Used when deleting the Azure VMware Solution Cluster.

## Import

Azure VMware Solution Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vmware_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AVS/privateClouds/privateCloud1/clusters/cluster1
```
