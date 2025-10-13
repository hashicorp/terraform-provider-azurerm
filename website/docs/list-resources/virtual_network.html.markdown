---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network"
description: |-
  Lists Virtual Network resources.
---

# List resource: azurerm_virtual_network

~> **Note:** The `azurerm_virtual_network` List Resource is in beta. Its interface and behaviour may change as the feature evolves, and breaking changes are possible. It is offered as a technical preview without compatibility guarantees until Terraform 1.14 is generally available.

Lists Virtual Network resources.

## Example Usage

```hcl
list "azurerm_virtual_network" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following attributes:

* `resource_group_name` - (Required) The name of the resource group to query.
