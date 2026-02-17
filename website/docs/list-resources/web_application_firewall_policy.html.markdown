---
subcategory: "Web"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_application_firewall_policy"
description: |-
    Lists Web Application Firewall Policy resources.
---

# List resource: azurerm_web_application_firewall_policy

Lists Web Application Firewall Policy resources.

## Example Usage

### List all Web Application Firewall Policys

```hcl
list "azurerm_web_application_firewall_policy" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Web Application Firewall Policys in a Resource Group

```hcl
list "azurerm_web_application_firewall_policy" "example" {
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
