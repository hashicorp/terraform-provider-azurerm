---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_connection_type"
description: |-
  Lists Automation Connection Type resources.
---

# List resource: azurerm_automation_connection_type

Lists Automation Connection Type resources.

## Example Usage

### List all Automation Connection Types in a specific Automation Connection

```hcl
list "azurerm_automation_connection_type" "example" {
  provider = azurerm
  config {
    resource_group_name     = "example-rg"
    automation_account_name = "example-automation-account"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Required) The name of the resource group to query.

* `automation_account_name` - (Required) The name of the automation account to query

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
