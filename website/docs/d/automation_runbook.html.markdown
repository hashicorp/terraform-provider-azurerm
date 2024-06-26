---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_runbook"
description: |-
  Gets information about an existing Automation Runbook.
---

# Data Source: azurerm_automation_runbook

Use this data source to access information about an existing Automation Runbook.

## Example Usage

```hcl
data "azurerm_automation_runbook" "example" {
  name                    = "existing-runbook"
  resource_group_name     = "existing"
  automation_account_name = "existing-automation"
}

output "id" {
  value = data.azurerm_automation_runbook.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Automation Runbook.

* `automation_account_name` - (Required) The name of the Automation Account the runbook belongs to.

* `resource_group_name` - (Required) The name of the Resource Group where the Automation exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The Automation Runbook ID.

* `content` - The content of the Runbook.

* `description` - The description of the Runbook.

* `location` - The Azure Region where the Runbook exists.

* `log_activity_trace_level` - The activity-level tracing of the Runbook.

* `log_progress` - The Progress log option of the Runbook.

* `log_verbose` - The Verbose log option of the Runbook.

* `runbook_type` - The type of Runbook.

* `tags` - A mapping of tags assigned to the Runbook.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
