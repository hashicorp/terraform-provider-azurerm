---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_custom_domain_association"
description: |-
  Manages the association between a Front Door (standard/premium) Custom Domain and one or more Front Door (standard/premium) Routes.
---

# azurerm_cdn_frontdoor_custom_domain_association

Manages the association between a Front Door (standard/premium) Custom Domain and one or more Front Door (standard/premium) Routes.

## Example Usage

```hcl
resource "azurerm_cdn_frontdoor_custom_domain_association" "example" {
  cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.contoso.id
  cdn_frontdoor_route_ids        = [azurerm_cdn_frontdoor_route.contoso.id, azurerm_cdn_frontdoor_route.fabrikam.id]
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_custom_domain_id` - (Required) The ID of the Front Door Custom Domain that should be managed by the association resource. Changing this forces a new association resource to be created.

* `cdn_frontdoor_route_ids` - (Required) One or more IDs of the Front Door Route to which the Front Door Custom Domain is associated with.

-> **NOTE:** This should include all of the Front Door Route resources that the Front Door Custom Domain is associated with. If the list of Front Door Routes is not complete you will receive the service side error `This resource is still associated with a route. Please delete the association with the route first before deleting this resource` when you attempt to `destroy`/`delete` your Front Door Custom Domain.

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
terraform import azurerm_cdn_frontdoor_custom_domain_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/associations/assoc1
```
