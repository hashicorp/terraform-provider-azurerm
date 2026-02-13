---
subcategory: "Ip"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ip_group"
description: |-
    Lists Ip Group resources.
---

# List resource: azurerm_ip_group

Lists Ip Group resources.

## Example Usage

### List all Ip Groups

```hcl
list "azurerm_ip_group" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Ip Groups in a Resource Group

```hcl
list "azurerm_ip_group" "example" {
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
