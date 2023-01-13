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

The following arguments are supported:

* `name` - (Required) Name of the existing Private DNS Resolver Dns Forwarding Ruleset.

* `resource_group_name` - (Required) Name of the Resource Group where the Private DNS Resolver Dns Forwarding Ruleset exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Dns Forwarding Ruleset.

* `private_dns_resolver_outbound_endpoint_ids` - The IDs list of the Private DNS Resolver Outbound Endpoints that are linked to the Private DNS Resolver Dns Forwarding Ruleset.

* `location` - The Azure Region where the Private DNS Resolver Dns Forwarding Ruleset exists.

* `tags` - The tags assigned to the Private DNS Resolver Dns Forwarding Ruleset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Dns Forwarding Ruleset.
