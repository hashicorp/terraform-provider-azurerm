---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_azure_endpoint"
description: |-
   Manages an Azure Endpoint within a Traffic Manager Profile..
---

# azurerm_traffic_manager_azure_endpoint

Manages an Azure Endpoint within a Traffic Manager Profile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  domain_name_label   = "example-public-ip"
}


resource "azurerm_traffic_manager_profile" "example" {
  name                   = "example-profile"
  resource_group_name    = azurerm_resource_group.example.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "example-profile"
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

resource "azurerm_traffic_manager_azure_endpoint" "example" {
  name                 = "example-endpoint"
  profile_id           = azurerm_traffic_manager_profile.example.id
  always_serve_enabled = true
  weight               = 100
  target_resource_id   = azurerm_public_ip.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Endpoint. Changing this forces a new resource to be created.

* `profile_id` - (Required) The ID of the Traffic Manager Profile that this Azure Endpoint should be created within. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the Azure Resource which should be used as a target.

* `weight` - (Optional) Specifies how much traffic should be distributed to this endpoint, this must be specified for Profiles using the Weighted traffic routing method. Valid values are between `1` and `1000`. Defaults to `1`.

* `always_serve_enabled` - (Optional) If Always Serve is enabled, probing for endpoint health will be disabled and endpoints will be included in the traffic routing method. Defaults to `false`.

* `custom_header` - (Optional) One or more `custom_header` blocks as defined below.

* `enabled` - (Optional) Is the endpoint enabled? Defaults to `true`.

* `geo_mappings` - (Optional) A list of Geographic Regions used to distribute traffic, such as `WORLD`, `UK` or `DE`. The same location can't be specified in two endpoints. [See the Geographic Hierarchies documentation for more information](https://docs.microsoft.com/rest/api/trafficmanager/geographichierarchies/getdefault).

* `priority` - (Optional) Specifies the priority of this Endpoint, this must be specified for Profiles using the `Priority` traffic routing method. Supports values between 1 and 1000, with no Endpoints sharing the same value. If omitted the value will be computed in order of creation.

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

* `id` - The ID of the Azure Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Endpoint.

## Import

Azure Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_traffic_manager_azure_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Network/trafficManagerProfiles/example-profile/AzureEndpoints/example-endpoint
```
