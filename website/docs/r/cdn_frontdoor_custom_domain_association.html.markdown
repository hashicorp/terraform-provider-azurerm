---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_custom_domain_association"
description: |-
  Manages the Custom Domain associations between one or more Front Door (standard/premium) Custom Domains and a Front Door (standard/premium) Route.
---

# azurerm_cdn_frontdoor_custom_domain_association

Manages the Custom Domain associations between one or more Front Door (standard/premium) Custom Domains and a Front Door (standard/premium) Route.

## Example Usage

```hcl
resource "azurerm_cdn_frontdoor_custom_domain_association" "example" {
  cdn_frontdoor_route_id          = azurerm_cdn_frontdoor_route.contoso.id
  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.contoso.id]
  link_to_default_domain          = false
}
```

## Example Usage using the `azurerm_cdn_frontdoor_route` Data Source

```hcl
data "azurerm_cdn_frontdoor_route" "example" {
  name                      = azurerm_cdn_frontdoor_route.example.name
  cdn_frontdoor_endpoint_id = azurerm_cdn_frontdoor_endpoint.example.id
}

resource "azurerm_cdn_frontdoor_custom_domain_association" "example" {
  cdn_frontdoor_route_id          = data.azurerm_cdn_frontdoor_route.example.id
  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.example.id]
  link_to_default_domain          = false
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_route_id` - (Required) The Resource ID of the Front Door Route which this Front Door Custom Domain Association resource will manage. Changing this forces a new Front Door Custom Domain Association resource to be created.

* `cdn_frontdoor_custom_domain_ids` - (Required) The Resource IDs of the Front Door Custom Domain(s) which are to be associated with the Front Door Route.

* `link_to_default_domain` - (Optional) Should this Front Door Route be linked to the default endpoint? Possible values include `true` or `false`. Defaults to `true`.

-> **NOTE:** The `link_to_default_domain` value will be automatically toggled to `true` on deletion of this resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Custom Domain Association.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Custom Domain Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Custom Domain Association.
* `update` - (Defaults to 30 minutes) Used when retrieving the Front Door Custom Domain Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Custom Domain Association.

## Import

Front Door Custom Domain Associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_custom_domain_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/associations/route1
```
