---
subcategory: "Cdn"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_afd_origin_group"
description: |-
  Manages a cdn AFDOriginGroup.
---

# azurerm_cdn_afd_origin_group

Manages a cdn AFDOriginGroup.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-cdn"
  location = "West Europe"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_afd_origin_group" "test" {
  name           = "acctest-c-%d"
  cdn_profile_id = azurerm_cdn_profile.test.id
  health_probe_settings {
    probe_interval_in_seconds = 0
    probe_path                = ""
    probe_protocol            = ""
    probe_request_type        = ""
  }
  load_balancing_settings {
    additional_latency_in_milliseconds = 0
    sample_size                        = 0
    successful_samples_required        = 0
  }
  response_based_afd_origin_error_detection_settings {
    http_error_ranges {
      begin = 0
      end   = 0
    }
    response_based_detected_error_types          = ""
    response_based_failover_threshold_percentage = 0
  }
  session_affinity_state                                         = ""
  traffic_restoration_time_to_healed_or_new_endpoints_in_minutes = 0
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cdn AFDOriginGroup. Changing this forces a new Cdn AFDOriginGroup to be created.

* `cdn_profile_id` - (Required) The ID of the cdn AFDOriginGroup. Changing this forces a new cdn AFDOriginGroup to be created.

* `health_probe_settings` - (Optional) A `health_probe_settings` block as defined below.

* `load_balancing_settings` - (Optional) A `load_balancing_settings` block as defined below.

* `response_based_afd_origin_error_detection_settings` - (Optional) A `response_based_afd_origin_error_detection_settings` block as defined below.

* `session_affinity_state` - (Optional) Whether to allow session affinity on this host. Valid options are 'Enabled' or 'Disabled'

* `traffic_restoration_time_to_healed_or_new_endpoints_in_minutes` - (Optional) Time in minutes to shift the traffic to the endpoint gradually when an unhealthy endpoint comes healthy or a new endpoint is added. Default is 10 mins. This property is currently not supported.

---

A `health_probe_settings` block supports the following:

* `probe_interval_in_seconds` - (Optional) The number of seconds between health probes.Default is 240sec.

* `probe_path` - (Optional) The path relative to the origin that is used to determine the health of the origin.

* `probe_protocol` - (Optional) Protocol to use for health probe.

* `probe_request_type` - (Optional) The type of health probe request that is made.

---

A `load_balancing_settings` block supports the following:

* `additional_latency_in_milliseconds` - (Optional) The additional latency in milliseconds for probes to fall into the lowest latency bucket

* `sample_size` - (Optional) The number of samples to consider for load balancing decisions

* `successful_samples_required` - (Optional) The number of samples within the sample period that must succeed

---

A `response_based_afd_origin_error_detection_settings` block supports the following:

* `http_error_ranges` - (Optional) A `http_error_ranges` block as defined below.

* `response_based_detected_error_types` - (Optional) Type of response errors for real user requests for which origin will be deemed unhealthy

* `response_based_failover_threshold_percentage` - (Optional) The percentage of failed requests in the sample where failover should trigger.

---

A `http_error_ranges` block supports the following:

* `begin` - (Optional) The inclusive start of the http status code range.

* `end` - (Optional) The inclusive end of the http status code range.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the cdn AFDOriginGroup.

* `deployment_status` - 

* `profile_name` - The name of the profile which holds the origin group.

* `provisioning_state` - Provisioning status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the cdn AFDOriginGroup.
* `read` - (Defaults to 5 minutes) Used when retrieving the cdn AFDOriginGroup.
* `update` - (Defaults to 30 minutes) Used when updating the cdn AFDOriginGroup.
* `delete` - (Defaults to 30 minutes) Used when deleting the cdn AFDOriginGroup.

## Import

cdn AFDOriginGroups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_afd_origin_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CDN/profiles/profile1/originGroups/originGroup1
```
