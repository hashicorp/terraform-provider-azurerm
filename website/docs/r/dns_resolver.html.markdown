---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_resolver"
description: |-
  Manages a DNS Resolver.
---

# azurerm_dns_resolver

Manages a DNS Resolver.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_dns_resolver" "test" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_network_id  = azurerm_virtual_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this DNS Resolver. Changing this forces a new DNS Resolver to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the DNS Resolver should exist. Changing this forces a new DNS Resolver to be created.

* `location` - (Required) Specifies the Azure Region where the DNS Resolver should exist. Changing this forces a new DNS Resolver to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network that is linked to the DNS Resolver. Changing this forces a new DNS Resolver to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the DNS Resolver.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the DNS Resolver.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS Resolver.
* `read` - (Defaults to 5 minutes) Used when retrieving the DNS Resolver.
* `update` - (Defaults to 30 minutes) Used when updating the DNS Resolver.
* `delete` - (Defaults to 30 minutes) Used when deleting the DNS Resolver.

## Import

DNS Resolver can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_resolver.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsResolvers/dnsResolver1
```
