---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_outbound_endpoint"
description: |-
  Manages a Private DNS Resolver Outbound Endpoint.
---

# azurerm_private_dns_resolver_outbound_endpoint

Manages a Private DNS Resolver Outbound Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_private_dns_resolver_private_dns_resolver" "example" {
  name                = "example-drdr"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_resolver_outbound_endpoint" "example" {
  name                                         = "example-droe"
  private_dns_resolver_private_dns_resolver_id = azurerm_private_dns_resolver_private_dns_resolver.test.id
  location                                     = "West Europe"
  subnet {
    id = ""
  }
  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Private DNS Resolver Outbound Endpoint. Changing this forces a new Private DNS Resolver Outbound Endpoint to be created.

* `private_dns_resolver_private_dns_resolver_id` - (Required) Specifies the ID of the Private DNS Resolver Outbound Endpoint. Changing this forces a new Private DNS Resolver Outbound Endpoint to be created.

* `location` - (Required) Specifies the Azure Region where the Private DNS Resolver Outbound Endpoint should exist. Changing this forces a new Private DNS Resolver Outbound Endpoint to be created.

* `subnet_id` - (Required)  The ID of the Subnet that is linked to the Private DNS Resolver Outbound Endpoint.

* `tags` - (Optional) A mapping of tags which should be assigned to the Private DNS Resolver Outbound Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Outbound Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS Resolver Outbound Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Outbound Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS Resolver Outbound Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS Resolver Outbound Endpoint.

## Import

Private DNS Resolver Outbound Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_resolver_outbound_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsResolvers/dnsResolver1/outboundEndpoints/outboundEndpoint1
```
