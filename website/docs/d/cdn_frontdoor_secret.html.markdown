---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_secret"
description: |-
  Gets information about an existing Front Door (standard/premium) Secret.
---

# Data Source: azurerm_cdn_frontdoor_secret

Use this data source to access information about an existing Front Door (standard/premium) Secret.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_secret" "example" {
  name                = "example-secret"
  profile_name        = "example-profile"
  resource_group_name = "example-resources"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door Secret.

* `profile_name` - (Required) The name of the Front Door Profile within which the Front Door Secret exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Front Door Profile exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Secret.

* `cdn_frontdoor_profile_id` - Specifies the ID of the Front Door Profile within which this Front Door Secret exists.

* `secret` - A `secret` block as defined below.

---

A `secret` block exports the following:

* `customer_certificate` - A `customer_certificate` block as defined below.

---

A `customer_certificate` block exports the following:


* `expiration_date` - The key vault certificate expiration date.
 
* `key_vault_certificate_id` - The key vault certificate ID.

* `subject_alternative_names` - One or more `subject alternative names` contained within the key vault certificate.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Secret.
