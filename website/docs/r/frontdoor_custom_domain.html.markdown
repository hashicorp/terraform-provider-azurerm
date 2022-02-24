---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_profile_custom_domain"
description: |-
  Manages a Frontdoor Custom Domain.
---

# azurerm_frontdoor_profile_custom_domain

Manages a Frontdoor Custom Domain.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-cdn"
  location = "West Europe"
}

resource "azurerm_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_frontdoor_profile_custom_domain" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
  dns_zone_id          = ""
  host_name            = ""

  pre_validated_custom_domain_resource_id {
    id = ""
  }

  tls_settings {
    certificate_type    = ""
    minimum_tls_version = ""
    secret_id           = ""
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Custom Domain. Changing this forces a new Frontdoor Custom Domain to be created.

* `frontdoor_profile_id` - (Required) The ID of the Frontdoor Profile. Changing this forces a new Frontdoor Profile to be created.

* `host_name` - (Required) The host name of the domain. Must be a domain name. Changing this forces a new Frontdoor Custom Domain to be created.

* `dns_zone_id` - (Optional) Resource ID.

* `pre_validated_custom_domain_resource_id` - (Optional) A `pre_validated_custom_domain_resource_id` block as defined below.

* `tls_settings` - (Optional) A `tls_settings` block as defined below.

---

A `pre_validated_custom_domain_resource_id` block supports the following:

* `id` - (Optional) Resource ID.

---

A `tls_settings` block supports the following:

* `certificate_type` - (Required) Defines the source of the SSL certificate.

* `minimum_tls_version` - (Optional) TLS protocol version that will be used for Https

* `secret_id` - (Optional) Resource ID.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Custom Domain.

* `deployment_status` - 

* `domain_validation_state` - Provisioning substate shows the progress of custom HTTPS enabling/disabling process step by step. DCV stands for DomainControlValidation.

* `profile_name` - The name of the profile which holds the domain.

* `provisioning_state` - Provisioning status

* `validation_properties` - A `validation_properties` block as defined below.

---

A `validation_properties` block exports the following:

* `expiration_date` - The date time that the token expires

* `validation_token` - Challenge used for DNS TXT record or file based validation

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Custom Domain.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Custom Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Custom Domain.

## Import

Frontdoor Custom Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor_profile_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1
```
