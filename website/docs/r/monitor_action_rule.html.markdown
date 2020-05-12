---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_action_rule"
description: |-
  Manages an Monitor Action Rule.
---

# azurerm_monitor_action_rule

Manages an Monitor Action Rule.

## Diagnostics Action Rule Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_action_rule" "example" {
  name                = "example-amar"
  resource_group_name = azurerm_resource_group.example.name
  type                = "Diagnostics"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.example.id]
  }

  tags = {
    foo = "bar"
  }
}
```

## Suppression Action Rule Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_action_rule" "example" {
  name                = "example-amar"
  resource_group_name = azurerm_resource_group.example.name
  type                = "Suppression"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.example.id]
  }

  suppression {
    recurrence_type = "Weekly"

    schedule {
      start_date = "12/09/2018"
      start_time = "06:00:00"
      end_date   = "12/18/2018"
      end_time   = "14:00:00"

      recurrence_weekly = ["Sunday", "Monday", "Friday", "Saturday"]
    }
  }

  tags = {
    foo = "bar"
  }
}
```

## ActionGroup Action Rule Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_action_group" "example" {
  name                = "example-action-group"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "exampleactiongroup"
}

resource "azurerm_monitor_action_rule" "example" {
  name                = "example-amar"
  resource_group_name = azurerm_resource_group.example.name
  type                = "ActionGroup"
  action_group_id     = azurerm_monitor_action_group.example.id

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.example.id]
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

* `type` - (Required) Specifies the type of the Action Rule. Possible values are `Suppression`, `ActionGroup` and `Diagnostics`. Changing this forces a new resource to be created.

* `description` - (Optional) Specifies a description for the Action Rule.

* `enabled` - (Optional) Is the Action Rule enabled? Defaults to `true`.

* `action_group_id` - (Optional) Specifies the resource id of monitor action group. Required only if `type` is `ActionGroup`.

* `suppression` - (Optional) One `suppression` as defined below. Required only if `type` is `Suppression`.

* `scope` - (Optional) One `scope` block as defined below.

* `condition` - (Optional) One `condition` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `suppression` block supports the following:

* `recurrence_type` - (Required) Specifies the type of suppression. Possible values are `Always`, `Daily`, `Monthly`, `Once` and `Weekly`.

* `schedule` - (Optional) One `schedule` block as defined below. Required if `recurrence_type` is `Daily`, `Monthly`, `Once` or `Weekly`.
---

The `scope` block supports the following:

* `type` - (Required) Specifies the type of target scope. Possible values are `ResourceGroup` and `Resource`.

* `resource_ids` - (Required) A list of resource IDs of the given scope type which will be the target of action rule.

---

The `condition` block supports the following:

* `alert_context` - (Optional) One `alert_context` block as defined below.

* `alert_rule_id` - (Optional) One `alert_rule_id` block as defined below.

* `description` - (Optional) One `description` block as defined below.

* `monitor` - (Optional) One `monitor` block as defined below.

* `monitor_service` - (Optional) One `monitor_service` as block defined below.

* `severity` - (Optional) One `severity` block as defined below.

* `target_resource_type` - (Optional) One `target_resource_type` block as defined below.

---

The `schedule` block supports the following:

* `start_date` - (Required) specifies the UTC start date of recurrence range. The format should be `MM/DD/YYYY`.

* `start_time` - (Required) specifies the UTC start time of recurrence range. The format should be `HH:MM:SS`.

* `end_date` - (Required) specifies the UTC end date of recurrence range. The format should be `MM/DD/YYYY`.

* `end_time` - (Required) specifies the UTC end time of recurrence range. The format should be `HH:MM:SS`.

* `recurrence_weekly` - (Optional) specifies the list of dayOfWeek to recurrence. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and  `Saturday`.

* `recurrence_monthly` - (Optional) specifies the list of dayOfMonth to recurrence. Possible values are between `1` - `31`. Required if `recurrence_type` is `Monthly`.

---

The `alert_context` block supports the following:

* `operator` - (Required) operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains` and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition.

---

The `alert_rule_id` block supports the following:

* `operator` - (Required) operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains` and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition.

---

The `description` block supports the following:

* `operator` - (Required) operator for a given condition. Possible values are `Equals`, `NotEquals`, `Contains` and `DoesNotContain`.

* `values` - (Required) A list of values to match for a given condition.

---

The `monitor` block supports the following:

* `operator` - (Required) operator for a given condition. Possible values are `Equals`, `NotEquals`.

* `values` - (Required) A list of values to match for a given condition. Possible values are `Fired`, `Resolved`.

---

The `monitor_service` block supports the following:

* `operator` - (Required) operator for a given condition. Possible values are `Equals`, `NotEquals`.

* `values` - (Required) A list of values to match for a given condition. Possible values are `ActivityLog Administrative`, `ActivityLog Autoscale`, `ActivityLog Policy`, `ActivityLog Recommendation`, `ActivityLog Security`, `Application Insights`, `Azure Backup`, `Data Box Edge`, `Data Box Gateway`, `Health Platform`, `Log Analytics`, `Platform` and `Resource Health`.

---

The `severity` block supports the following:

* `operator` - (Required) operator for a given condition. Possible values are `Equals`, `NotEquals`.

* `values` - (Required) A list of values to match for a given condition. Possible values are `Sev0`, `Sev1`, `Sev2`, `Sev3` and `Sev4`.

---

The `target_resource_type` block supports the following:

* `operator` - (Required) operator for a given condition. Possible values are `Equals`, `NotEquals`.

* `values` - (Required) list of values to match for a given condition. The values should be valid resource types.

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
$ terraform import azurerm_monitor_action_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AlertsManagement/actionRules/actionRule1
```
