---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerum_cdn_frontdoor_custom_domain_route_association"
description: |-
  Manages a part of the workflow when creating a Frontdoor Custom Domain.
---

# azurerum_cdn_frontdoor_custom_domain_route_association

This resource is used with the `azurerm_cdn_frontdoor_custom_domain` and the `azurerm_cdn_frontdoor_route` resources to associate the Frontdoor Route with the Frontdoor Custom Domain. The successful creation of this resource represents that the Frontdoor Route has been successfully associated with the Frontdoor Custom Domain.

-> **NOTE:** The custom domain association functionality has been separated from the Frontdoor Route resource to accommodate the Custom Domain workflow logic.

## Example Usage

```hcl
resource "azurerm_cdn_frontdoor_custom_domain_route_association" "example" {
  cdn_frontdoor_route_id = azurerm_cdn_frontdoor_route.example.id

  cdn_frontdoor_custom_domain_txt_validator_ids = [azurerm_cdn_frontdoor_custom_domain_txt_validator.example.id]
  cdn_frontdoor_custom_domain_ids               = [azurerm_cdn_frontdoor_custom_domain.example.id]
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_route_id` - (Required) The resource ID of the Frontdoor Route to associate the Frontdoor Custom Domain(s) with. Changing this forces a new Frontdoor Custom Domain Route Association to be created.

* `cdn_frontdoor_custom_domain_txt_validator_ids` - (Required) One or more resource IDs of the Frontdoor Custom Domain Validator.

* `cdn_frontdoor_custom_domain_ids` - (Required) One or more resource IDs of the Frontdoor Custom Domain to associate with the Frontdoor Route.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Endpoint.

* `cdn_frontdoor_custom_domains_active_status` - One or more `cdn_frontdoor_custom_domains_active_status` blocks as defined below.

---
A `cdn_frontdoor_custom_domains_active_status` block exports the following:

* `id` - (Computed) The resource ID of the Frontdoor Custom Domain.

* `active` - (Computed) Is the Frontdoor Custom Domain active?

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Custom Domain Route Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Custom Domain Route Association.
* `update` - (Defaults to 30 minutes) Used when retrieving the Frontdoor Custom Domain Route Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Custom Domain Route Association.
