---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_security_group"
description: |-
    Lists Application Security Group resources.
---

# List resource: azurerm_application_security_group

Lists Application Security Group resources.

## Example Usage

### List all Application Security Groups

```hcl
list "azurerm_application_security_group" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Application Security Groups in a Resource Group

```hcl
list "azurerm_application_security_group" "example" {
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
