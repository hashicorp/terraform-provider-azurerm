---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_endpoint"
description: |-
  Manages a Front Door (standard/premium) Endpoint.
---

# azurerm_cdn_frontdoor_endpoint

Manages a Front Door (standard/premium) Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                     = "example-endpoint"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  tags = {
    ENV = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Front Door Endpoint. Changing this forces a new Front Door Endpoint to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the Front Door Profile within which this Front Door Endpoint should exist. Changing this forces a new Front Door Endpoint to be created.

---

* `enabled` - (Optional) Specifies if this Front Door Endpoint is enabled? Defaults to `true`.

* `tags` - (Optional) Specifies a mapping of tags which should be assigned to the Front Door Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this Front Door Endpoint.

* `host_name` - The host name of the Front Door Endpoint, in the format `{endpointName}.{dnsZone}` (for example, `contoso.azureedge.net`).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Front Door Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Endpoint.

## Import

Front Door Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1
```
