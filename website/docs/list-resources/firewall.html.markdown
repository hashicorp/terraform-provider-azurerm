---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall"
description: |-
    Lists firewall resources.
---

# List resource: azurerm_firewall

Lists firewall resources.

## Example Usage

### List all firewalls

```hcl
list "azurerm_firewall" "example" {
  provider = azurerm
  config {
  }
}
```

### List all firewalls in a resource group

```hcl
list "azurerm_firewall" "example" {
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
