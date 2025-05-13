---
subcategory: "Azure VMware Solution"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_vmware_private_cloud"
description: |-
  Gets information about an existing Azure VMware Solution Private Cloud.
---

# Data Source: azurerm_vmware_private_cloud

Use this data source to access information about an existing Azure VMware Solution Private Cloud.

## Example Usage

~> **Note:** Normal `terraform apply` could ignore this note. Please disable correlation request id for continuous operations in one build (like acctest). The continuous operations like `update` or `delete` could not be triggered when it shares the same `correlation-id` with its previous operation.

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

* `name` - (Required) The name of this Azure VMware Solution Private Cloud.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure VMware Solution Private Cloud exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure VMware Solution Private Cloud.

* `location` - The Azure Region where the Azure VMware Solution Private Cloud exists.

* `circuit` - A `circuit` block as defined below.

* `internet_connection_enabled` - Is the Azure VMware Solution Private Cloud connected to the internet?

* `management_cluster` - A `management_cluster` block as defined below.

* `network_subnet_cidr` - The subnet CIDR of the Azure VMware Solution Private Cloud.

* `hcx_cloud_manager_endpoint` - The endpoint for the VMware HCX Cloud Manager.

* `nsxt_manager_endpoint` - The endpoint for the VMware NSX Manager.

* `vcsa_endpoint` - The endpoint for VMware vCenter Server Appliance.

* `sku_name` - The Name of the SKU used for this Azure VMware Solution Private Cloud.

* `nsxt_certificate_thumbprint` - The thumbprint of the VMware NSX Manager SSL certificate.

* `vcenter_certificate_thumbprint` - The thumbprint of the VMware vCenter Server SSL certificate.

* `management_subnet_cidr` - The network used to access VMware vCenter Server and NSX Manager.

* `provisioning_subnet_cidr` - The network which isused for virtual machine cold migration, cloning, and snapshot migration.

* `vmotion_subnet_cidr` - The network which is used for live migration of virtual machines.

* `tags` - A mapping of tags assigned to the Azure VMware Solution Private Cloud.

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the VMware Private Cloud.
