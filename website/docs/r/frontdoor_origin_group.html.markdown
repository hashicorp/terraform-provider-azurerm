---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: frontdoor_origin_group"
description: |-
  Manages a Frontdoor Origin Group.
---

# frontdoor_origin_group

Manages a Frontdoor Origin Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-cdn"
  location = "West Europe"
}

resource "azurerm_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "frontdoor_origin_group" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id

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

  response_based_origin_error_detection {
    http_error_ranges {
      begin = 300
      end   = 599
    }

    detected_error_types          = "TcpAndHttpErrors"
    failover_threshold_percentage = 10
  }

  session_affinity                      = true
  restore_traffic_or_new_endpoints_time = 10
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Origin Group. Changing this forces a new Frontdoor Origin Group to be created.

* `frontdoor_profile_id` - (Required) The ID of the Frontdoor Profile. Changing this forces a new Frontdoor Origin Group to be created.

* `health_probe` - (Optional) A `health_probe` block as defined below.

* `load_balancing` - (Optional) A `load_balancing` block as defined below.

* `response_based_origin_error_detection` - (Optional) A `response_based_origin_error_detection` block as defined below.

* `session_affinity` - (Optional) Whether to allow session affinity on this host. Possible values are `true` or `false`. Defaults to `true`.

* `restore_traffic_or_new_endpoints_time` - (Optional) Time in minutes to shift the traffic to another endpoint when an healthy endpoint becomes unhealthy or a new endpoint is added. Default is `10` minutes.

~> **NOTE:** This property is currently not used, but will be in the near future.

---

A `health_probe` block supports the following:

* `interval_in_seconds` - (Optional) The number of seconds between health probes. Default is `240` seconds. Possible values are between `5` and `31536000` seconds(inclusive).

* `path` - (Optional) The path relative to the origin that is used to determine the health of the origin. Defaults to `/`.

* `protocol` - (Optional) Protocol to use for health probe. Possible values are `Http`, `Https` or `NotSet`. Defaults to `Https`.

* `request_type` - (Optional) The type of health probe request that is made. Possible values are `GET`, `HEAD` or `NotSet`. Defaults to `GET`.

---

A `load_balancing` block supports the following:

* `additional_latency_in_milliseconds` - (Optional) The additional latency in milliseconds for probes to fall into the lowest latency bucket. Possible values are between `0` and `1000` seconds(inclusive). Defaults to `0`.

* `sample_size` - (Optional) The number of samples to consider for load balancing decisions. Possible values are between `0` and `255`(inclusive). Defaults to `16`.

* `successful_samples_required` - (Optional) The number of samples within the sample period that must succeed. Possible values are between `0` and `255`(inclusive). Defaults to `3`.

---

A `response_based_origin_error_detection` block supports the following:

* `http_error_ranges` - (Optional) A `http_error_ranges` block as defined below.

* `detected_error_types` - (Optional) Type of response errors for real user requests for which origin will be deemed unhealthy. Possible values are `None`, `TcpAndHttpErrors` or `TcpErrorsOnly`. Defaults to `TcpAndHttpErrors`.

* `failover_threshold_percentage` - (Optional) The percentage of failed requests in the sample where failover should trigger. Possible values are between `0` and `100`(inclusive). Defaults to `10`.

---

A `http_error_ranges` block supports the following:

* `begin` - (Optional) The inclusive start of the http status code range. Possible values are between `100` and `999`(inclusive). Defaults to `300`.

* `end` - (Optional) The inclusive end of the http status code range. Possible values are between `100` and `999`(inclusive). Defaults to `599`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Origin Group.

* `deployment_status` - Deployment status

* `frontdoor_profile_name` - The name of the profile which holds the origin group.

* `provisioning_state` - Provisioning status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Origin Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Origin Group.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Origin Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Origin Group.

## Import

Frontdoor Origin Groups can be imported using the `resource id`, e.g.

```shell
terraform import frontdoor_origin_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1
```
