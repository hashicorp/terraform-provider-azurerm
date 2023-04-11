---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_forwarding_rule"
description: |-
  Gets information about an existing Private DNS Resolver Forwarding Rule.
---

# Data Source: azurerm_private_dns_resolver_forwarding_rule

Gets information about an existing Private DNS Resolver Forwarding Rule.

## Example Usage

```hcl
data "azurerm_private_dns_resolver_forwarding_rule" "example" {
  name                      = "example-rule"
  dns_forwarding_ruleset_id = "example-forwarding-rulset-id"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the Private DNS Resolver Forwarding Rule.

* `dns_forwarding_ruleset_id` - (Required) ID of the Private DNS Resolver Forwarding Ruleset.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Forwarding Rule.

* `domain_name` - The domain name for the Private DNS Resolver Forwarding Rule.

* `target_dns_servers` - A list of `target_dns_servers` block as defined below.

* `enabled` - Is the Private DNS Resolver Forwarding Rule enabled?

* `metadata` - The metadata attached to the Private DNS Resolver Forwarding Rule.

---

A `target_dns_servers` block exports the following:

* `ip_address` - The DNS server IP address.

* `port` - The DNS server port.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Forwarding Rule.
