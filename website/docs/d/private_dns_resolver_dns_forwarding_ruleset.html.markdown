---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_dns_forwarding_ruleset"
description: |-
  Gets information about an existing Private DNS Resolver Dns Forwarding Ruleset.
---

# Data Source: azurerm_private_dns_resolver_dns_forwarding_ruleset

Gets information about an existing Private DNS Resolver Dns Forwarding Ruleset.

## Example Usage

```hcl
data "azurerm_private_dns_resolver_dns_forwarding_ruleset" "example" {
  name                = "example-ruleset"
  resource_group_name = "example-ruleset-resourcegroup"
}
```

## Arguments Reference

The following arguments are required:

* `name` - Name of the existing Private DNS Resolver Dns Forwarding Ruleset.

* `resource_group_name` - Name of the Resource Group where the Private DNS Resolver Dns Forwarding Ruleset exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Dns Forwarding Ruleset.

* `private_dns_resolver_outbound_endpoint_ids` - List of IDs of the Private DNS Resolver Outbound Endpoint that is linked to the Private DNS Resolver Dns Forwarding Ruleset.

* `location` - Azure Region where the Private DNS Resolver Dns Forwarding Ruleset exists.

* `tags` - Mapping of tags assigned to the Private DNS Resolver Dns Forwarding Ruleset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS SRV Record.
