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

* `content` - The content of the Runbook.

* `description` - The description of the Runbook.

* `draft` - A `draft` block as defined below.

* `id` - The Automation Runbook ID.

* `job_schedule` - One or more `job_schedule` block as defined below.

* `location` - The Azure Region where the Runbook exists.

* `log_activity_trace_level` - The activity-level tracing of the Runbook.

* `log_progress` - The Progress log option of the Runbook.

* `log_verbose` - The Verbose log option of the Runbook.

* `publish_content_link` -  One `publish_content_link` block as defined below.

* `runbook_type` - The type of Runbook.

* `runtime_environment_name` - The runtime environment name for the runbook.

* `tags` - A mapping of tags assigned to the Runbook.

---

The `publish_content_link` block supports the following:

* `uri` -  The URI of the runbook content.

* `version` -  Specifies the version of the content

* `hash` - A `hash` block as defined below.

The `draft` block supports:

* `edit_mode_enabled` -  Whether the draft in edit mode.

* `content_link` - A `publish_content_link` block as defined above.

* `output_types` - Specifies the output types of the runbook.

* `parameters` - A list of `parameters` block as defined below.

---

The `parameters` block supports:

* `key` -  The name of the parameter.

* `type` -  Specifies the type of this parameter.

* `mandatory` -  Whether this parameter is mandatory.

* `position` -  Specifies the position of the parameter.

* `default_value` -  Specifies the default value of the parameter.

---

The `job_schedule` block supports:

* `schedule_name` -  The name of the Schedule.

* `parameters` -  A map of key/value pairs corresponding to the arguments that can be passed to the Runbook.

-> **Note:** The parameter keys/names must strictly be in lowercase, even if this is not the case in the runbook. This is due to a limitation in Azure Automation where the parameter names are normalized. The values specified don't have this limitation.

* `run_on` - Name of a Hybrid Worker Group the Runbook will be executed on.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Automation` - 2024-10-23
