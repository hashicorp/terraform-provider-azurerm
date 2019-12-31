---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_registration_info"
sidebar_current: "docs-azurerm-datasource-automation-registration-info"
description: |-
  Gets information about an existing Automation Account Registration Information.
---

# Data Source: azurerm_automation_account_registration_info

Use this data source to access information about an existing Automation Account Registration Information.

## Example Usage

```hcl
data "azurerm_automation_account_registration_info" "example" {
  name                = "automation-account"
  resource_group_name = "automation-resource-group"
}
output "automation_account_id" {
  value = "${data.azurerm_automation_account_registration_info.example.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the Automation Account.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Automation Account exists.

## Attributes Reference

* `id` - The ID of the Automation Account

* `primary_key` - The primary key for the Automation Account Registration information

* `secondary_key` - The primary key for the Automation Account Registration information

* `endpoint` - The Assigned Automation Account Registration endpoint
