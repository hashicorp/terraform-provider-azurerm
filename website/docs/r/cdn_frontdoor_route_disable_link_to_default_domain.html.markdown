---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_route_disable_link_to_default_domain"
description: |-
  Manages the Link To Default Domain property of a Front Door (standard/premium) Route.
---

# azurerm_cdn_frontdoor_route_disable_link_to_default_domain

Manages the Link To Default Domain property of a Front Door (standard/premium) Route.

!>**IMPORTANT:** This resource has been deprecated and should not be used for new deployments. The `azurerm_cdn_frontdoor_route_disable_link_to_default_domain` resource will be removed from the 4.0 AzureRM provider. Please use the `link_to_default_domain` field in the `azurerm_cdn_frontdoor_route` resource to control this value.

## Example Usage

```hcl
resource "azurerm_cdn_frontdoor_route_disable_link_to_default_domain" "example" {
  cdn_frontdoor_route_id          = azurerm_cdn_frontdoor_route.example.id
  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.contoso.id, azurerm_cdn_frontdoor_custom_domain.fabrikam.id]
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_route_id` - (Required) The resource ID of the Front Door Route where the Link To Default Domain property should be `disabled`. Changing this forces a new Front Door Route Disable Link To Default Domain to be created.

* `cdn_frontdoor_custom_domain_ids` - (Required) The resource IDs of the Front Door Custom Domains which are associated with this Front Door Route.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Route Disable Link To Default Domain.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Route Disable Link To Default Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Route Disable Link To Default Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Route Disable Link To Default Domain.
* `update` - (Defaults to 30 minutes) Used when updating the Cdn Frontdoor Route Disable Link To Default Domain.

## Import

FrontDoor Route Disable Link To Default Domain can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_route_disable_link_to_default_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/routes/route1/disableLinkToDefaultDomain/disableLinkToDefaultDomain1
```
