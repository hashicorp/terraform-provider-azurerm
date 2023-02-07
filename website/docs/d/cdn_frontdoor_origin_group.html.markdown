---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin_group"
description: |-
  Gets information about an existing Front Door (standard/premium) Origin Group.
---

# Data Source: azurerm_cdn_frontdoor_origin_group

Use this data source to access information about an existing Front Door (standard/premium) Origin Group.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_origin_group" "example" {
  name                = "example-origin-group"
  profile_name        = "example-profile"
  resource_group_name = "example-resources"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door Origin Group.

* `profile_name` - (Required) The name of the Front Door Profile within which Front Door Origin Group exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Front Door Profile exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Origin Group.

* `cdn_frontdoor_profile_id` - Specifies the ID of the Front Door Profile within which this Front Door Origin Group exists.

* `health_probe` - A `health_probe` block as defined below.

* `load_balancing` - A `load_balancing` block as defined below.

* `session_affinity_enabled` - Specifies whether session affinity is enabled on this host.

---

A `health_probe` block exports the following:

* `protocol` - Specifies the protocol to use for health probe.

* `request_type` - Specifies the type of health probe request that is made.

* `interval_in_seconds` - Specifies the number of seconds between health probes.

* `path` - Specifies the path relative to the origin that is used to determine the health of the origin.

---

A `load_balancing` block exports the following:

* `additional_latency_in_milliseconds` - Specifies the additional latency in milliseconds for probes to fall into the lowest latency bucket.

* `sample_size` - Specifies the number of samples to consider for load balancing decisions.

* `successful_samples_required` - Specifies the number of samples within the sample period that must succeed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Origin Group.
