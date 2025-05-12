---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_dns_servers"
description: |-
  Manages the DNS servers associated with a virtual network.
---

# azurerm_virtual_network_dns_servers

Manages the DNS servers associated with a virtual network.

~> **Note:** Terraform currently provides both a standalone [virtual network DNS Servers resource](virtual_network_dns_servers.html), and allows for DNS servers to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line DNS servers in conjunction with any Virtual Network DNS Servers resources. Doing so will cause a conflict of Virtual Network DNS Servers configurations and will overwrite virtual networks DNS servers.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_virtual_network_dns_servers" "example" {
  virtual_network_id = azurerm_virtual_network.example.id
  dns_servers        = ["10.7.7.2", "10.7.7.7", "10.7.7.1"]
}
```

## Argument Reference

The following arguments are supported:

* `virtual_network_id` - (Required) The ID of the Virtual Network that should be linked to the DNS Zone. Changing this forces a new resource to be created.

* `dns_servers` - (Optional) List of IP addresses of DNS servers

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The virtual network DNS server ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Network.

## Import

Virtual Network DNS Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_network_dns_servers.exampleNetwork /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/dnsServers/default
```
