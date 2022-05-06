---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerum_cdn_frontdoor_custom_domain_secret_validator"
description: |-
  Manages a part of the workflow when creating a Frontdoor Custom Domain.
---

# azurerum_cdn_frontdoor_custom_domain_secret_validator

The successful creation of this resource represents that the TLS certificate(s) have been successfully provisioned to each Frontdoor Custom Domain by the Frontdoor service.

This resource is used with the `azurerm_cdn_frontdoor_custom_domain` and the `cdn_frontdoor_route` resources to verify that the certificate defined in the `azurerm_cdn_frontdoor_custom_domain` `tls_settings` field has been successfully provisioned or not. The resource will halt the provisioning of other CDN Frontdoor resources while it waits for the Frontdoor service to provision the TLS certificate to there custom domain(s) before proceeding.

~> **WARNING:** This resource implements a part of the validation workflow logic for Frontdoor Custom Domains. It does not represent a real-world entity in Azure, therefore changing or deleting this resource on its own has no immediate effect on the Azure CDN Frontdoor resource itself. If you are not implementing Frontdoor Custom Domains you do not need to add this resource to your configuration file.

## Example Usage

```hcl
resource "azurerm_cdn_frontdoor_custom_domain_secret_validator" "example" {
  cdn_frontdoor_route_id                        = azurerm_cdn_frontdoor_route.example.id
  cdn_frontdoor_custom_domain_ids               = [azurerm_cdn_frontdoor_custom_domain.contoso.id, azurerm_cdn_frontdoor_custom_domain.fabrikam.id]
  cdn_frontdoor_custom_domain_txt_validator_ids = [azurerm_cdn_frontdoor_custom_domain_txt_validator.contoso.id, azurerm_cdn_frontdoor_custom_domain_txt_validator.fabrikam.id]
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_route_id` - (Required) The resource ID of the Frontdoor Route which holds the Custom Domain association information. Changing this forces a new Frontdoor Custom Domain Secret Validator to be created.

* `cdn_frontdoor_custom_domain_ids` - (Required) One or more resource IDs of the Frontdoor Custom Domain(s) to validate. Changing this forces a new Frontdoor Custom Domain Secret Validator to be created.

* `cdn_frontdoor_custom_domain_txt_validator_ids` (Required) One or more resource IDs of the Frontdoor Custom Domain TXT validators that will trigger the Frontdoor Custom Domain Secret Validator to begin polling. Changing this forces a new Frontdoor Custom Domain Secret Validator to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Custom Domain Secret Validator.

* `cdn_frontdoor_custom_domain_secrets_state` - One or more `cdn_frontdoor_custom_domain_secrets_state` blocks as defined below.

---
A `cdn_frontdoor_custom_domains_active_status` block exports the following:

* `cdn_frontdoor_custom_domain_id` - (Computed) The resource ID of the CDN frontdoor custom domain.

* `cdn_frontdoor_secret_id` - (Computed) The resource ID of the TLS secret.

* `cdn_frontdoor_secret_provisioning_state` - (Computed) The provisioning state of the TLS certificate.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Frontdoor Custom Domain Secret Validator.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Custom Domain Secret Validator.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Custom Domain Secret Validator.
