---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ip_group"
description: |-
    Lists ip group resources.
---

# List resource: azurerm_ip_group

Lists ip group resources.

## Example Usage

### List all ip groups

```hcl
list "azurerm_ip_group" "example" {
  provider = azurerm
  config {
  }
}
```

### List all ip groups in a resource group

```hcl
list "azurerm_ip_group" "example" {
  provider = azurerm
  config {
    resource_group_name = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The id of the ip subscription to query.

* `resource_group_name` - (Optional) The name of the ip resource group to query.
