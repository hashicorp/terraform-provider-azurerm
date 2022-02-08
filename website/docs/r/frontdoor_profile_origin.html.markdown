---
subcategory: "Cdn"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_profile_origin"
description: |-
  Manages a Frontdoor Profile Origin.
---

# azurerm_frontdoor_profile_origin

Manages a Frontdoor Profile Origin.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-frontdoor-profile"
  location = "West Europe"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_frontdoor_profile_origin_group" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
}

resource "azurerm_frontdoor_profile_origin" "test" {
  name                              = "acctest-c-%d"
  frontdoor_profile_origin_group_id = azurerm_frontdoor_profile_origin_group.test.id
  azure_origin_id                   = ""

  enabled_state                  = ""
  enforce_certificate_name_check = false
  host_name                      = ""
  http_port                      = 0
  https_port                     = 0
  origin_host_header             = ""
  priority                       = 0
  weight                         = 0
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Profile Origin. Changing this forces a new Frontdoor Profile Origin to be created.

* `frontdoor_profile_origin_group_id` - (Required) The ID of the Frontdoor Profile Origin Group. Changing this forces a new Frontdoor Profile Origin Group to be created.

* `host_name` - (Required) The address of the origin. Domain names, IPv4 addresses, and IPv6 addresses are supported.This should be unique across all origins in an endpoint.

* `azure_origin_id` - (Optional) Resource ID..

* `enabled_state` - (Optional) Whether to enable health probes to be made against backends defined under backendPools. Health probes can only be disabled if there is a single enabled backend in single enabled backend pool.

* `enforce_certificate_name_check` - (Optional) Whether to enable certificate name check at origin level

* `http_port` - (Optional) The value of the HTTP port. Must be between 1 and 65535.

* `https_port` - (Optional) The value of the HTTPS port. Must be between 1 and 65535.

* `origin_host_header` - (Optional) The host header value sent to the origin with each request. If you leave this blank, the request hostname determines this value. Azure CDN origins, such as Web Apps, Blob Storage, and Cloud Services require this host header value to match the origin hostname by default. This overrides the host header defined at Endpoint

* `priority` - (Optional) Priority of origin in given origin group for load balancing. Higher priorities will not be used for load balancing if any lower priority origin is healthy.Must be between 1 and 5

* `weight` - (Optional) Weight of the origin in given origin group for load balancing. Must be between 1 and 1000

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Profile Origin.

* `deployment_status` - 

* `origin_group_name` - The name of the origin group which contains this origin.

* `provisioning_state` - Provisioning status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Profile Origin.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Profile Origin.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Profile Origin.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Profile Origin.

## Import

Frontdoor Profile Origins can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor_profile_origin.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1/origins/origin1
```
