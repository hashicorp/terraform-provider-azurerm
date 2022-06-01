---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_profile"
description: |-
  Manages a CDN FrontDoor Profile which contains a collection of CDN FrontDoor Endpoints.
---

# azurerm_cdn_frontdoor_profile

Manages a CDN FrontDoor Profile which contains a collection of CDN FrontDoor Endpoints.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-cdn-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the FrontDoor Profile. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where this FrontDoor Profile should exist. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU for this CDN FrontDoor Profile. Possible values include `Standard_AzureFrontDoor` and `Premium_AzureFrontDoor`. Changing this forces a new resource to be created.

* `response_timeout_seconds` - Specifies the maximum response timeout in seconds. Possible values are between `16` and `240` seconds (inclusive). Defaults to `120` seconds.

* `tags` - (Optional) Specifies a mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this CDN FrontDoor Profile.

* `resource_guid` - The UUID of this CDN FrontDoor Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN FrontDoor Profile.
* `update` - (Defaults to 30 minutes) Used when updating the CDN FrontDoor Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN FrontDoor Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN FrontDoor Profile.

## Import

CDN FrontDoor Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/myprofile1
```
