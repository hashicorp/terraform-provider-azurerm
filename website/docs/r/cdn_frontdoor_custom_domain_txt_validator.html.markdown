---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerum_cdn_frontdoor_custom_domain_txt_validator"
description: |-
  Manages a part of the workflow when creating a Frontdoor Custom Domain.
---

# azurerum_cdn_frontdoor_custom_domain_txt_validator

The successful creation of this resource represents that the DNS TXT record for a Frontdoor Custom Domain has been validated by the Frontdoor service.

This resource is used with the `azurerm_cdn_frontdoor_custom_domain` resource to verify domain ownership via a DNS TXT record. The resource will halt the provisioning of other CDN Frontdoor resources while it waits for the Frontdoor service to validate the DNS TXT record in the custom domain before proceeding.

~> **WARNING:** This resource implements a part of the validation workflow logic for CDN Frontdoor Custom Domains. It does not represent a real-world entity in Azure, therefore changing or deleting this resource on its own has no immediate effect on the Azure CDN Frontdoor resource itself. If you are not implementing Frontdoor Custom Domains you do not need to add this resource to your configuration file.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "westeurope"
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Premium_AzureFrontDoor"

  response_timeout_seconds = 120

  tags = {
    environment = "Test"
  }
}

resource "azurerm_cdn_frontdoor_custom_domain" "contoso" {
  name                     = "contoso-custom-domain"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  dns_zone_id              = azurerm_dns_zone.example.id
  host_name                = join(".", ["contoso", azurerm_dns_zone.example.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_dns_txt_record" "contoso" {
  name                = join(".", ["_dnsauth", "contoso"])
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 3600

  record {
    value = azurerm_cdn_frontdoor_custom_domain.contoso.validation_token
  }
}

resource "azurerm_cdn_frontdoor_custom_domain_txt_validator" "example" {
  cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_custom_domain_id` - (Required) The resource ID of the Frontdoor Custom Domain to validate. Changing this forces a new Frontdoor Custom Domain Txt Validator to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor custom domain txt validator.

* `cdn_frontdoor_custom_domain_validation_state` - The state of the Frontdoor Custom Domain DNS TXT record validation process. Possible return values include `Approved`, `InternalError`, `Pending`, `PendingRevalidation`, `RefreshingValidationToken`, `Rejected`, `Submitting`, `TimedOut` and `Unknown`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Frontdoor Custom Domain Txt Validator.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Custom Domain Txt Validator.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Custom Domain Txt Validator.
