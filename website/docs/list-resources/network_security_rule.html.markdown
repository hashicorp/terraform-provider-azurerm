---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_rule"
description: |-
    Lists network security rule resources.
---

# List resource: azurerm_network_security_rule

Lists network security rule resources.

## Example Usage

### List network security rules in a network network security group

```hcl
list "azurerm_network_security_rule" "example" {
  provider = azurerm
  config {
    network_security_group_id = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `network_security_group_id` - (Required) The id of the network security group to query.

