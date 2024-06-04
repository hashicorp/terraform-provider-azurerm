---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_default_origin_group"
description: |-
  Manages the default origin group for a CDN Endpoint.
---

# azurerm_cdn_default_origin_group

Manages the default origin group for a CDN Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_profile" "example" {
  name                = "example-cdn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "example" {
  name                = "example"
  profile_name        = azurerm_cdn_profile.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  origin {
    name      = "example"
    host_name = "www.contoso.com"
  }
}

resource "azurerm_cdn_origin_groups" "example" {
    endpoint_id             = azurerm_cdn_endpoint.example.id
    default_origin_group_id = azurerm_cdn_origin_groups.example.origin.[0].id
}

resource "azurerm_cdn_default_origin_group" "example" {
    endpoint_id             = azurerm_cdn_endpoint.example.id
    default_origin_group_id = azurerm_cdn_origin_groups.example.origin.[0].id
}
```

## Arguments Reference

The following arguments are supported:

* `endpoint_id` - (Required) Specifies the resource ID of the CDN Endpoint. Changing this forces a new resource to be created.

* `default_origin_group_id` - (Required) The resource ID of the CDN origin group which should be used as the default origin group for the CDN Endpoint.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN Default Origin Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN Default Origin Group.
* `update` - (Defaults to 30 minutes) Used when updating the CDN Default Origin Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Default Origin Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN Default Origin Group.

## Import

CDN Default Origin Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_default_origin_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/myprofile1/endpoints/myendpoint1
```
