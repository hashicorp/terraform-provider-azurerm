---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_rule"
description: |-
    Lists Network Security Rule resources.
---

# List resource: azurerm_network_security_rule

Lists Network Security Rule resources.

## Example Usage

### List Network Security Rules in a Network Security Group

```hcl
list "azurerm_network_security_rule" "example" {
  provider = azurerm
  config {
    network_security_group_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkSecurityGroups/mySecurityGroup"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `network_security_group_id` - (Required) The ID of the Network Security Group to query.
