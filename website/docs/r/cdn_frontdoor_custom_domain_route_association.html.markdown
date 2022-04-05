---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerum_cdn_frontdoor_custom_domain_route_association"
description: |-
  Manages a part of the workflow when creating a Frontdoor Custom Domain.
---

# azurerum_cdn_frontdoor_custom_domain_route_association

This resource is used with the `azurerm_cdn_frontdoor_custom_domain` and the `azurerm_cdn_frontdoor_route` to associate the Frontdoor Route with the Frontdoor Custom Domain. The successful creation of this resource represents that the Frontdoor Route has been successfully associated with the Frontdoor Custom Domain.

-> **NOTE:** The custom domain association functionality has been separated from the Frontdoor Route resource to accommodate for the Custom Domain workflow logic.

## Example Usage

```hcl
resource "azurerm_cdn_frontdoor_custom_domain_route_association" "example" {
  cdn_frontdoor_route_id                       = azurerm_cdn_frontdoor_route.example.id
  cdn_frontdoor_custom_domain_txt_validator_id = azurerm_cdn_frontdoor_custom_domain_txt_validator.example.id

  custom_domains {
    id = azurerm_cdn_frontdoor_custom_domain.example.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_route_id` - (Required) The resource ID of the Frontdoor Custom Domain to validate. Changing this forces a new Frontdoor Custom Domain Route Association to be created.

* `cdn_frontdoor_custom_domain_txt_validator_id` - (Required) The resource ID of the Frontdoor Custom Domain to validate. Changing this forces a new Frontdoor Custom Domain Route Association to be created.

* `custom_domains` - (Required) One or more `custom_domains` block as defined below.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Endpoint.

[comment]: <> (TODO: This needs to be a list map structure like custom domains)
* `cdn_frontdoor_custom_domain_validation_state` - The state of the Frontdoor Custom Domain TXT record validation process. Possible return values include `Approved`, `InternalError`, `Pending`, `PendingRevalidation`, `RefreshingValidationToken`, `Rejected`, `Submitting`, `TimedOut` and `Unknown`.

---

A `custom_domains` block supports the following:

* `id` - (Optional) Resource ID.

* `active` - (Computed) Is the custom domain active?

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Custom Domain Route Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Custom Domain Route Association.
* `update` - (Defaults to 30 minutes) Used when retrieving the Frontdoor Custom Domain Route Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Custom Domain Route Association.