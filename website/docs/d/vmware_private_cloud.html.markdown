---
subcategory: "VMware (AVS)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_vmware_private_cloud"
description: |-
  Gets information about an existing Vmware Private Cloud.
---

# Data Source: azurerm_vmware_private_cloud

Use this data source to access information about an existing Vmware Private Cloud.

## Example Usage

```hcl
data "azurerm_vmware_private_cloud" "example" {
  name                = "existing-vmware-private-cloud"
  resource_group_name = "existing-resgroup"
}

output "id" {
  value = data.azurerm_vmware_private_cloud.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Vmware Private Cloud.

* `resource_group_name` - (Required) The name of the Resource Group where the Vmware Private Cloud exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Vmware Private Cloud.

* `location` - The Azure Region where the Vmware Private Cloud exists.

* `circuit` - A `circuit` block as defined below.

* `internet_connection_enabled` - Is the Vmware Private Cluster connected to the internet?

* `management_cluster` - A `management_cluster` block as defined below.

* `network_subnet_cidr` - The subnet cidr of the Vmware Private Cloud.

* `hcx_cloud_manager_endpoint` - The endpoint for the HCX Cloud Manager.

* `nsxt_manager_endpoint` - The endpoint for the NSX-T Data Center manager.

* `vcsa_endpoint` - The endpoint for Virtual Center Server Appliance.

* `sku_name` - The Name of the SKU used for this Private Cloud.

* `nsxt_certificate_thumbprint` - The thumbprint of the NSX-T Manager SSL certificate.

* `vcenter_certificate_thumbprint` - The thumbprint of the vCenter Server SSL certificate.

* `management_subnet_cidr` - The network used to access vCenter Server and NSX-T Manager.

* `provisioning_subnet_cidr` - The network which isused for virtual machine cold migration, cloning, and snapshot migration.

* `vmotion_subnet_cidr` - The network which is used for live migration of virtual machines.

* `tags` - A mapping of tags assigned to the Vmware Private Cloud.

---

A `circuit` block exports the following:

* `express_route_id` - The ID of the ExpressRoute Circuit.

* `express_route_private_peering_id` - The ID of the ExpressRoute Circuit private peering.

* `primary_subnet_cidr` - The CIDR of the primary subnet.

* `secondary_subnet_cidr` - The CIDR of the secondary subnet.

---

A `management_cluster` block exports the following:

* `id` - The ID of the management cluster.

* `size` - The size of the management cluster.

* `hosts` - The list of the hosts in the management cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Vmware Private Cloud.
