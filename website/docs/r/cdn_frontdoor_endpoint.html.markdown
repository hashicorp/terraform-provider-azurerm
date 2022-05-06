---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_endpoint"
description: |-
  Manages a Frontdoor Endpoint.
---

# azurerm_cdn_frontdoor_endpoint

Manages a Frontdoor Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                            = "example-endpoint"
  cdn_frontdoor_profile_id        = azurerm_cdn_frontdoor_profile.example.id
  enabled                         = true
  origin_response_timeout_seconds = 120

  tags = {
    ENV = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Endpoint. Changing this forces a new Frontdoor Endpoint to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the Frontdoor Profile. Changing this forces a new Frontdoor Endpoint to be created.

* `enabled` - (Optional) Should this Frontdoor Endpoint be used? Possible values include `true` or `false`. Defaults to `true`.

* `origin_response_timeout_seconds` - (Optional) Send and receive timeout on forwarding request to the origin. When timeout is reached, the request fails and returns. Defaults to `120` seconds.

~> **NOTE:** Due to a bug in the service code the `origin_response_timeout_seconds` will always be set to the default value of `120` seconds.

* `tags` - (Optional) A mapping of tags which should be assigned to the Frontdoor Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Endpoint.

* `host_name` - The host name of the Frontdoor Endpoint structured as `\[endpointName\].\[DNSZone\]`(e.g. contoso.azureedge.net).

* `frontdoor_profile_name` - The name of the Frontdoor Profile which holds the endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Endpoint.

## Import

Frontdoor Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1
```
