---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ip"
description: |-
  Gets information about an existing Public IP Address.

---

# Data Source: azurerm_public_ip

Use this data source to access information about an existing Public IP Address.

## Example Usage (reference an existing)

```hcl
data "azurerm_public_ip" "example" {
  name                = "name_of_public_ip"
  resource_group_name = "name_of_resource_group"
}

output "domain_name_label" {
  value = data.azurerm_public_ip.example.domain_name_label
}

output "public_ip_address" {
  value = data.azurerm_public_ip.example.ip_address
}
```

## Example Usage (Retrieve the Dynamic Public IP of a new VM)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "test-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "test-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "acctsub"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                    = "test-pip"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  allocation_method       = "Dynamic"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}

resource "azurerm_network_interface" "example" {
  name                = "test-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.2.5"
    public_ip_address_id          = azurerm_public_ip.example.id
  }
}

resource "azurerm_virtual_machine" "example" {
  name                  = "test-vm"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  # ...
}

data "azurerm_public_ip" "example" {
  name                = azurerm_public_ip.example.name
  resource_group_name = azurerm_virtual_machine.example.resource_group_name
}

output "public_ip_address" {
  value = data.azurerm_public_ip.example.ip_address
}
```

## Argument Reference

* `name` - Specifies the name of the public IP address.
* `resource_group_name` - Specifies the name of the resource group.

## Attributes Reference

* `id` - The ID of the Public IP address.
* `allocation_method` - The allocation method for this IP address. Possible values are `Static` or `Dynamic`.
* `domain_name_label` - The label for the Domain Name.
* `idle_timeout_in_minutes` - Specifies the timeout for the TCP idle connection.
* `ddos_protection_mode` - The DDoS protection mode of the public IP.
* `ddos_protection_plan_id` - The ID of DDoS protection plan associated with the public IP. 
* `fqdn` - Fully qualified domain name of the A DNS record associated with the public IP. This is the concatenation of the domainNameLabel and the regionalized DNS zone.
* `ip_address` - The IP address value that was allocated.
* `ip_version` - The IP version being used, for example `IPv4` or `IPv6`.
* `location` - The region that this public ip exists.
* `reverse_fqdn` - The fully qualified domain name that resolves to this public IP address.
* `sku` - The SKU of the Public IP.
* `ip_tags` - A mapping of tags to assigned to the resource.
* `tags` - A mapping of tags to assigned to the resource.
* `zones` - A list of Availability Zones in which this Public IP is located.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Public IP Address.
