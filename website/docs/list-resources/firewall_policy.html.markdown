---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_policy"
description: |-
    Lists firewall policy resources.
---

# List resource: azurerm_firewall_policy

Lists firewall policy resources.

## Example Usage

### List all firewall policys

```hcl
list "azurerm_firewall_policy" "example" {
  provider = azurerm
  config {
  }
}
```

### List all firewall policys in a resource group

```hcl
list "azurerm_firewall_policy" "example" {
  provider = azurerm
  config {
    resource_group_name = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The id of the firewall subscription to query.

* `resource_group_name` - (Optional) The name of the firewall resource group to query.
