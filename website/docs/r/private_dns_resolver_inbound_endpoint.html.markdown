---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_inbound_endpoint"
description: |-
  Manages a Private DNS Resolver Inbound Endpoint.
---

# azurerm_private_dns_resolver_inbound_endpoint

Manages a Private DNS Resolver Inbound Endpoint.

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

resource "azurerm_private_dns_resolver" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_network_id  = azurerm_virtual_network.example.id
}

resource "azurerm_subnet" "example" {
  name                 = "inbounddns"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azurerm_private_dns_resolver_inbound_endpoint" "example" {
  name                    = "example-drie"
  private_dns_resolver_id = azurerm_private_dns_resolver.example.id
  location                = azurerm_private_dns_resolver.example.location
  ip_configurations {
    private_ip_allocation_method = "Dynamic"
    subnet_id                    = azurerm_subnet.example.id
  }
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Private DNS Resolver Inbound Endpoint. Changing this forces a new Private DNS Resolver Inbound Endpoint to be created.

* `private_dns_resolver_id` - (Required) Specifies the ID of the Private DNS Resolver Inbound Endpoint. Changing this forces a new Private DNS Resolver Inbound Endpoint to be created.

* `ip_configurations` - (Required) One `ip_configurations` block as defined below. Changing this forces a new Private DNS Resolver Inbound Endpoint to be created.

* `location` - (Required) Specifies the Azure Region where the Private DNS Resolver Inbound Endpoint should exist. Changing this forces a new Private DNS Resolver Inbound Endpoint to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Private DNS Resolver Inbound Endpoint.

---

An `ip_configurations` block supports the following:

* `subnet_id` - (Required) The subnet ID of the IP configuration.

* `private_ip_address` - (Optional) Private IP address of the IP configuration.

* `private_ip_allocation_method` - (Optional) Private IP address allocation method. Allowed value is `Dynamic` and `Static`. Defaults to `Dynamic`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private DNS Resolver Inbound Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS Resolver Inbound Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Resolver Inbound Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS Resolver Inbound Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS Resolver Inbound Endpoint.

## Import

Private DNS Resolver Inbound Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_resolver_inbound_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsResolvers/dnsResolver1/inboundEndpoints/inboundEndpoint1
```
