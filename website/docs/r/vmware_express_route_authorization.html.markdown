---
subcategory: "Azure VMware Solution"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vmware_express_route_authorization"
description: |-
  Manages an Azure VMware Solution ExpressRoute Circuit Authorization.
---

# azurerm_vmware_express_route_authorization

Manages an Azure VMware Solution ExpressRoute Circuit Authorization.

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

* `name` - (Required) The name which should be used for this Azure VMware Solution ExpressRoute Circuit Authorization. Changing this forces a new Azure VMware Solution ExpressRoute Circuit Authorization to be created.

* `private_cloud_id` - (Required) The ID of the Azure VMware Solution Private Cloud in which to create this Azure VMware Solution ExpressRoute Circuit Authorization. Changing this forces a new Azure VMware Solution ExpressRoute Circuit Authorization to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure VMware Solution ExpressRoute Circuit Authorization.

* `express_route_authorization_id` - The ID of the Azure VMware Solution ExpressRoute Circuit Authorization.

* `express_route_authorization_key` - The key of the Azure VMware Solution ExpressRoute Circuit Authorization.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure VMware Solution ExpressRoute Circuit Authorization.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure VMware Solution ExpressRoute Circuit Authorization.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure VMware Solution ExpressRoute Circuit Authorization.

## Import

Azure VMware Solution ExpressRoute Circuit Authorizations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vmware_express_route_authorization.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AVS/privateClouds/privateCloud1/authorizations/authorization1
```
