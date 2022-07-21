---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_endpoint"
description: |-
  Manages a FrontDoor Endpoint.
---

# azurerm_cdn_frontdoor_endpoint

Manages a FrontDoor Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"

}

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                     = "example-endpoint"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  tags = {
    ENV = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this CDN FrontDoor Endpoint. Changing this forces a new CDN FrontDoor Endpoint to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the FrontDoor Profile within which this FrontDoor Endpoint should exist. Changing this forces a new CDN FrontDoor Endpoint to be created.

---

* `enabled` - (Optional) Specifies if this CDN FrontDoor Endpoint is enabled? Defaults to `true`.

* `tags` - (Optional) Specifies a mapping of tags which should be assigned to the CDN FrontDoor Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this CDN FrontDoor Endpoint.

* `host_name` - The host name of the CDN FrontDoor Endpoint, in the format `{endpointName}.{dnsZone}` (for example, `contoso.azureedge.net`).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN FrontDoor Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN FrontDoor Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the CDN FrontDoor Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN FrontDoor Endpoint.

## Import

CDN FrontDoor Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1
```
