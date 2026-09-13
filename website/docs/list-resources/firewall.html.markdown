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

### List all Firewalls in the subscription

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
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query. Defaults to the value specified in the Provider Configuration.

* `resource_group_name` - (Optional) The name of the Resource Group to query.
