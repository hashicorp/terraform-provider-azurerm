---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_custom_domain"
description: |-
  Gets information about an existing Front Door (standard/premium) Custom Domain.
---

# Data Source: azurerm_cdn_frontdoor_custom_domain

Use this data source to access information about an existing Front Door (standard/premium) Custom Domain.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_custom_domain" "example" {
  name                = azurerm_cdn_frontdoor_custom_domain.example.name
  profile_name        = azurerm_cdn_frontdoor_profile.example.name
  resource_group_name = azurerm_cdn_frontdoor_profile.example.resource_group_name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Front Door Custom Domain.

* `profile_name` - (Required) The name of the Front Door Profile which the Front Door Custom Domain is bound to.

* `resource_group_name` - (Required) The name of the Resource Group where the Front Door Profile exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Custom Domain.

* `cdn_frontdoor_profile_id` - The ID of the Front Door Profile which the Front Door Custom Domain is bound to.

* `expiration_date` - The date time that the token expires.

* `host_name` - The host name of the domain.

* `tls` - A `tls` block as defined below.

* `validation_token` - The challenge used for DNS TXT record or file based validation.

---

A `tls` block exports the following:

* `cdn_frontdoor_secret_id` - The Resource ID of the Front Door Secret.

* `certificate_type` - The SSL certificate type.

* `minimum_tls_version` - The TLS protocol version that will be used for Https connections.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Custom Domain.
