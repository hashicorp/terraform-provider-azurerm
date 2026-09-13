---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_peering"
description: |-
    Lists Virtual Network Peering resources.
---

# List resource: azurerm_virtual_network_peering

Lists Virtual Network Peering resources.

## Example Usage

### List Virtual Network Peerings in a Virtual Network

```hcl
list "azurerm_virtual_network_peering" "example" {
  provider = azurerm
  config {
    virtual_network_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `virtual_network_id` - (Required) The ID of the Virtual Network to query.
