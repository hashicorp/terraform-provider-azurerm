---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_inbound_endpoint"
description: |-
  Gets information about an existing Private DNS Resolver Inbound Endpoint.
---

# Data Source: azurerm_private_dns_resolver_inbound_endpoint

Gets information about an existing Private DNS Resolver Inbound Endpoint.

## Example Usage

```hcl
data "azurerm_private_dns_resolver_inbound_endpoint" "example" {
  name                    = "example-drie"
  private_dns_resolver_id = "example-private-dns-resolver-id"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the Private DNS Resolver Inbound Endpoint.

* `private_dns_resolver_id` - (Required) ID of the Private DNS Resolver.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Inbound Endpoint.

* `ip_configurations` - A list of `ip_configurations` block as defined below.

* `location` - The Azure Region where the Private DNS Resolver Inbound Endpoint exists.

* `tags` - The tags assigned to the Private DNS Resolver Inbound Endpoint.

---

An `ip_configurations` block exports the following:

* `private_ip_allocation_method` - The private IP address allocation method.

* `subnet_id` - The subnet ID of the IP configuration.

* `private_ip_address` - The private IP address of the IP configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Inbound Endpoint.
