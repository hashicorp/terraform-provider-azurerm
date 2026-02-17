---
subcategory: "Public"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ip"
description: |-
    Lists Public Ip resources.
---

# List resource: azurerm_public_ip

Lists Public Ip resources.

## Example Usage

### List all Public Ips

```hcl
list "azurerm_public_ip" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Public Ips in a Resource Group

```hcl
list "azurerm_public_ip" "example" {
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
