---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_forwarding_rule"
description: |-
  Manages a Private DNS Resolver Forwarding Rule.
---

# azurerm_private_dns_resolver_forwarding_rule

Manages a Private DNS Resolver Forwarding Rule.

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
  name                                       = "example-drdfr"
  resource_group_name                        = azurerm_resource_group.example.name
  location                                   = azurerm_resource_group.example.location
  private_dns_resolver_outbound_endpoint_ids = [azurerm_private_dns_resolver_outbound_endpoint.example.id]
}

resource "azurerm_private_dns_resolver_forwarding_rule" "example" {
  name                      = "example-rule"
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.example.id
  domain_name               = "onprem.local."
  enabled                   = true
  target_dns_servers {
    ip_address = "10.10.0.1"
    port       = 53
  }
  metadata = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Private DNS Resolver Forwarding Rule. Changing this forces a new Private DNS Resolver Forwarding Rule to be created.

* `dns_forwarding_ruleset_id` - (Required) Specifies the ID of the Private DNS Resolver Forwarding Ruleset. Changing this forces a new Private DNS Resolver Forwarding Rule to be created.

* `domain_name` - (Required) Specifies the domain name for the Private DNS Resolver Forwarding Rule. Changing this forces a new Private DNS Resolver Forwarding Rule to be created.

* `target_dns_servers` - (Required) Can be specified multiple times to define multiple target DNS servers. Each `target_dns_servers` block as defined below.

* `enabled` - (Optional) Specifies the state of the Private DNS Resolver Forwarding Rule. Defaults to `true`.

* `metadata` - (Optional) Metadata attached to the Private DNS Resolver Forwarding Rule.

---

A `target_dns_servers` block supports the following:

* `ip_address` - (Required) DNS server IP address.

* `port` - (Optional) DNS server port.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Forwarding Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS Resolver Forwarding Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Forwarding Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS Resolver Forwarding Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS Resolver Forwarding Rule.

## Import

Private DNS Resolver Forwarding Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_resolver_forwarding_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsForwardingRulesets/dnsForwardingRuleset1/forwardingRules/forwardingRule1
```
