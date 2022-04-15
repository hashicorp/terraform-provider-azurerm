---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin"
description: |-
  Manages a Frontdoor Origin.
---

# azurerm_cdn_frontdoor_origin

Manages a Frontdoor Origin.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-frontdoor-profile"
  location = "West Europe"
}

resource "azurerm_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "example-originGroup"
  cdn_frontdoor_profile_id = azurerm_frontdoor_profile.example.id
}

resource "azurerm_cdn_frontdoor_origin" "example" {
  name                                  = "example-origin"
  cdn_frontdoor_profile_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id

  health_probes_enabled          = true
  certificate_name_check_enabled = false

  host_name          = "contoso.com"
  http_port          = 80
  https_port         = 443
  origin_host_header = "www.contoso.com"
  priority           = 1
  weight             = 1
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Origin. Changing this forces a new Frontdoor Origin to be created.

* `cdn_frontdoor_origin_group_id` - (Required) The ID of the Frontdoor Origin Group. Changing this forces a new Frontdoor Origin Group to be created.

* `host_name` - (Required) The address of the origin. Domain names, IPv4 addresses, and IPv6 addresses are supported. This should be unique across all Frontdoor Origins in an Frontdoor Endpoints.

* `cdn_frontdoor_origin_id` - (Optional) Resource ID.

* `health_probes_enabled` - (Optional) Are health probes enabled against backends defined under the backendPools? Health probes can only be disabled if there is a single enabled backend in single enabled backend pool. Possible values are `true` or `false`. Defaults to `true`.

* `certificate_name_check_enabled` - (Optional) Whether to enable certificate name check at origin level. Possible values are `true` or `false`. Defaults to `false`.

* `private_link` - (Optional) TBD.

->**NOTE:** To include a `private_link` in the Frontdoor Origin your Frontdoor Profile must be a `Premium_AzureFrontDoor` SKU and the `certificate_name_check_enabled` field must be set to `true`.

* `http_port` - (Optional) The value of the HTTP port. Must be between `1` and `65535`. Defaults to `80`.

* `https_port` - (Optional) The value of the HTTPS port. Must be between `1` and `65535`. Defaults to `443`.

* `origin_host_header` - (Optional) The host header value sent to the origin with each request. Possible values include an `IPv4` IP address, `IPv6` IP address or a valid `domain name`. If you leave field undefined, the requests hostname determines this value. Azure Frontdoor Origins, such as Web Apps, Blob Storage, and Cloud Services require this host header value to match the origins hostname. This fields value overrides the host header defined in the Frontdoor Endpoint. For more information on how to properly set the origin host header value please see the [product documentation](https://docs.microsoft.com/azure/frontdoor/origin?pivots=front-door-standard-premium#origin-host-header).

* `priority` - (Optional) Priority of origin in given origin group for load balancing. Higher priorities will not be used for load balancing if any lower priority origin is healthy.Must be between `1` and `5`(inclusive). Defaults to `1`.

* `weight` - (Optional) Weight of the origin in given origin group for load balancing. Must be between `1` and `1000`(inclusive). Defaults to `500`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Origin.

* `origin_group_name` - The name of the origin group which contains this Frontdoor Origin.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Origin.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Origin.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Origin.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Origin.

## Import

Frontdoor Origins can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_origin.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1/origins/origin1
```
