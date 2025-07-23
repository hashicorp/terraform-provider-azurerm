---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_network_interface"
description: |-
  Manages an Azure Stack HCI Network Interface.
---

# azurerm_stack_hci_network_interface

Manages an Azure Stack HCI Network Interface.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_stack_hci_logical_network" "example" {
  name                = "example-hci-ln"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["10.0.0.7", "10.0.0.8"]

  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "10.0.0.0/24"
    route {
      name                = "example-route"
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.20.1"
    }
    vlan_id = 123
  }

  tags = {
    foo = "bar"
  }
}

resource "azurerm_stack_hci_network_interface" "example" {
  name                = "example-ni"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  dns_servers         = ["10.0.0.8"]

  ip_configuration {
    private_ip_address = "10.0.0.2"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }

  tags = {
    foo = "bar"
  }

  lifecycle {
    ignore_changes = [mac_address]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Network Interface. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Network Interface should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Stack HCI Network Interface should exist. Changing this forces a new resource to be created.

* `custom_location_id` - (Required) The ID of the Custom Location where the Azure Stack HCI Network Interface should exist. Changing this forces a new resource to be created.

* `ip_configuration` - (Required) An `ip_configuration` block as defined below. Changing this forces a new resource to be created.

* `dns_servers` - (Optional) A list of IPv4 addresses of DNS servers available to VMs deployed in the Network Interface. Changing this forces a new resource to be created.

* `mac_address` - (Optional) The MAC address of the Network Interface. Changing this forces a new resource to be created.

-> **Note:** If `mac_address` is not specified, it will be assigned by the server. If you experience a diff you may need to add this to `ignore_changes`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Stack HCI Network Interface.

---

An `ip_configuration` block supports the following:

* `subnet_id` - (Required) The resource ID of the Stack HCI Logical Network bound to the IP configuration. Changing this forces a new resource to be created.

* `private_ip_address` - (Optional) The IPv4 address of the IP configuration. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The resource ID of the Azure Stack HCI Network Interface.

* `ip_configuration` - An `ip_configuration` as defined below.

---

An `ip_configuration` block exports the following:

* `gateway` - The IPv4 address of the gateway for the Network Interface.

* `prefix_length` - The prefix length for the address of the Network Interface.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Network Interface.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Network Interface.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Network Interface.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Network Interface.

## Import

Azure Stack HCI Network Interfaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_network_interface.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/networkInterfaces/ni1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.AzureStackHCI`: 2024-01-01
