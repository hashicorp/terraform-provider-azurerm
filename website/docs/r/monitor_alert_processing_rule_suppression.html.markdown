---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_alert_processing_rule_suppression"
description: |-
  Manages an Alert Processing Rule which suppress notifications.
---

# azurerm_monitor_alert_processing_rule_suppression

Manages an Alert Processing Rule which suppress notifications.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_alert_processing_rule_suppression" "example" {
  name                = "example"
  resource_group_name = "example"
  scopes              = [azurerm_resource_group.example.id]

  condition {
    target_resource_type {
      operator = "Equals"
      values   = ["Microsoft.Compute/VirtualMachines"]
    }
    severity {
      operator = "Equals"
      values   = ["Sev0", "Sev1", "Sev2"]
    }
  }

  schedule {
    effective_from  = "2022-01-01T01:02:03"
    effective_until = "2022-02-02T01:02:03"
    time_zone       = "Pacific Standard Time"
    recurrence {
      daily {
        start_time = "17:00:00"
        end_time   = "09:00:00"
      }
      weekly {
        days_of_week = ["Saturday", "Sunday"]
      }
    }
  }

  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Alert Processing Rule. Changing this forces a new Alert Processing Rule to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Alert Processing Rule should exist. Changing this forces a new Alert Processing Rule to be created.

* `scopes` - (Required) A list of resource IDs which will be the target of Alert Processing Rule.

---

* `condition` - (Optional) A `condition` block as defined below.

* `description` - (Optional) Specifies a description for the Alert Processing Rule.

* `enabled` - (Optional) Should the Alert Processing Rule be enabled? Defaults to `true`.

* `schedule` - (Optional) A `schedule` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Alert Processing Rule.

---

A `alert_context` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) Specifies a list of values to match for a given condition.

---

A `alert_rule_id` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) Specifies a list of values to match for a given condition.

---

A `alert_rule_name` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) Specifies a list of values to match for a given condition.

---

A `condition` block supports the following:

* `alert_context` - (Optional) A `alert_context` block as defined above.

* `alert_rule_id` - (Optional) A `alert_rule_id` block as defined above.

* `alert_rule_name` - (Optional) A `alert_rule_name` block as defined above.

* `description` - (Optional) A `description` block as defined below.

* `monitor_condition` - (Optional) A `monitor_condition` block as defined below.

* `monitor_service` - (Optional) A `monitor_service` block as defined below.

* `severity` - (Optional) A `severity` block as defined below.

* `signal_type` - (Optional) A `signal_type` block as defined below.

* `target_resource` - (Optional) A `target_resource` block as defined below.

* `target_resource_group` - (Optional) A `target_resource_group` block as defined below.

* `target_resource_type` - (Optional) A `target_resource_type` block as defined below.

---

A `daily` block supports the following:

* `start_time` - (Required) Specifies the recurrence start time (H:M:S).

* `end_time` - (Required) Specifies the recurrence end time (H:M:S).

---

A `description` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) Specifies a list of values to match for a given condition.

---

A `monitor_condition` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals` and `NotEquals`.

* `values` - (Required) Specifies a list of values to match for a given condition. Possible values are `Fired` and `Resolved`.

---

A `monitor_service` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals` and `NotEquals`.

* `values` - (Required) A list of values to match for a given condition. Possible values are `ActivityLog Administrative`, `ActivityLog Autoscale`, `ActivityLog Policy`, `ActivityLog Recommendation`, `ActivityLog Security`, `Application Insights`, `Azure Backup`, `Azure Stack Edge`, `Azure Stack Hub`, `Custom`, `Data Box Gateway`, `Health Platform`, `Log Alerts V2`, `Log Analytics`, `Platform`, `Prometheus`, `Resource Health`, `Smart Detector`, and `VM Insights - Health`.

---

A `monthly` block supports the following:

* `days_of_month` - (Required) Specifies a list of dayOfMonth to recurrence. Possible values are integers between `1` - `31`.

* `start_time` - (Optional) Specifies the recurrence start time (H:M:S).

* `end_time` - (Optional) Specifies the recurrence end time (H:M:S).

---

A `recurrence` block supports the following:

* `daily` - (Optional) One or more `daily` blocks as defined above.

* `weekly` - (Optional) One or more `weekly` blocks as defined below.

* `monthly` - (Optional) One or more `monthly` blocks as defined above.

---

A `schedule` block supports the following:

* `effective_from` - (Optional) Specifies the Alert Processing Rule effective start time (Y-m-d'T'H:M:S).

* `effective_until` - (Optional) Specifies the Alert Processing Rule effective end time (Y-m-d'T'H:M:S).

* `recurrence` - (Optional) A `recurrence` block as defined above.

* `time_zone` - (Optional) The time zone (e.g. Pacific Standard time, Eastern Standard Time). Defaults to `UTC`. [possible values are defined here](https://docs.microsoft.com/en-us/previous-versions/windows/embedded/ms912391(v=winembedded.11)).

---

A `severity` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals` and `NotEquals`.

* `values` - (Required) Specifies list of values to match for a given condition. Possible values are `Sev0`, `Sev1`, `Sev2`, `Sev3`, and `Sev4`.

---

A `signal_type` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals` and `NotEquals`.

* `values` - (Required) Specifies a list of values to match for a given condition. Possible values are `Metric`, `Log`, `Unknown`, and `Health`.

---

A `target_resource` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition. The values should be valid resource IDs.

---

A `target_resource_group` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition. The values should be valid resource group IDs.

---

A `target_resource_type` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition. The values should be valid resource types. (e.g. Microsoft.Compute/VirtualMachines)

---

A `weekly` block supports the following:

* `days_of_week` - (Required) Specifies a list of dayOfWeek to recurrence. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, and `Saturday`.

* `start_time` - (Optional) Specifies the recurrence start time (H:M:S).

* `end_time` - (Optional) Specifies the recurrence end time (H:M:S).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Alert Processing Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Alert Processing Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Alert Processing Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Alert Processing Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Alert Processing Rule.

## Import

Alert Processing Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_alert_processing_rule_suppression.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.AlertsManagement/actionRules/actionRule1
```
