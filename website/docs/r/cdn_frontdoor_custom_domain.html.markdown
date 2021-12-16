---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_custom_domain"
description: |-
  Manages an Azure Front Door (Standard/Premium) instance. (currently in public preview)
---

# azurerm_cdn_frontdoor_custom_domain

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "afdpremv2"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Premium_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_custom_domain" "example" {
  name       = "mycustomdomain"
  profile_id = azurerm_cdn_frontdoor_profile.example.id
  host_name  = "mycustomdomain.com"

  tls {
    certificate_type = "ManagedCertificate"
  }
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Custom Domain. I.e. `mydomain-com`.

* `profile_id` - (Required) Resource ID of the Azure Front Door Profile.

* `host_name` - (Required) Hostname i.e. `mydomain.com` 

---

The `tls` block supports the following:

* `certificate_type` can be either `CustomerCertificate` or `ManagedCertificate`.

* `minimum_tls_version` can be either `TLS10` or `TLS12`. 

* `secret_id` refers to the certificate used. Can only be used combination with `certificate_type` set to `CustomerCertificate`.

---

The following attributes are exported:

* `validation_token` - The validation token used for domain validation via a TXT DNS record.
