---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_avs_private_cloud"
description: |-
  Manages a avs PrivateCloud.
---

# azurerm_avs_private_cloud

Manages a avs PrivateCloud.

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
  network_block = "192.168.48.0/22"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this avs PrivateCloud. Changing this forces a new avs PrivateCloud to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the avs PrivateCloud should exist. Changing this forces a new avs PrivateCloud to be created.

* `location` - (Required) The Azure Region where the avs PrivateCloud should exist. Changing this forces a new avs PrivateCloud to be created.

* `sku` - (Required)  A `sku` block as defined below. Changing this forces a new avs PrivateCloud to be created.

* `management_cluster` - (Required)  A `management_cluster` block as defined below.

* `network_block` - (Required) The block of addresses should be unique across VNet in your subscription as well as on-premise. Make sure the CIDR format is conformed to (A.B.C.D/X) where A,B,C,D are between 0 and 255, and X is between 0 and 22.

---

* `identity_source` - (Optional)  A `identity_source` block as defined below.

* `internet` - (Optional) Connectivity to internet is enabled or disabled. Possible values are "true" and "false" is allowed.

* `nsxt_password` - (Optional) Optionally, set the NSX-T Manager password when the private cloud is created.

* `vcenter_password` - (Optional) Optionally, set the vCenter admin password when the private cloud is created.

* `tags` - (Optional) A mapping of tags which should be assigned to the avs PrivateCloud.

---

An `sku` block exports the following:

* `name` - (Required) The name which should be used for this sku. Changing this forces a new avs PrivateCloud to be created.

---

An `management_cluster` block exports the following:

* `cluster_size` - (Required) The cluster size.

---

An `identity_source` block exports the following:

* `name` - (Optional) The name which should be used for this identity_source.

* `alias` - (Optional) The domain's NetBIOS name.

* `base_group_dn` - (Optional) The base distinguished name for groups.

* `base_user_dn` - (Optional) The base distinguished name for users.

* `domain` - (Optional) The domain's dns name.

* `password` - (Optional) The password of the Active Directory user with a minimum of read-only access to Base DN for users and groups.

* `primary_server` - (Optional) Primary server URL.

* `secondary_server` - (Optional) Secondary server URL.

* `ssl` - (Optional) Protect LDAP communication using SSL certificate (LDAPS). Possible values are "true" and "false" is allowed.

* `username` - (Optional) The ID of an Active Directory user with a minimum of read-only access to Base DN for users and group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the avs PrivateCloud.

* `circuit` - An ExpressRoute Circuit. A `circuit` block as defined below.

* `endpoint` - The endpoints. A `endpoint` block as defined below.

* `management_network` - Network used to access vCenter Server and NSX-T Manager.

* `nsxt_certificate_thumbprint` - Thumbprint of the NSX-T Manager SSL certificate.

* `provisioning_network` - Used for virtual machine cold migration, cloning, and snapshot migration.

* `type` - Resource type.

* `vcenter_certificate_thumbprint` - Thumbprint of the vCenter Server SSL certificate.

* `vmotion_network` - Used for live migration of virtual machines.

---

An `circuit` block exports the following:

* `express_route_id` - The ID of the express_route.

* `express_route_private_peering_id` - The ID of the express_route_private_peering.

* `primary_subnet` - CIDR of primary subnet.

* `secondary_subnet` - CIDR of secondary subnet.

---

An `endpoint` block exports the following:

* `hcx_cloud_manager` - Endpoint for the HCX Cloud Manager.

* `nsxt_manager` - Endpoint for the NSX-T Data Center manager.

* `vcsa` - Endpoint for Virtual Center Server Appliance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the avs PrivateCloud.
* `read` - (Defaults to 5 minutes) Used when retrieving the avs PrivateCloud.
* `update` - (Defaults to 30 minutes) Used when updating the avs PrivateCloud.
* `delete` - (Defaults to 30 minutes) Used when deleting the avs PrivateCloud.

## Import

avs PrivateClouds can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_avs_private_cloud.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/privateCloud1
```