---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_outbound_endpoint"
description: |-
  Gets information about an existing Private DNS Resolver Outbound Endpoint.
---

# Data Source: azurerm_private_dns_resolver_outbound_endpoint

Gets information about an existing Private DNS Resolver Outbound Endpoint.

## Example Usage

```hcl
data "azurerm_private_dns_resolver_outbound_endpoint" "example" {
  name                    = "example-endpoint"
  private_dns_resolver_id = "example-private-dns-resolver-id"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the Private DNS Resolver Outbound Endpoint.

* `private_dns_resolver_id` - (Required) ID of the Private DNS Resolver Outbound Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Outbound Endpoint.

* `location` - The Azure Region where the Private DNS Resolver Outbound Endpoint exists.

* `subnet_id` - The ID of the Subnet that is linked to the Private DNS Resolver Outbound Endpoint.

* `tags` - The tags assigned to the Private DNS Resolver Outbound Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Outbound Endpoint.
