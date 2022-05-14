---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin_group"
description: |-
  Manages a Frontdoor Origin Group.
---

# azurerm_cdn_frontdoor_origin_group

Manages a Frontdoor Origin Group.

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

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                         = "example-originGroup"
  cdn_cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  session_affinity_enabled     = true

  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = 10

  health_probe {
    interval_in_seconds = 240
    path                = "/healthProbe"
    protocol            = "Https"
    request_type        = "GET"
  }

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Origin Group. Possible values must be between 2 and 90 characters in length, begin with a letter or number, end with a letter or number and may contain only letters, numbers and hyphens. Changing this forces a new Frontdoor Origin Group to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the Frontdoor Profile. Changing this forces a new Frontdoor Origin Group to be created.

* `load_balancing` - (Required) A `load_balancing` block as defined below.

* `health_probe` - (Optional) A `health_probe` block as defined below.

* `session_affinity_enabled` - (Optional) Whether to allow session affinity on this host. Possible values are `true` or `false`. Defaults to `true`.

* `restore_traffic_time_to_healed_or_new_endpoint_in_minutes` - (Optional) Time in minutes to shift the traffic to the endpoint gradually when an unhealthy endpoint becomes healthy or a new endpoint is added. Default is `10` minutes.

~> **NOTE:** This property is currently not supported.

---

A `health_probe` block supports the following:

* `protocol` - (Required) Protocol to use for health probe. Possible values are `Http` or `Https`.

* `request_type` - (Required) The type of health probe request that is made. Possible values are `GET` or `HEAD`.

* `interval_in_seconds` - (Optional) The number of seconds between health probes. Default is `100` seconds. Possible values are between `5` and `31536000` seconds(inclusive).

* `path` - (Optional) The path relative to the origin that is used to determine the health of the origin. Defaults to `/`.

---

A `load_balancing` block supports the following:

* `additional_latency_in_milliseconds` - (Optional) The additional latency in milliseconds for probes to fall into the lowest latency bucket. Possible values are between `0` and `1000` seconds(inclusive). Defaults to `50`.

* `sample_size` - (Optional) The number of samples to consider for load balancing decisions. Possible values are between `0` and `255`(inclusive). Defaults to `4`.

* `successful_samples_required` - (Optional) The number of samples within the sample period that must succeed. Possible values are between `0` and `255`(inclusive). Defaults to `3`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Origin Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Origin Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Origin Group.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Origin Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Origin Group.

## Import

Frontdoor Origin Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_origin_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1
```
