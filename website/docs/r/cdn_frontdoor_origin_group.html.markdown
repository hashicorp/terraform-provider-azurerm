---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin_group"
description: |-
  Manages a Front Door (standard/premium) Origin Group.
---

# azurerm_cdn_frontdoor_origin_group

Manages a Front Door (standard/premium) Origin Group.

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

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "example-origin-group"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  session_affinity_enabled = true

  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = 10

  health_probe {
    interval_in_seconds = 240
    path                = "/healthProbe"
    protocol            = "Https"
    request_type        = "HEAD"
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

* `name` - (Required) The name which should be used for this Front Door Origin Group. Changing this forces a new Front Door Origin Group to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the Front Door Profile within which this Front Door Origin Group should exist. Changing this forces a new Front Door Origin Group to be created.

* `load_balancing` - (Required) A `load_balancing` block as defined below.

---

* `health_probe` - (Optional) A `health_probe` block as defined below.

* `restore_traffic_time_to_healed_or_new_endpoint_in_minutes` - (Optional) Specifies the amount of time which should elapse before shifting traffic to another endpoint when a healthy endpoint becomes unhealthy or a new endpoint is added. Possible values are between `0` and `50` minutes (inclusive). Default is `10` minutes.

-> **Note:** This property is currently not used, but will be in the near future.

* `session_affinity_enabled` - (Optional) Specifies whether session affinity should be enabled on this host. Defaults to `true`.

---

A `health_probe` block supports the following:

* `protocol` - (Required) Specifies the protocol to use for health probe. Possible values are `Http` and `Https`.

* `interval_in_seconds` - (Required) Specifies the number of seconds between health probes. Possible values are between `1` and `255` seconds (inclusive).

* `request_type` - (Optional) Specifies the type of health probe request that is made. Possible values are `GET` and `HEAD`. Defaults to `HEAD`.

* `path` - (Optional) Specifies the path relative to the origin that is used to determine the health of the origin. Defaults to `/`.

-> **Note:** Health probes can only be disabled if there is a single enabled origin in a single enabled origin group. For more information about the `health_probe` settings please see the [product documentation](https://docs.microsoft.com/azure/frontdoor/health-probes).

---

A `load_balancing` block supports the following:

* `additional_latency_in_milliseconds` - (Optional) Specifies the additional latency in milliseconds for probes to fall into the lowest latency bucket. Possible values are between `0` and `1000` milliseconds (inclusive). Defaults to `50`.

* `sample_size` - (Optional) Specifies the number of samples to consider for load balancing decisions. Possible values are between `0` and `255` (inclusive). Defaults to `4`.

* `successful_samples_required` - (Optional) Specifies the number of samples within the sample period that must succeed. Possible values are between `0` and `255` (inclusive). Defaults to `3`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Origin Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Origin Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Origin Group.
* `update` - (Defaults to 30 minutes) Used when updating the Front Door Origin Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Origin Group.

## Import

Front Door Origin Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_origin_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1
```
