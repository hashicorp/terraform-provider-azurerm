---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_avs_authorization"
description: |-
  Manages a avs Authorization.
---

# azurerm_avs_authorization

Manages a avs Authorization.

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
  network_block = ""
}

resource "azurerm_avs_authorization" "example" {
  name = "example-authorization"
  resource_group_name = azurerm_resource_group.example.name
  private_cloud_name = azurerm_avs_private_cloud.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this avs Authorization. Changing this forces a new avs Authorization to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the avs Authorization should exist. Changing this forces a new avs Authorization to be created.

* `private_cloud_name` - (Required) The name of the private cloud. Changing this forces a new avs Authorization to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the avs Authorization.

* `express_route_authorization_id` - The ID of the express_route_authorization.

* `express_route_authorization_key` - The key of the ExpressRoute Circuit Authorization.

* `type` - Resource type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the avs Authorization.
* `read` - (Defaults to 5 minutes) Used when retrieving the avs Authorization.
* `delete` - (Defaults to 30 minutes) Used when deleting the avs Authorization.

## Import

avs Authorizations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_avs_authorization.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/privateCloud1/authorizations/authorization1
```