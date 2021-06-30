---
subcategory: "VMware (AVS)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vmware_express_route_authorization"
description: |-
  Manages an Express Route Vmware Authorization.
---

# azurerm_vmware_express_route_authorization

Manages an Express Route Vmware Authorization.

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

resource "azurerm_vmware_express_route_authorization" "example" {
  name             = "example-authorization"
  private_cloud_id = azurerm_vmware_private_cloud.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Express Route Vmware Authorization. Changing this forces a new Vmware Authorization to be created.

* `private_cloud_id` - (Required) The ID of the Vmware Private Cloud in which to create this Express Route Vmware Authorization. Changing this forces a new Vmware Authorization to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Vmware Authorization.

* `express_route_authorization_id` - The ID of the Express Route Circuit Authorization.

* `express_route_authorization_key` - The key of the Express Route Circuit Authorization.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Vmware Authorization.
* `read` - (Defaults to 5 minutes) Used when retrieving the Vmware Authorization.
* `delete` - (Defaults to 30 minutes) Used when deleting the Vmware Authorization.

## Import

Vmware Authorizations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vmware_express_route_authorization.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AVS/privateClouds/privateCloud1/authorizations/authorization1
```
