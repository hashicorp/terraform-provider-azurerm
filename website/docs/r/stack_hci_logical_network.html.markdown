---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_logical_network"
description: |-
  Manages an Azure Stack HCI Logical Network.
---

# azurerm_stack_hci_logical_network

Manages an Azure Stack HCI Logical Network.

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
    vlan_id              = 123
    route {
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "10.0.0.1"
    }
  }

  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Logical Network. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Logical Network should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Stack HCI Logical Network should exist. Changing this forces a new resource to be created.

* `custom_location_id` - (Required) The ID of Custom Location where the Azure Stack HCI Logical Network should exist. Changing this forces a new resource to be created.

* `virtual_switch_name` - (Required) The name of the virtual switch on the cluster used to associate with the Azure Stack HCI Logical Network. Possible switch names can be retrieved by following this [Azure guide](https://learn.microsoft.com/azure-stack/hci/manage/create-logical-networks?tabs=azurecli#prerequisites). Changing this forces a new resource to be created.

* `subnet` - (Required) A `subnet` block as defined below. Changing this forces a new resource to be created.

* `dns_servers` - (Optional) A list of IPv4 addresses of DNS servers available to VMs deployed in the Logical Networks. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Stack HCI Logical Network.

---

A `ip_pool` block supports the following:

* `start` - (Required) The IPv4 address of the start of the IP address pool. Changing this forces a new resource to be created.

* `end` - (Required) The IPv4 address of the end of the IP address pool. Changing this forces a new resource to be created.

---

A `route` block supports the following:

* `address_prefix` - (Optional) The Address in CIDR notation. Changing this forces a new resource to be created.

* `next_hop_ip_address` - (Optional) The IPv4 address of the next hop. Changing this forces a new resource to be created.

* `name` - (Optional) The name of the route. Changing this forces a new resource to be created.

---

A `subnet` block supports the following:

* `ip_allocation_method` - (Required) The IP address allocation method for the subnet. Possible values are `Dynamic` and `Static`. Changing this forces a new resource to be created.

* `address_prefix` - (Optional) The address prefix in CIDR notation. Changing this forces a new resource to be created.

* `ip_pool` - (Optional) One or more `ip_pool` block as defined above. Changing this forces a new resource to be created.

* `route` - (Optional) A `route` block as defined above. Changing this forces a new resource to be created.

* `vlan_id` - (Optional) The VLAN ID for the Logical Network. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The resource ID of the Azure Stack HCI Logical Network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Logical Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Logical Network.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Logical Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Logical Network.

## Import

Azure Stack HCI Logical Networks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_logical_network.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/logicalNetworks/ln1
```
