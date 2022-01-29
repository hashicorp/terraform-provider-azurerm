---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_azure_endpoint"
description: |-
  Manages a Traffic Manager Azure Endpoint.
---

# azurerm_traffic_manager_azure_endpoint

Manages a Traffic Manager Azure Endpoint.

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

resource "azurerm_traffic_manager_azure_endpoint" "example" {
  name                = random_id.server.hex
  resource_group_name = azurerm_resource_group.example.name
  profile_name        = azurerm_traffic_manager_profile.example.name
  weight              = 100
  target_resource_id  = azurerm_public_ip.test.id
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

* `custom_header` - (Optional) One or more `custom_header` blocks as defined below.

* `priority` - (Optional) Specifies the priority of this Endpoint, this must be
  specified for Profiles using the `Priority` traffic routing method. Supports
  values between 1 and 1000, with no Endpoints sharing the same value. If
  omitted the value will be computed in order of creation.

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

* `id` - The ID of the Traffic Manager Azure Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Traffic Manager Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Traffic Manager Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Traffic Manager Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Traffic Manager Endpoint.

## Import

Traffic Manager Azure Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_traffic_manager_azure_endpoint.exampleEndpoints /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/trafficManagerProfiles/mytrafficmanagerprofile1/AzureEndpoints/mytrafficmanagerendpoint
```
