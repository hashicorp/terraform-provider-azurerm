---
subcategory: "Cdn"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_endpoint"
description: |-
  Manages a Frontdoor Endpoint.
---

# azurerm_frontdoor_endpoint

Manages a Frontdoor Endpoint.

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

resource "azurerm_frontdoor_endpoint" "test" {
  name                            = "acctest-c-%d"
  frontdoor_profile_id            = azurerm_frontdoor_profile.test.id
  enabled_state                   = ""
  location                        = "%s"
  origin_response_timeout_seconds = 0

  tags = {
    ENV = "Test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Endpoint. Changing this forces a new Frontdoor Endpoint to be created.

* `frontdoor_profile_id` - (Required) The ID of the Frontdoor Profile. Changing this forces a new Frontdoor Endpoint to be created.

* `location` - (Required) The Azure Region where the Frontdoor Endpoint should exist. Changing this forces a new Frontdoor Endpoint to be created.

* `enabled_state` - (Optional) Whether to enable use of this rule. Permitted values are 'Enabled' or 'Disabled'

* `origin_response_timeout_seconds` - (Optional) Send and receive timeout on forwarding request to the origin. When timeout is reached, the request fails and returns.

* `tags` - (Optional) A mapping of tags which should be assigned to the Frontdoor Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Endpoint.

* `deployment_status` - 

* `host_name` - The host name of the endpoint structured as {endpointName}.{DNSZone}, e.g. contoso.azureedge.net

* `profile_name` - The name of the profile which holds the endpoint.

* `provisioning_state` - Provisioning status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Endpoint.

## Import

Frontdoor Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1
```
