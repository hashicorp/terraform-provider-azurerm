---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall"
description: |-
    Lists Firewall resources.
---

# List resource: azurerm_firewall

Lists Firewall resources.

## Example Usage

### List all Firewalls

```hcl
list "azurerm_firewall" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Firewalls in a Resource Group

```hcl
list "azurerm_firewall" "example" {
  provider = azurerm
  config {
    resource_group_name = "resource_group_name-example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query.

* `resource_group_name` - (Optional) The name of the Resource Group to query.
