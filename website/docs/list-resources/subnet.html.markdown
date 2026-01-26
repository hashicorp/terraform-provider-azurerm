---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet"
description: |-
  Lists Subnet resources.
---

# List resource: azurerm_subnet

Lists Subnet resources.

## Example Usage

### List all Subnets for a Specific Virtual Network

```hcl
list "azurerm_subnet" "example" {
  provider = azurerm
  config {
    virtual_network_id = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{virtualNetworkName}"
  }
}
```

### List all Subnets for Returned Virtual Networks

```hcl
list "azurerm_virtual_network" "example" {
  provider = azurerm

  include_resource = true

  config {
    resource_group_name = "example"
  }
}

list "azurerm_subnet" "example" {
  for_each = toset([for vnet in list.azurerm_virtual_network.list.data : vnet.state.id])

  provider = azurerm
  config {
    virtual_network_id = each.key
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `virtual_network_id` - (Optional) The ID of the virtual network for which to list subnets.
