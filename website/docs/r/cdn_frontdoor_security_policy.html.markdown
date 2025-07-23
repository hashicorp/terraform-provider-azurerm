---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_security_policy"
description: |-
  Manages a Front Door (standard/premium) Security Policy.
---

# azurerm_cdn_frontdoor_security_policy

Manages a Front Door (standard/premium) Security Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_firewall_policy" "example" {
  name                              = "exampleWAF"
  resource_group_name               = azurerm_resource_group.example.name
  sku_name                          = azurerm_cdn_frontdoor_profile.example.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

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
      match_values       = ["192.168.1.0/24", "10.0.1.0/24"]
    }
  }
}

resource "azurerm_dns_zone" "example" {
  name                = "sub-domain.domain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_custom_domain" "example" {
  name                     = "example-customDomain"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  dns_zone_id              = azurerm_dns_zone.example.id
  host_name                = "contoso.fabrikam.com"

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_cdn_frontdoor_security_policy" "example" {
  name                     = "Example-Security-Policy"
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Front Door Security Policy. Possible values must not be an empty string. Changing this forces a new Front Door Security Policy to be created.

* `cdn_frontdoor_profile_id` - (Required) The Front Door Profile Resource Id that is linked to this Front Door Security Policy. Changing this forces a new Front Door Security Policy to be created.

* `security_policies` - (Required) An `security_policies` block as defined below. Changing this forces a new Front Door Security Policy to be created.

---

A `security_policies` block supports the following:

* `firewall` - (Required) An `firewall` block as defined below. Changing this forces a new Front Door Security Policy to be created.

---

A `firewall` block supports the following:

* `cdn_frontdoor_firewall_policy_id` - (Required) The Resource Id of the Front Door Firewall Policy that should be linked to this Front Door Security Policy. Changing this forces a new Front Door Security Policy to be created.

* `association` - (Required) An `association` block as defined below. Changing this forces a new Front Door Security Policy to be created.

---

An `association` block supports the following:

* `domain` - (Required) One or more `domain` blocks as defined below. Changing this forces a new Front Door Security Policy to be created.

* `patterns_to_match` - (Required) The list of paths to match for this firewall policy. Possible value includes `/*`. Changing this forces a new Front Door Security Policy to be created.

---

A `domain` block supports the following:

~> **Note:** The number of `domain` blocks that maybe included in the configuration file varies depending on the `sku_name` field of the linked Front Door Profile. The `Standard_AzureFrontDoor` sku may contain up to 100 `domain` blocks and a `Premium_AzureFrontDoor` sku may contain up to 500 `domain` blocks.

* `cdn_frontdoor_domain_id` - (Required) The Resource Id of the **Front Door Custom Domain** or **Front Door Endpoint** that should be bound to this Front Door Security Policy. Changing this forces a new Front Door Security Policy to be created.

* `active` - (Computed) Is the Front Door Custom Domain/Endpoint activated?

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Security Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Security Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Security Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Security Policy.

## Import

Front Door Security Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_security_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/securityPolicies/policy1
```
