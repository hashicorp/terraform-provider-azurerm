---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_avs_private_cloud"
description: |-
  Gets information about an existing Azure Vmware Solution Private Cloud.
---

# Data Source: azurerm_avs_private_cloud

Use this data source to access information about an existing Azure Vmware Solution Private Cloud.

## Example Usage

```hcl
data "azurerm_avs_private_cloud" "example" {
  name                = "existing-avs-private-cloud"
  resource_group_name = "existing-resgroup"
}

output "id" {
  value = data.azurerm_avs_private_cloud.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Azure Vmware Solution Private Cloud.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Vmware Solution Private Cloud exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Vmware Solution Private Cloud.

* `location` - The Azure Region where the Azure Vmware Solution Private Cloud exists.

* `circuit` - A `circuit` block as defined below.

* `internet_connection_enabled` - Is the Private CLuster connected to the internet?

* `management_cluster` - A `management_cluster` block as defined below.

* `network_subnet` - The subnet which should be unique across virtual network in your subscription as well as on-premise.

* `hcx_cloud_manager_endpoint` - The endpoint for the HCX Cloud Manager.

* `nsxt_manager_endpoint` - The endpoint for the NSX-T Data Center manager.

* `vcsa_endpoint` - The endpoint for Virtual Center Server Appliance.

* `sku_name` - The Name of the SKU used for this Private Cloud.

* `nsxt_certificate_thumbprint` - The thumbprint of the NSX-T Manager SSL certificate.

* `vcenter_certificate_thumbprint` - The thumbprint of the vCenter Server SSL certificate.

* `management_network` - The network used to access vCenter Server and NSX-T Manager.

* `provisioning_subnet` - The network which isused for virtual machine cold migration, cloning, and snapshot migration.

* `vmotion_subnet` - The network which is used for live migration of virtual machines.

* `tags` - A mapping of tags assigned to the Azure Vmware Solution Private Cloud.

---

A `circuit` block exports the following:

* `express_route_id` - The ID of the ExpressRoute Circuit.

* `express_route_private_peering_id` - The ID of the ExpressRoute Circuit private peering.

* `primary_subnet` - The CIDR of the primary subnet.

* `secondary_subnet` - The CIDR of the secondary subnet.

---

A `management_cluster` block exports the following:

* `id` - The ID of the cluster.

* `size` - The size of the cluster.

* `hosts` - The set of the hosts in the cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Vmware Solution Private Cloud.
