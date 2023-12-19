---
subcategory: "VMware (AVS)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vmware_private_cloud"
description: |-
  Manages a VMware Private Cloud.
---

# azurerm_vmware_private_cloud

Manages a VMware Private Cloud.

## Example Usage

~> **NOTE :**  Normal `terraform apply` could ignore this note. Please disable correlation request id for continuous operations in one build (like acctest). The continuous operations like `update` or `delete` could not be triggered when it shares the same `correlation-id` with its previous operation.

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this VMware Private Cloud. Changing this forces a new VMware Private Cloud to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the VMware Private Cloud should exist. Changing this forces a new VMware Private Cloud to be created.

* `location` - (Required) The Azure Region where the VMware Private Cloud should exist. Changing this forces a new VMware Private Cloud to be created.

* `management_cluster` - (Required) A `management_cluster` block as defined below.
~> **NOTE :** `internet_connection_enabled` and `management_cluster.0.size` cannot be updated at the same time.

* `network_subnet_cidr` - (Required) The subnet which should be unique across virtual network in your subscription as well as on-premise. Changing this forces a new VMware Private Cloud to be created.

* `sku_name` - (Required) The Name of the SKU used for this Private Cloud. Possible values are `av20`, `av36`, `av36t`, `av36p`, `av36pt`, `av52`, `av52t`, and `av64`. Changing this forces a new VMware Private Cloud to be created.

* `internet_connection_enabled` - (Optional) Is the Private Cluster connected to the internet? This field can not updated with `management_cluster.0.size` together.
~> **NOTE :** `internet_connection_enabled` and `management_cluster.0.size` cannot be updated at the same time.

* `nsxt_password` - (Optional) The password of the NSX-T Manager. Changing this forces a new VMware Private Cloud to be created.

* `vcenter_password` - (Optional) The password of the vCenter admin. Changing this forces a new VMware Private Cloud to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the VMware Private Cloud.

---

A `management_cluster` block supports the following:

* `size` - (Required) The size of the management cluster. This field can not updated with `internet_connection_enabled` together.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the VMware Private Cloud.

* `circuit` - A `circuit` block as defined below.

* `hcx_cloud_manager_endpoint` - The endpoint for the HCX Cloud Manager.

* `nsxt_manager_endpoint` - The endpoint for the NSX-T Data Center manager.

* `vcsa_endpoint` - The endpoint for Virtual Center Server Appliance.

* `nsxt_certificate_thumbprint` - The thumbprint of the NSX-T Manager SSL certificate.

* `vcenter_certificate_thumbprint` - The thumbprint of the vCenter Server SSL certificate.

* `management_subnet_cidr` - The network used to access vCenter Server and NSX-T Manager.

* `provisioning_subnet_cidr` - The network which is used for virtual machine cold migration, cloning, and snapshot migration.

* `vmotion_subnet_cidr` - The network which is used for live migration of virtual machines.

---

A `circuit` block exports the following:

* `express_route_id` - The ID of the ExpressRoute Circuit.

* `express_route_private_peering_id` - The ID of the ExpressRoute Circuit private peering.

* `primary_subnet_cidr` - The CIDR of the primary subnet.

* `secondary_subnet_cidr` - The CIDR of the secondary subnet.

---

A `management_cluster` block exports the following:

* `id` - The ID of the management cluster.

* `hosts` - A list of hosts in the management cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 hours) Used when creating the VMware Private Cloud.
* `read` - (Defaults to 5 minutes) Used when retrieving the VMware Private Cloud.
* `update` - (Defaults to 10 hours) Used when updating the VMware Private Cloud.
* `delete` - (Defaults to 10 hours) Used when deleting the VMware Private Cloud.

## Import

VMware Private Clouds can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vmware_private_cloud.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/privateCloud1
```
