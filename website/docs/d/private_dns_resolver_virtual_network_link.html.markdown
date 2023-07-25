---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_virtual_network_link"
description: |-
  Gets information about an existing Private DNS Resolver Virtual Network Link.
---

# Data Source: azurerm_private_dns_resolver_virtual_network_link

Gets information about an existing Private DNS Resolver Virtual Network Link.

## Example Usage

```hcl
data "azurerm_private_dns_resolver_virtual_network_link" "example" {
  name                      = "example-link"
  dns_forwarding_ruleset_id = "example-dns-forwarding-ruleset-id"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the Private DNS Resolver Virtual Network Link.

* `dns_forwarding_ruleset_id` - (Required) ID of the Private DNS Resolver DNS Forwarding Ruleset.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Virtual Network Link.

* `virtual_network_id` - The ID of the Virtual Network that is linked to the Private DNS Resolver Virtual Network Link.

* `metadata` - The metadata attached to the Private DNS Resolver Virtual Network Link.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Virtual Network Link.
