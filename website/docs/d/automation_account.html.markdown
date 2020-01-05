---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_account"
sidebar_current: "docs-azurerm-datasource-automation"
description: |-
  Gets information about an existing Automation Account.
---

# Data Source: azurerm_automation_account

Use this data source to access information about an existing Automation Account.

## Example Usage

```hcl
data "azurerm_automation_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resources"
}
output "automation_account_id" {
  value = "${data.azurerm_automation_account.example.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the Automation Account.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Automation Account exists.

## Attributes Reference

* `id` - The ID of the Automation Account

* `primary_key` - The Primary Access Key for the Automation Account.

* `secondary_key` - The Secondary Access Key for the Automation Account.

* `endpoint` - The Endpoint for this Auomation Account.
