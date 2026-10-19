---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_network_rules"
description: |-
    Lists Storage Account Network Rules resources.
---

# List resource: azurerm_storage_account_network_rules

Lists Storage Account Network Rules resources.

## Example Usage

### List all Storage Account Network Ruless in the subscription

```hcl
list "azurerm_storage_account_network_rules" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Storage Account Network Ruless in a Resource Group

```hcl
list "azurerm_storage_account_network_rules" "example" {
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
