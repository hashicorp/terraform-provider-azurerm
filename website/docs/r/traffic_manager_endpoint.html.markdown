---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_endpoint"
description: |-
  Manages a Traffic Manager Endpoint.
---

# azurerm_traffic_manager_endpoint

Manages a Traffic Manager Endpoint.

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

resource "azurerm_traffic_manager_profile" "example" {
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

resource "azurerm_traffic_manager_endpoint" "example" {
  name                = random_id.server.hex
  resource_group_name = azurerm_resource_group.example.name
  profile_name        = azurerm_traffic_manager_profile.example.name
  target              = "terraform.io"
  type                = "externalEndpoints"
  weight              = 100
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Traffic Manager endpoint. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the Traffic Manager Profile exists.

* `profile_name` - (Required) The name of the Traffic Manager Profile to attach
    create the Traffic Manager endpoint.

* `endpoint_status` - (Optional) The status of the Endpoint, can be set to
    either `Enabled` or `Disabled`. Defaults to `Enabled`.

* `type` - (Required) The Endpoint type, must be one of:
    - `azureEndpoints`
    - `externalEndpoints`
    - `nestedEndpoints`

* `target` - (Optional) The FQDN DNS name of the target. This argument must be
    provided for an endpoint of type `externalEndpoints`, for other types it
    will be computed.

* `target_resource_id` - (Optional) The resource id of an Azure resource to
    target. This argument must be provided for an endpoint of type
    `azureEndpoints` or `nestedEndpoints`.

* `weight` - (Optional) Specifies how much traffic should be distributed to this
    endpoint, this must be specified for Profiles using the  `Weighted` traffic
    routing method. Supports values between 1 and 1000.

* `priority` - (Optional) Specifies the priority of this Endpoint, this must be
    specified for Profiles using the `Priority` traffic routing method. Supports
    values between 1 and 1000, with no Endpoints sharing the same value. If
    omitted the value will be computed in order of creation.

* `endpoint_location` - (Optional) Specifies the Azure location of the Endpoint,
    this must be specified for Profiles using the `Performance` routing method
    if the Endpoint is of either type `nestedEndpoints` or `externalEndpoints`.
    For Endpoints of type `azureEndpoints` the value will be taken from the
    location of the Azure target resource.

* `min_child_endpoints` - (Optional) This argument specifies the minimum number
    of endpoints that must be ‘online’ in the child profile in order for the
    parent profile to direct traffic to any of the endpoints in that child
    profile. This argument only applies to Endpoints of type `nestedEndpoints`
    and has to be larger than `0`.

~>**NOTE**: If `min_child_endpoints` is less than either `minimum_required_child_endpoints_ipv4` or `minimum_required_child_endpoints_ipv6`, then it won't have any effect.

* `minimum_required_child_endpoints_ipv4` - (Optional) This argument specifies the minimum number of IPv4 (DNS record type A) endpoints that must be ‘online’ in the child profile in order for the parent profile to direct traffic to any of the endpoints in that child profile. This argument only applies to Endpoints of type `nestedEndpoints` and defaults to `1`.

* `minimum_required_child_endpoints_ipv6` - (Optional) This argument specifies the minimum number of IPv6 (DNS record type AAAA) endpoints that must be ‘online’ in the child profile in order for the parent profile to direct traffic to any of the endpoints in that child profile. This argument only applies to Endpoints of type `nestedEndpoints` and defaults to `1`.

* `geo_mappings` - (Optional) A list of Geographic Regions used to distribute traffic, such as `WORLD`, `UK` or `DE`. The same location can't be specified in two endpoints. [See the Geographic Hierarchies documentation for more information](https://docs.microsoft.com/en-us/rest/api/trafficmanager/geographichierarchies/getdefault).

* `custom_header` - (Optional) One or more `custom_header` blocks as defined below

* `subnet` - (Optional) One or more `subnet` blocks as defined below

---
A `custom_header` block supports the following:

* `name` - (Required) The name of the custom header.

* `value` - (Required) The value of custom header. Applicable for Http and Https protocol.

A `subnet` block supports the following:

* `first` - (Required) The First IP....

* `last` - (Optional) The Last IP...

* `scope` - (Optional) The Scope...

-> **NOTE:** One and only one of either `last` (in case of IP range) or `scope` (in case of CIDR) must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Traffic Manager Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Traffic Manager Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Traffic Manager Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Traffic Manager Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Traffic Manager Endpoint.

## Import

Traffic Manager Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_traffic_manager_endpoint.exampleEndpoints /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/trafficManagerProfiles/mytrafficmanagerprofile1/azureEndpoints/mytrafficmanagerendpoint
```

-> **NOTE:** `azureEndpoints` in the above shell command should be replaced with `externalEndpoints` or `nestedEndpoints` while using other endpoint types.
