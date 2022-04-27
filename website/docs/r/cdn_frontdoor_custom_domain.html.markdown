---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_profile_custom_domain"
description: |-
  Manages a Frontdoor Custom Domain.
---

# azurerm_cdn_frontdoor_profile_custom_domain

Manages a Frontdoor Custom Domain.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "afdx-terraform.azfdtest.xyz "
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_custom_domain" "example" {
  name                     = "example-customDomain"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  dns_zone_id              = azurerm_dns_zone.example.id
  host_name                = "contoso.com"

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Custom Domain. Changing this forces a new Frontdoor Custom Domain to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the Frontdoor Profile. Changing this forces a new Frontdoor Profile to be created.

* `host_name` - (Required) The host name of the domain. Changing this forces a new Frontdoor Custom Domain to be created.

* `dns_zone_id` - (Optional) The Resource ID of the DNS Zone that is to be used for the Frontdoor Custom Domain.

* `pre_validated_cdn_frontdoor_custom_domain_id` - (Optional) Resource ID.

* `tls` - (Required) A `tls` block as defined below.

---

A `tls` block supports the following:

* `certificate_type` - (Optional) Defines the source of the SSL certificate. Possible values include `CustomerCertificate` and `ManagedCertificate`. Defaults to `ManagedCertificate`.

* `minimum_tls_version` - (Optional) TLS protocol version that will be used for Https. Possible values include `TLS10` and `TLS12`. Defaults to `TLS12`.

* `cdn_frontdoor_secret_id` - (Optional) Resource ID of the Frontdoor Secrect.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Custom Domain.

* `deployment_status` - 

* `domain_validation_state` - Provisioning substate shows the progress of custom HTTPS enabling/disabling process.

* `cdn_frontdoor_profile_name` - The name of the Frontdoor Profile which holds the domain.

* `provisioning_state` - Provisioning status

* `validation_properties` - A `validation_properties` block as defined below.

---

A `validation_properties` block exports the following:

* `expiration_date` - The date time that the token expires

* `validation_token` - Challenge used for DNS TXT record or file based validation

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 12 hours) Used when creating the Frontdoor Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Custom Domain.
* `update` - (Defaults to 24 hours) Used when updating the Frontdoor Custom Domain.
* `delete` - (Defaults to 12 hours) Used when deleting the Frontdoor Custom Domain.

## Import

Frontdoor Custom Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1
```
