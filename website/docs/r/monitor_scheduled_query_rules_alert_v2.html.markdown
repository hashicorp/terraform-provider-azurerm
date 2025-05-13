---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_alert_v2"
description: |-
  Manages an AlertingAction Scheduled Query Rules Version 2 resource within Azure Monitor
---

# azurerm_monitor_scheduled_query_rules_alert_v2

Manages an AlertingAction Scheduled Query Rules Version 2 resource within Azure Monitor

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "example" {
  name                = "example-mag"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "test mag"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_application_insights.example.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "example" {
  name                = "example-msqrv2"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  evaluation_frequency = "PT10M"
  window_duration      = "PT10M"
  scopes               = [azurerm_application_insights.example.id]
  severity             = 4
  criteria {
    query                   = <<-QUERY
      requests
        | summarize CountByCountry=count() by client_CountryOrRegion
      QUERY
    time_aggregation_method = "Maximum"
    threshold               = 17.5
    operator                = "LessThan"

    resource_id_column    = "client_CountryOrRegion"
    metric_measure_column = "CountByCountry"
    dimension {
      name     = "client_CountryOrRegion"
      operator = "Exclude"
      values   = ["123"]
    }
    failing_periods {
      minimum_failing_periods_to_trigger_alert = 1
      number_of_evaluation_periods             = 1
    }
  }

  auto_mitigation_enabled          = true
  workspace_alerts_storage_enabled = false
  description                      = "example sqr"
  display_name                     = "example-sqr"
  enabled                          = true
  query_time_range_override        = "PT1H"
  skip_query_validation            = true
  action {
    action_groups = [azurerm_monitor_action_group.example.id]
    custom_properties = {
      key  = "value"
      key2 = "value2"
    }
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id,
    ]
  }
  tags = {
    key  = "value"
    key2 = "value2"
  }

  depends_on = [azurerm_role_assignment.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Monitor Scheduled Query Rule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Monitor Scheduled Query Rule should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Monitor Scheduled Query Rule should exist. Changing this forces a new resource to be created.

* `criteria` - (Required) A `criteria` block as defined below.

* `evaluation_frequency` - (Optional) How often the scheduled query rule is evaluated, represented in ISO 8601 duration format. Possible values are `PT1M`, `PT5M`, `PT10M`, `PT15M`, `PT30M`, `PT45M`, `PT1H`, `PT2H`, `PT3H`, `PT4H`, `PT5H`, `PT6H`, `P1D`.

-> **Note:** `evaluation_frequency` cannot be greater than the query look back which is `window_duration`*`number_of_evaluation_periods`.

-> **Note:** `evaluation_frequency` cannot be greater than the `mute_actions_after_alert_duration`.

* `scopes` - (Required) Specifies the list of resource IDs that this scheduled query rule is scoped to. Changing this forces a new resource to be created. Currently, the API supports exactly 1 resource ID in the scopes list.

* `severity` - (Required) Severity of the alert. Should be an integer between 0 and 4. Value of 0 is severest.

* `window_duration` - (Required) Specifies the period of time in ISO 8601 duration format on which the Scheduled Query Rule will be executed (bin size). If `evaluation_frequency` is `PT1M`, possible values are `PT1M`, `PT5M`, `PT10M`, `PT15M`, `PT30M`, `PT45M`, `PT1H`, `PT2H`, `PT3H`, `PT4H`, `PT5H`, and `PT6H`. Otherwise, possible values are `PT5M`, `PT10M`, `PT15M`, `PT30M`, `PT45M`, `PT1H`, `PT2H`, `PT3H`, `PT4H`, `PT5H`, `PT6H`, `P1D`, and `P2D`.

* `action` - (Optional) An `action` block as defined below.

* `auto_mitigation_enabled` - (Optional) Specifies the flag that indicates whether the alert should be automatically resolved or not. Value should be `true` or `false`. The default is `false`.

* `workspace_alerts_storage_enabled` - (Optional) Specifies the flag which indicates whether this scheduled query rule check if storage is configured. Value should be `true` or `false`. The default is `false`.

* `description` - (Optional) Specifies the description of the scheduled query rule.

* `display_name` - (Optional) Specifies the display name of the alert rule.

* `enabled` - (Optional) Specifies the flag which indicates whether this scheduled query rule is enabled. Value should be `true` or `false`. Defaults to `true`.

* `mute_actions_after_alert_duration` - (Optional) Mute actions for the chosen period of time in ISO 8601 duration format after the alert is fired. Possible values are `PT5M`, `PT10M`, `PT15M`, `PT30M`, `PT45M`, `PT1H`, `PT2H`, `PT3H`, `PT4H`, `PT5H`, `PT6H`, `P1D` and `P2D`.

-> **Note:** `auto_mitigation_enabled` and `mute_actions_after_alert_duration` are mutually exclusive and cannot both be set.

* `query_time_range_override` - (Optional) Set this if the alert evaluation period is different from the query time range. If not specified, the value is `window_duration`*`number_of_evaluation_periods`. Possible values are `PT5M`, `PT10M`, `PT15M`, `PT20M`, `PT30M`, `PT45M`, `PT1H`, `PT2H`, `PT3H`, `PT4H`, `PT5H`, `PT6H`, `P1D` and `P2D`.

-> **Note:** `query_time_range_override` cannot be less than the query look back which is `window_duration`*`number_of_evaluation_periods`.

* `skip_query_validation` - (Optional) Specifies the flag which indicates whether the provided query should be validated or not. The default is false.

* `tags` - (Optional) A mapping of tags which should be assigned to the Monitor Scheduled Query Rule.

* `target_resource_types` - (Optional) List of resource type of the target resource(s) on which the alert is created/updated. For example if the scope is a resource group and targetResourceTypes is `Microsoft.Compute/virtualMachines`, then a different alert will be fired for each virtual machine in the resource group which meet the alert criteria.

* `identity` - (Optional) An `identity` block as defined below.

---

An `action` block supports the following:

* `action_groups` - (Optional) List of Action Group resource IDs to invoke when the alert fires.

* `custom_properties` - (Optional) Specifies the properties of an alert payload.

---

A `criteria` block supports the following:

* `operator` - (Required) Specifies the criteria operator. Possible values are `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`,and `LessThanOrEqual`.

* `query` - (Required) The query to run on logs. The results returned by this query are used to populate the alert.

* `threshold` - (Required) Specifies the criteria threshold value that activates the alert.

* `time_aggregation_method` - (Required) The type of aggregation to apply to the data points in aggregation granularity. Possible values are `Average`, `Count`, `Maximum`, `Minimum`,and `Total`.

* `dimension` - (Optional) A `dimension` block as defined below.

* `failing_periods` - (Optional) A `failing_periods` block as defined below.

* `metric_measure_column` - (Optional) Specifies the column containing the metric measure number.

-> **Note:** `metric_measure_column` is required if `time_aggregation_method` is `Average`, `Maximum`, `Minimum`, or `Total`. And `metric_measure_column` can not be specified if `time_aggregation_method` is `Count`.

* `resource_id_column` - (Optional) Specifies the column containing the resource ID. The content of the column must be an uri formatted as resource ID.

---

A `dimension` block supports the following:

* `name` - (Required) Name of the dimension.

* `operator` - (Required) Operator for dimension values. Possible values are `Exclude`,and `Include`.

* `values` - (Required) List of dimension values. Use a wildcard `*` to collect all.

---

A `failing_periods` block supports the following:

* `minimum_failing_periods_to_trigger_alert` - (Required) Specifies the number of violations to trigger an alert. Should be smaller or equal to `number_of_evaluation_periods`. Possible value is integer between 1 and 6.

* `number_of_evaluation_periods` - (Required) Specifies the number of aggregated look-back points. The look-back time window is calculated based on the aggregation granularity `window_duration` and the selected number of aggregated points. Possible value is integer between 1 and 6.

-> **Note:** The query look back which is `window_duration`*`number_of_evaluation_periods` cannot exceed 48 hours.

-> **Note:** `number_of_evaluation_periods` must be `1` for queries that do not project timestamp column

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Scheduled Query Rule. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Scheduled Query Rule.

~> **Note:** This is required when `type` is set to `UserAssigned`. The identity associated must have required roles, read the [Azure documentation](https://learn.microsoft.com/en-us/azure/azure-monitor/alerts/alerts-create-log-alert-rule#configure-the-alert-rule-details) for more information.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Monitor Scheduled Query Rule.

* `created_with_api_version` - The api-version used when creating this alert rule.

* `is_a_legacy_log_analytics_rule` - True if this alert rule is a legacy Log Analytic Rule.

* `is_workspace_alerts_storage_configured` - The flag indicates whether this Scheduled Query Rule has been configured to be stored in the customer's storage.

* `identity` - An `identity` block as defined below.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this App Service slot.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this App Service slot.

-> **Note:** You can access the Principal ID via `azurerm_monitor_scheduled_query_rules_alert_v2.example.identity[0].principal_id` and the Tenant ID via `azurerm_monitor_scheduled_query_rules_alert_v2.example.identity[0].tenant_id`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Monitor Scheduled Query Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Monitor Scheduled Query Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Monitor Scheduled Query Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Monitor Scheduled Query Rule.

## Import

Monitor Scheduled Query Rule Alert can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_scheduled_query_rules_alert_v2.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Insights/scheduledQueryRules/rule1
```
