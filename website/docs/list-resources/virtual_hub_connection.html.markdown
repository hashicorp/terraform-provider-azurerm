---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_connection"
description: |-
    Lists Virtual Hub Connection resources.
---

# List resource: azurerm_virtual_hub_connection

Lists Virtual Hub Connection resources.

## Example Usage

### List Virtual Hub Connections in a Virtual Hub

```hcl
list "azurerm_virtual_hub_connection" "example" {
  provider = azurerm
  config {
    virtual_hub_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/hub1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `virtual_hub_id` - (Required) The ID of the Virtual Hub to query.
