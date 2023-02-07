---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_virtual_network_link"
description: |-
  Manages a Private DNS Resolver Virtual Network Link.
---

# azurerm_private_dns_resolver_virtual_network_link

Manages a Private DNS Resolver Virtual Network Link.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "west europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "outbounddns"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.64/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azurerm_private_dns_resolver" "example" {
  name                = "example-resolver"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_network_id  = azurerm_virtual_network.example.id
}

resource "azurerm_private_dns_resolver_outbound_endpoint" "example" {
  name                    = "example-endpoint"
  private_dns_resolver_id = azurerm_private_dns_resolver.example.id
  location                = azurerm_private_dns_resolver.example.location
  subnet_id               = azurerm_subnet.example.id
  tags = {
    key = "value"
  }
}

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "example" {
  name                                       = "example-ruleset"
  resource_group_name                        = azurerm_resource_group.example.name
  location                                   = azurerm_resource_group.example.location
  private_dns_resolver_outbound_endpoint_ids = [azurerm_private_dns_resolver_outbound_endpoint.example.id]
  tags = {
    key = "value"
  }
}

resource "azurerm_private_dns_resolver_virtual_network_link" "example" {
  name                      = "example-link"
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.example.id
  virtual_network_id        = azurerm_virtual_network.example.id
  metadata = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Private DNS Resolver Virtual Network Link. Changing this forces a new Private DNS Resolver Virtual Network Link to be created.

* `dns_forwarding_ruleset_id` - (Required) Specifies the ID of the Private DNS Resolver DNS Forwarding Ruleset. Changing this forces a new Private DNS Resolver Virtual Network Link to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network that is linked to the Private DNS Resolver Virtual Network Link. Changing this forces a new resource to be created.

* `metadata` - (Optional) Metadata attached to the Private DNS Resolver Virtual Network Link.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Virtual Network Link.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS Resolver Virtual Network Link.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Virtual Network Link.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS Resolver Virtual Network Link.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS Resolver Virtual Network Link.

## Import

Private DNS Resolver Virtual Network Link can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_resolver_virtual_network_link.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsForwardingRulesets/dnsForwardingRuleset1/virtualNetworkLinks/virtualNetworkLink1
```
