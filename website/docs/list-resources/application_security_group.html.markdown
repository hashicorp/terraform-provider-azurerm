---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_security_group"
description: |-
    Lists application security group resources.
---

# List resource: azurerm_application_security_group

Lists application security group resources.

## Example Usage

### List all application security groups

```hcl
list "azurerm_application_security_group" "example" {
  provider = azurerm
  config {
  }
}
```

### List all application security groups in a resource group

```hcl
list "azurerm_application_security_group" "example" {
  provider = azurerm
  config {
    resource_group_name = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The id of the application subscription to query.

* `resource_group_name` - (Optional) The name of the application resource group to query.
