---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_action_rule_suppression"
description: |-
  Manages an Monitor Action Rule which type is suppression.
---

# azurerm_monitor_action_rule_suppression

Manages an Monitor Action Rule which type is suppression.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_action_rule_suppression" "example" {
  name                = "example-amar"
  resource_group_name = azurerm_resource_group.example.name

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.example.id]
  }

  suppression {
    recurrence_type = "Weekly"

    schedule {
      start_date_utc    = "2019-01-01T01:02:03Z"
      end_date_utc      = "2019-01-03T15:02:07Z"
      recurrence_weekly = ["Sunday", "Monday", "Friday", "Saturday"]
    }
  }

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Monitor Action Rule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group in which the Monitor Action Rule should exist. Changing this forces a new resource to be created.

* `suppression` - (Required) A `suppression` block as defined below.

* `description` - (Optional) Specifies a description for the Action Rule.

* `enabled` - (Optional) Is the Action Rule enabled? Defaults to `true`.

* `scope` - (Optional) A `scope` block as defined below.

* `condition` - (Optional) A `condition` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `suppression` block supports the following:

* `recurrence_type` - (Required) Specifies the type of suppression. Possible values are `Always`, `Daily`, `Monthly`, `Once`, and `Weekly`.

* `schedule` - (Optional) A `schedule` block as defined below. Required if `recurrence_type` is `Daily`, `Monthly`, `Once` or `Weekly`.

---

The `schedule` block supports the following:

* `start_date_utc` - (Required) specifies the recurrence UTC start datetime (Y-m-d'T'H:M:S'Z').

* `end_date_utc` - (Required) specifies the recurrence UTC end datetime (Y-m-d'T'H:M:S'Z').

* `recurrence_weekly` - (Optional) specifies the list of dayOfWeek to recurrence. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and  `Saturday`.

* `recurrence_monthly` - (Optional) specifies the list of dayOfMonth to recurrence. Possible values are between `1` - `31`. Required if `recurrence_type` is `Monthly`.

---

The `scope` block supports the following:

* `type` - (Required) Specifies the type of target scope. Possible values are `ResourceGroup` and `Resource`.

* `resource_ids` - (Required) A list of resource IDs of the given scope type which will be the target of action rule.

---

The `condition` block supports the following:

* `alert_context` - (Optional) A `alert_context` block as defined below.

* `alert_rule_id` - (Optional) A `alert_rule_id` block as defined below.

* `description` - (Optional) A `description` block as defined below.

* `monitor` - (Optional) A `monitor` block as defined below.

* `monitor_service` - (Optional) A `monitor_service` as block defined below.

* `severity` - (Optional) A `severity` block as defined below.

* `target_resource_type` - (Optional) A `target_resource_type` block as defined below.

---

The `alert_context` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition.

---

The `alert_rule_id` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition.

---

The `description` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains`, and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition.

---

The `monitor` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals` and `NotEquals`.

* `values` - (Required) A list of values to match for a given condition. Possible values are `Fired` and `Resolved`.

---

The `monitor_service` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals` and `NotEquals`.

* `values` - (Required) A list of values to match for a given condition. Possible values are `ActivityLog Administrative`, `ActivityLog Autoscale`, `ActivityLog Policy`, `ActivityLog Recommendation`, `ActivityLog Security`, `Application Insights`, `Azure Backup`, `Data Box Edge`, `Data Box Gateway`, `Health Platform`, `Log Analytics`, `Platform`, and `Resource Health`.

---

The `severity` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals`and `NotEquals`.

* `values` - (Required) A list of values to match for a given condition. Possible values are `Sev0`, `Sev1`, `Sev2`, `Sev3`, and `Sev4`.

---

The `target_resource_type` block supports the following:

* `operator` - (Required) The operator for a given condition. Possible values are `Equals` and `NotEquals`.

* `values` - (Required) A list of values to match for a given condition. The values should be valid resource types.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Monitor Action Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Monitor Action Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Monitor Action Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Monitor Action Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Monitor Action Rule.

## Import

Monitor Action Rule can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_monitor_action_rule_suppression.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AlertsManagement/actionRules/actionRule1
```
