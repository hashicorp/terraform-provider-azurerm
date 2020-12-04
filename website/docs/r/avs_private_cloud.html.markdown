---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_avs_private_cloud"
description: |-
  Manages a Azure Vmware Solution Private Cloud.
---

# azurerm_avs_private_cloud

Manages a Azure Vmware Solution Private Cloud.

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

resource "azurerm_avs_private_cloud" "example" {
  name                = "example-avs-private-cloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "AV36"

  management_cluster {
    cluster_size = 3
  }

  network_block      = "192.168.48.0/22"
  internet_connected = false
  nsxt_password      = "QazWsx13$Edc"
  vcenter_password   = "QazWsx13$Edc"
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Vmware Solution Private Cloud. Changing this forces a new Azure Vmware Solution Private Cloud to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Vmware Solution Private Cloud should exist. Changing this forces a new Azure Vmware Solution Private Cloud to be created.

* `location` - (Required) The Azure Region where the Azure Vmware Solution Private Cloud should exist. Changing this forces a new Azure Vmware Solution Private Cloud to be created.

* `management_cluster` - (Required) A `management_cluster` block as defined below.

* `network_block` - (Required) The block of addresses which should be unique across virtual network in your subscription as well as on-premise. Changing this forces a new Azure Vmware Solution Private Cloud to be created.

* `sku_name` - (Required) The name of the SKU. Changing this forces a new Azure Vmware Solution Private Cloud to be created. Possible values are "av20", "av36" and "av36t".

* `internet_connected` - (Optional) Is connected to the internet?

* `nsxt_password` - (Optional) The password of the NSX-T Manager. Changing this forces a new Azure Vmware Solution Private Cloud to be created.

* `vcenter_password` - (Optional) The password of the vCenter admin. Changing this forces a new Azure Vmware Solution Private Cloud to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Vmware Solution Private Cloud.

---

A `management_cluster` block supports the following:

* `cluster_size` - (Required) The size of the cluster.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Vmware Solution Private Cloud.

* `circuit` - A `circuit` block as defined below.

* `hcx_cloud_manager_endpoint` - The endpoint for the HCX Cloud Manager.

* `nsxt_manager_endpoint` - The endpoint for the NSX-T Data Center manager.

* `vcsa_endpoint` - The endpoint for Virtual Center Server Appliance.

* `nsxt_certificate_thumbprint` - The thumbprint of the NSX-T Manager SSL certificate.

* `vcenter_certificate_thumbprint` - The thumbprint of the vCenter Server SSL certificate.

* `management_network` - The network used to access vCenter Server and NSX-T Manager.

* `provisioning_network` - The network which isused for virtual machine cold migration, cloning, and snapshot migration.

* `vmotion_network` - The network which is used for live migration of virtual machines.

---

A `circuit` block exports the following:

* `express_route_id` - The ID of the ExpressRoute Circuit.

* `express_route_private_peering_id` - The ID of the ExpressRoute Circuit private peering.

* `primary_subnet` - The CIDR of primary subnet.

* `secondary_subnet` - The CIDR of secondary subnet.

---

A `management_cluster` block exports the following:

* `cluster_id` - The ID of the cluster.

* `hosts` - The set of the hosts in the cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the Azure Vmware Solution Private Cloud.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Vmware Solution Private Cloud.
* `update` - (Defaults to 6 hours) Used when updating the Azure Vmware Solution Private Cloud.
* `delete` - (Defaults to 6 hours) Used when deleting the Azure Vmware Solution Private Cloud.

## Import

Azure Vmware Solution Private Clouds can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_avs_private_cloud.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/PrivateClouds/privateCloud1
```
