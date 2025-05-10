---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_nested_endpoint"
description: |-
  Manages a Nested Endpoint within a Traffic Manager Profile.
---

# azurerm_traffic_manager_nested_endpoint

Manages a Nested Endpoint within a Traffic Manager Profile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "example-publicip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  domain_name_label   = "example-pip"
}


resource "azurerm_traffic_manager_profile" "parent" {
  name                   = "parent-profile"
  resource_group_name    = azurerm_resource_group.example.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "parent-profile"
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

resource "azurerm_traffic_manager_profile" "nested" {
  name                   = "nested-profile"
  resource_group_name    = azurerm_resource_group.example.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "nested-profile"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTP"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_nested_endpoint" "example" {
  name                    = "example-endpoint"
  target_resource_id      = azurerm_traffic_manager_profile.nested.id
  priority                = 1
  profile_id              = azurerm_traffic_manager_profile.parent.id
  minimum_child_endpoints = 9
  weight                  = 5
}
```

## Argument Reference

The following arguments are supported:

* `minimum_child_endpoints` - (Required) This argument specifies the minimum number of endpoints that must be ‘online’ in the child profile in order for the parent profile to direct traffic to any of the endpoints in that child profile. This value must be larger than `0`.

~> **Note:** If `min_child_endpoints` is less than either `minimum_required_child_endpoints_ipv4` or `minimum_required_child_endpoints_ipv6`, then it won't have any effect.

* `name` - (Required) The name of the External Endpoint. Changing this forces a new resource to be created.

* `profile_id` - (Required) The ID of the Traffic Manager Profile that this External Endpoint should be created within. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The resource id of an Azure resource to target.

* `weight` - (Optional) Specifies how much traffic should be distributed to this endpoint, this must be specified for Profiles using the Weighted traffic routing method. Valid values are between `1` and `1000`. Defaults to `1`.

---

* `custom_header` - (Optional) One or more `custom_header` blocks as defined below.

* `enabled` - (Optional) Is the endpoint enabled? Defaults to `true`.

* `endpoint_location` - (Optional) Specifies the Azure location of the Endpoint, this must be specified for Profiles using the `Performance` routing method.

* `minimum_required_child_endpoints_ipv4` - (Optional) This argument specifies the minimum number of IPv4 (DNS record type A) endpoints that must be ‘online’ in the child profile in order for the parent profile to direct traffic to any of the endpoints in that child profile. This argument only applies to Endpoints of type `nestedEndpoints` and 

* `minimum_required_child_endpoints_ipv6` - (Optional) This argument specifies the minimum number of IPv6 (DNS record type AAAA) endpoints that must be ‘online’ in the child profile in order for the parent profile to direct traffic to any of the endpoints in that child profile. This argument only applies to Endpoints of type `nestedEndpoints` and 

* `priority` - (Optional) Specifies the priority of this Endpoint, this must be specified for Profiles using the `Priority` traffic routing method. Supports values between 1 and 1000, with no Endpoints sharing the same value. If omitted the value will be computed in order of creation.

* `geo_mappings` - (Optional) A list of Geographic Regions used to distribute traffic, such as `WORLD`, `UK` or `DE`. The same location can't be specified in two endpoints. [See the Geographic Hierarchies documentation for more information](https://docs.microsoft.com/rest/api/trafficmanager/geographichierarchies/getdefault).

* `subnet` - (Optional) One or more `subnet` blocks as defined below. Changing this forces a new resource to be created.

---

A `custom_header` block supports the following:

* `name` - (Required) The name of the custom header.

* `value` - (Required) The value of custom header. Applicable for HTTP and HTTPS protocol.

---

A `subnet` block supports the following:

* `first` - (Required) The first IP Address in this subnet.

* `last` - (Optional) The last IP Address in this subnet.

* `scope` - (Optional) The block size (number of leading bits in the subnet mask).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Nested Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Nested Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Nested Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Nested Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Nested Endpoint.

## Import

Nested Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_traffic_manager_nested_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Network/trafficManagerProfiles/example-profile/NestedEndpoints/example-endpoint
```
