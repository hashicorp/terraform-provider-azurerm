---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall"
description: |-
  Manages an Azure Firewall.

---

# azurerm_firewall

Manages an Azure Firewall.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "North Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "testvnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "testpip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "example" {
  name                = "testfirewall"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.example.id
    public_ip_address_id = azurerm_public_ip.example.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Firewall. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the resource. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `ip_configuration` - (Required) An `ip_configuration` block as documented below.

* `dns_servers` - (Optional) A list of DNS servers that the Azure Firewall will direct DNS traffic to the for name resolution.

* `management_ip_configuration` - (Optional) A `management_ip_configuration` block as documented below, which allows force-tunnelling of traffic to be performed by the firewall. Adding or removing this block or changing the `subnet_id` in an existing block forces a new resource to be created.

* `threat_intel_mode` - (Optional) The operation mode for threat intelligence-based filtering. Possible values are: `Off`, `Alert` and `Deny`. Defaults to `Alert`

* `zones` - (Optional) Specifies the availability zones in which the Azure Firewall should be created. Changing this forces a new resource to be created.

-> **Please Note**: Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview).

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `ip_configuration` block supports the following:

* `name` - (Required) Specifies the name of the IP Configuration.

* `subnet_id` - (Optional) Reference to the subnet associated with the IP Configuration.

-> **NOTE** The Subnet used for the Firewall must have the name `AzureFirewallSubnet` and the subnet mask must be at least a `/26`.

-> **NOTE** At least one and only one `ip_configuration` block may contain a `subnet_id`.

* `public_ip_address_id` - (Required) The ID of the Public IP Address associated with the firewall.

-> **NOTE** The Public IP must have a `Static` allocation and `Standard` sku.

---

A `management_ip_configuration` block supports the following:

* `name` - (Required) Specifies the name of the IP Configuration.

* `subnet_id` - (Required) Reference to the subnet associated with the IP Configuration. Changing this forces a new resource to be created.

-> **NOTE** The Management Subnet used for the Firewall must have the name `AzureFirewallManagementSubnet` and the subnet mask must be at least a `/26`.

* `public_ip_address_id` - (Required) The ID of the Public IP Address associated with the firewall.

-> **NOTE** The Public IP must have a `Static` allocation and `Standard` sku.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Azure Firewall.

* `ip_configuration` - A `ip_configuration` block as defined below.

---

A `ip_configuration` block exports the following:

* `private_ip_address` - The Private IP address of the Azure Firewall.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Firewall.
* `update` - (Defaults to 90 minutes) Used when updating the Firewall.
* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall.
* `delete` - (Defaults to 90 minutes) Used when deleting the Firewall.

## Import

Azure Firewalls can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/azureFirewalls/testfirewall
```
