---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_stack_hci_logical_network"
description: |-
  Gets information about an existing Stack HCI Logical Network.
---

# Data Source: azurerm_stack_hci_logical_network

Use this data source to access information about an existing Stack HCI Logical Network.

## Example Usage

```hcl
data "azurerm_stack_hci_logical_network" "example" {
  name = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_stack_hci_logical_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Stack HCI Logical Network. Changing this forces a new Stack HCI Logical Network to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stack HCI Logical Network exists. Changing this forces a new Stack HCI Logical Network to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stack HCI Logical Network.

* `custom_location_id` - The ID of the Custom Location where the Azure Stack HCI Logical Network exists.

* `dns_servers` - A `dns_servers` block as defined below.

* `location` - The Azure Region where the Stack HCI Logical Network exists.

* `subnet` - A `subnet` block as defined below.

* `tags` - A mapping of tags assigned to the Stack HCI Logical Network.

* `virtual_switch_name` - The name of the virtual switch on the cluster associates with the Azure Stack HCI Logical Network.

---

A `ip_pool` block exports the following:

* `end` - The IPv4 address of the end of the IP address pool.

* `start` - The IPv4 address of the start of the IP address pool.

---

A `route` block exports the following:

* `address_prefix` - The address prefix in CIDR notation.

* `name` - The name of this route.

* `next_hop_ip_address` - The IPv4 address of the next hop.

---

A `subnet` block exports the following:

* `address_prefix` - The address prefix in CIDR notation.

* `ip_allocation_method` - The IP address allocation method of the subnet.

* `ip_pool` - A `ip_pool` block as defined above.

* `route` - A `route` block as defined above.

* `vlan_id` - The VLAN ID of the Logical Network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Stack HCI Logical Network.

