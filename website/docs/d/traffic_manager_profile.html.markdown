---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_profile"
description: |-
  Gets information about a specified Traffic Manager Profile.

---

# Data Source: azurerm_traffic_manager_profile

Use this data source to access information about an existing Traffic Manager Profile.

## Example Usage

```hcl
data "azurerm_traffic_manager_profile" "example" {
  name                = "test"
  resource_group_name = "test"
}

output "traffic_routing_method" {
  value = data.azurerm_traffic_manager_profile.traffic_routing_method
}
```

## Argument Reference

* `name` - Specifies the name of the Traffic Manager Profile.

* `resource_group_name` - Specifies the name of the resource group the Traffic Manager Profile is located in.

## Attributes Reference

* `id` - The ID of the Traffic Manager Profile.

* `location` - The Azure location where the Traffic Manager Profile exists.

* `fqdn` - The FQDN of the created Profile.

* `profile_status` - The status of the profile.

* `traffic_routing_method` - Specifies the algorithm used to route traffic.

* `traffic_view_enabled` - Indicates whether Traffic View is enabled for the Traffic Manager profile.

* `dns_config` - This block specifies the DNS configuration of the Profile.

* `monitor_config` - This block specifies the Endpoint monitoring configuration for the Profile.

* `tags` - A mapping of tags to assign to the resource.

The `dns_config` block provides:

* `relative_name` - The relative domain name, this is combined with the domain name used by Traffic Manager to form the FQDN which is exported as documented below.

* `ttl` - The TTL value of the Profile used by Local DNS resolvers and clients.

The `monitor_config` block provides:

* `protocol` - The protocol used by the monitoring checks.

* `port` - The port number used by the monitoring checks.

* `path` - The path used by the monitoring checks.

* `expected_status_code_ranges` - A list of status code ranges.

* `custom_header` - One or more `custom_header` blocks as defined below.

* `interval_in_seconds` - The interval used to check the endpoint health from a Traffic Manager probing agent.

* `timeout_in_seconds` - The amount of time the Traffic Manager probing agent should wait before considering that check a failure when a health check probe is sent to the endpoint.

* `tolerated_number_of_failures` - The number of failures a Traffic Manager probing agent tolerates before marking that endpoint as unhealthy.

A `custom_header` block supports the following:

* `name` - The name of the custom header.

* `value` - The value of custom header. Applicable for Http and Https protocol.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Location.
