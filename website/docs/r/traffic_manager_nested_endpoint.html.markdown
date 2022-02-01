---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_nested_endpoint"
description: |-
  Manages a Traffic Manager Nested Endpoint.
---

# azurerm_traffic_manager_nested_endpoint

Manages a Traffic Manager Nested Endpoint.

## Example Usage

```hcl
resource "random_id" "server" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_resource_group" "example" {
  name     = "trafficmanagerendpointTest"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "trafficmanagerendpointTest"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "trafficmanagerendpointTest"
}


resource "azurerm_traffic_manager_profile" "parent" {
  name                = random_id.server.hex
  resource_group_name = azurerm_resource_group.example.name

  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = random_id.server.hex
    ttl           = 100
  }

  monitor_config {
    protocol                     = "http"
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

resource "azurerm_traffic_manager_profile" "child" {
  name                   = random_id.server.hex
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = random_id.server.hex
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_nested_endpoint" "test" {
  name                = "trafficManagerNestedEndpoint"
  target_resource_id  = azurerm_traffic_manager_profile.child.id
  priority            = 1
  profile_name        = azurerm_traffic_manager_profile.parent.name
  resource_group_name = azurerm_resource_group.test.name
  min_child_endpoints = 5
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Traffic Manager endpoint. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the Traffic Manager Profile exists.

* `profile_name` - (Required) The name of the Traffic Manager Profile to attach
    create the Traffic Manager endpoint.

* `enabled` - (Optional) Is the endpoint enabled? Defaults to `true`.

* `target_resource_id` - (Required) The resource id of an Azure resource to
    target.

* `weight` - (Required) Specifies how much traffic should be distributed to this
    endpoint. Valid values are between `1` and `1000`.

* `minimum_child_endpoints` - (Required) This argument specifies the minimum number
  of endpoints that must be ‘online’ in the child profile in order for the
  parent profile to direct traffic to any of the endpoints in that child
  profile. This value must be larger than `0`.

~>**NOTE:** If `min_child_endpoints` is less than either `minimum_required_child_endpoints_ipv4` or `minimum_required_child_endpoints_ipv6`, then it won't have any effect.

* `minimum_required_child_endpoints_ipv4` - (Optional) This argument specifies the minimum number of IPv4 (DNS record type A) endpoints that must be ‘online’ in the child profile in order for the parent profile to direct traffic to any of the endpoints in that child profile. This argument only applies to Endpoints of type `nestedEndpoints` and defaults to `1`.

* `minimum_required_child_endpoints_ipv6` - (Optional) This argument specifies the minimum number of IPv6 (DNS record type AAAA) endpoints that must be ‘online’ in the child profile in order for the parent profile to direct traffic to any of the endpoints in that child profile. This argument only applies to Endpoints of type `nestedEndpoints` and defaults to `1`.

* `custom_header` - (Optional) One or more `custom_header` blocks as defined below.

* `priority` - (Optional) Specifies the priority of this Endpoint, this must be
    specified for Profiles using the `Priority` traffic routing method. Supports
    values between 1 and 1000, with no Endpoints sharing the same value. If
    omitted the value will be computed in order of creation.

* `endpoint_location` - (Optional) Specifies the Azure location of the Endpoint,
  this must be specified for Profiles using the `Performance` routing method.

* `geo_mappings` - (Optional) A list of Geographic Regions used to distribute traffic, such as `WORLD`, `UK` or `DE`. The same location can't be specified in two endpoints. [See the Geographic Hierarchies documentation for more information](https://docs.microsoft.com/en-us/rest/api/trafficmanager/geographichierarchies/getdefault).

* `subnet` - (Optional) One or more `subnet` blocks as defined below

---

A `custom_header` block supports the following:

* `name` - (Required) The name of the custom header.

* `value` - (Required) The value of custom header. Applicable for Http and Https protocol.

---

A `subnet` block supports the following:

* `first` - (Required) The first IP Address in this subnet.

* `last` - (Optional) The last IP Address in this subnet.

* `scope` - (Optional) The block size (number of leading bits in the subnet mask).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Traffic Manager Nested Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Traffic Manager Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Traffic Manager Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Traffic Manager Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Traffic Manager Endpoint.

## Import

Traffic Manager Nested Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_traffic_manager_nested_endpoint.exampleEndpoints /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/trafficManagerProfiles/mytrafficmanagerprofile1/NestedEndpoints/mytrafficmanagerendpoint
```
