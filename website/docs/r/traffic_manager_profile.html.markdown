---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_profile"
description: |-
  Manages a Traffic Manager Profile.

---

# azurerm_traffic_manager_profile

Manages a Traffic Manager Profile to which multiple endpoints can be attached.

## Example Usage

```hcl
resource "random_id" "server" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_resource_group" "example" {
  name     = "trafficmanagerProfile"
  location = "West Europe"
}

resource "azurerm_traffic_manager_profile" "example" {
  name                   = random_id.server.hex
  resource_group_name    = azurerm_resource_group.example.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = random_id.server.hex
    ttl           = 100
  }

  monitor_config {
    protocol                     = "HTTP"
    port                         = 80
    path                         = "/"
    interval_in_seconds          = 30
    timeout_in_seconds           = 9
    tolerated_number_of_failures = 3
  }

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Traffic Manager profile. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the resource group in which to create the Traffic Manager profile. Changing this forces a new resource to be created.

* `profile_status` - (Optional) The status of the profile, can be set to either `Enabled` or `Disabled`. Defaults to `Enabled`.

* `traffic_routing_method` - (Required) Specifies the algorithm used to route traffic. Possible values are `Geographic`, `Weighted`, `Performance`, `Priority`, `Subnet` and `MultiValue`.
  * `Geographic` - Traffic is routed based on Geographic regions specified in the Endpoint.
  * `MultiValue` - All healthy Endpoints are returned.  MultiValue routing method works only if all the endpoints of type `External` and are specified as IPv4 or IPv6 addresses.
  * `Performance` - Traffic is routed via the User's closest Endpoint
  * `Priority` - Traffic is routed to the Endpoint with the lowest `priority` value.
  * `Subnet` - Traffic is routed based on a mapping of sets of end-user IP address ranges to a specific Endpoint within a Traffic Manager profile.
  * `Weighted` - Traffic is spread across Endpoints proportional to their `weight` value.

* `traffic_view_enabled` - (Optional) Indicates whether Traffic View is enabled for the Traffic Manager profile.

* `dns_config` - (Required) This block specifies the DNS configuration of the Profile. One `dns_config` block as defined below.

* `monitor_config` - (Required) This block specifies the Endpoint monitoring configuration for the Profile. One `monitor_config` block as defined below.

* `max_return` - (Optional) The amount of endpoints to return for DNS queries to this Profile. Possible values range from `1` to `8`.

~> **Note:** `max_return` must be set when the `traffic_routing_method` is `MultiValue`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `dns_config` block supports:

* `relative_name` - (Required) The relative domain name, this is combined with the domain name used by Traffic Manager to form the FQDN which is exported as documented below. Changing this forces a new resource to be created.

* `ttl` - (Required) The TTL value of the Profile used by Local DNS resolvers and clients.

---

The `monitor_config` block supports:

* `protocol` - (Required) The protocol used by the monitoring checks, supported values are `HTTP`, `HTTPS` and `TCP`.

* `port` - (Required) The port number used by the monitoring checks.

* `path` - (Optional) The path used by the monitoring checks. Required when `protocol` is set to `HTTP` or `HTTPS` - cannot be set when `protocol` is set to `TCP`.

* `expected_status_code_ranges` - (Optional) A list of status code ranges in the format of `100-101`.

* `custom_header` - (Optional) One or more `custom_header` blocks as defined below.

* `interval_in_seconds` - (Optional) The interval used to check the endpoint health from a Traffic Manager probing agent. You can specify two values here: `30` (normal probing) and `10` (fast probing). The default value is `30`.

* `timeout_in_seconds` - (Optional) The amount of time the Traffic Manager probing agent should wait before considering that check a failure when a health check probe is sent to the endpoint. If `interval_in_seconds` is set to `30`, then `timeout_in_seconds` can be between `5` and `10`. The default value is `10`. If `interval_in_seconds` is set to `10`, then valid values are between `5` and `9` and `timeout_in_seconds` is required.

* `tolerated_number_of_failures` - (Optional) The number of failures a Traffic Manager probing agent tolerates before marking that endpoint as unhealthy. Valid values are between `0` and `9`. The default value is `3`

---

A `custom_header` block supports the following:

* `name` - (Required) The name of the custom header.

* `value` - (Required) The value of custom header. Applicable for HTTP and HTTPS protocol.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Traffic Manager Profile.

* `fqdn` - The FQDN of the created Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Traffic Manager Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Traffic Manager Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Traffic Manager Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Traffic Manager Profile.

## Import

Traffic Manager Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_traffic_manager_profile.exampleProfile /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/trafficManagerProfiles/mytrafficmanagerprofile1
```
