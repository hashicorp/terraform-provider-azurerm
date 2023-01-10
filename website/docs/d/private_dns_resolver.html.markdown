---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver"
description: |-
  Gets information about an existing Private DNS Resolver.
---

# Data Source: azurerm_private_dns_resolver

Gets information about an existing Private DNS Resolver.

## Example Usage

```hcl
data "azurerm_private_dns_resolver" "test" {
  name                = "example"
  resource_group_name = "example-resourcegroup-name"
}
```

## Arguments Reference

The following arguments are required:

* `name` - Name of the Private DNS Resolver.

* `resource_group_name` - Name of the Resource Group where the Private DNS Resolver exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - ID of the DNS Resolver.

* `location` - Azure Region where the Private DNS Resolver exists.

* `virtual_network_id` - ID of the Virtual Network that is linked to the Private DNS Resolver.

* `tags` - Mapping of tags which should be assigned to the Private DNS Resolver.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS SRV Record.
