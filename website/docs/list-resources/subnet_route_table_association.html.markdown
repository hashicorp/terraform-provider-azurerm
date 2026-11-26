---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet_route_table_association"
description: |-
    Lists Subnet Route Table Association resources.
---

# List resource: azurerm_subnet_route_table_association

Lists Subnet Route Table Association resources.

## Example Usage

### List Subnet Route Table Associations in a Virtual Network

```hcl
list "azurerm_subnet_route_table_association" "example" {
  provider = azurerm
  config {
    virtual_network_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `virtual_network_id` - (Required) The ID of the Virtual Network to query.
