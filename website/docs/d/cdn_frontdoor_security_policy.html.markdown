---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_cdn_frontdoor_security_policy"
description: |-
  Gets information about an existing Front Door (standard/premium) Security Policy.
---

# Data Source: azurerm_cdn_frontdoor_security_policy

Gets information about an existing Front Door (standard/premium) Security Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-frontdoor-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_firewall_policy" "example" {
  name                = "examplecdnfrontdoorfirewallpolicy"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = azurerm_cdn_frontdoor_profile.example.sku_name
  enabled             = true
  mode                = "Prevention"
  redirect_url        = "https://www.example.com"

  custom_rule {
    name                           = "Rule1"
    enabled                        = true
    priority                       = 1
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }
  }
}

resource "azurerm_dns_zone" "example" {
  name                = "example-frontdoor.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_custom_domain" "example" {
  name                     = "example-custom-domain"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  dns_zone_id              = azurerm_dns_zone.example.id
  host_name                = "www.example-frontdoor.com"

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_cdn_frontdoor_security_policy" "example" {
  name                     = "example-security-policy"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  security_policies {
    firewall {
      cdn_frontdoor_firewall_policy_id = azurerm_cdn_frontdoor_firewall_policy.example.id

      association {
        domain {
          cdn_frontdoor_domain_id = azurerm_cdn_frontdoor_custom_domain.example.id
        }

        patterns_to_match = ["/*"]
      }
    }
  }
}

data "azurerm_cdn_frontdoor_security_policy" "example" {
  name                = azurerm_cdn_frontdoor_security_policy.example.name
  profile_name        = azurerm_cdn_frontdoor_profile.example.name
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Front Door Security Policy.

* `profile_name` - (Required) The name of the Front Door Profile.

* `resource_group_name` - (Required) The name of the Resource Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Security Policy.

* `cdn_frontdoor_profile_id` - The ID of the Front Door Profile associated with this Front Door Security Policy.

* `security_policies` - A `security_policies` block as defined below.

---

A `security_policies` block exports the following:

* `firewall` - A `firewall` block as defined below.

---

A `firewall` block exports the following:

* `association` - An `association` block as defined below.

* `cdn_frontdoor_firewall_policy_id` - The ID of the Front Door Firewall Policy associated with this Front Door Security Policy.

---

An `association` block exports the following:

* `domain` - A `domain` block as defined below.

* `patterns_to_match` - The paths associated with this firewall policy.

---

A `domain` block exports the following:

* `active` - Is the Front Door Custom Domain or Front Door Endpoint active?

* `cdn_frontdoor_domain_id` - The ID of the Front Door Custom Domain or Front Door Endpoint associated with this Front Door Security Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Security Policy.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Cdn` - 2024-02-01
