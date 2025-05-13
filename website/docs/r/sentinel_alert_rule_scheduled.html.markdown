---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_alert_rule_scheduled"
description: |-
  Manages a Sentinel Scheduled Alert Rule.
---

# azurerm_sentinel_alert_rule_scheduled

Manages a Sentinel Scheduled Alert Rule.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_sentinel_alert_rule_scheduled" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.example.workspace_id
  display_name               = "example"
  severity                   = "High"
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller
QUERY
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel Scheduled Alert Rule. Changing this forces a new Sentinel Scheduled Alert Rule to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace this Sentinel Scheduled Alert Rule belongs to. Changing this forces a new Sentinel Scheduled Alert Rule to be created.

* `display_name` - (Required) The friendly name of this Sentinel Scheduled Alert Rule.

* `severity` - (Required) The alert severity of this Sentinel Scheduled Alert Rule. Possible values are `High`, `Medium`, `Low` and `Informational`.

* `query` - (Required) The query of this Sentinel Scheduled Alert Rule.

---

* `alert_details_override` - (Optional) An `alert_details_override` block as defined below.

* `alert_rule_template_guid` - (Optional) The GUID of the alert rule template which is used for this Sentinel Scheduled Alert Rule. Changing this forces a new Sentinel Scheduled Alert Rule to be created.

* `alert_rule_template_version` - (Optional) The version of the alert rule template which is used for this Sentinel Scheduled Alert Rule.

* `custom_details` - (Optional) A map of string key-value pairs of columns to be attached to this Sentinel Scheduled Alert Rule. The key will appear as the field name in alerts and the value is the event parameter you wish to surface in the alerts.

* `description` - (Optional) The description of this Sentinel Scheduled Alert Rule.

* `enabled` - (Optional) Should the Sentinel Scheduled Alert Rule be enabled? Defaults to `true`.

* `entity_mapping` - (Optional) A list of `entity_mapping` blocks as defined below.

* `event_grouping` - (Optional) A `event_grouping` block as defined below.

* `incident` - (Optional) A `incident` block as defined below.

* `query_frequency` - (Optional) The ISO 8601 timespan duration between two consecutive queries. Defaults to `PT5H`.

* `query_period` - (Optional) The ISO 8601 timespan duration, which determine the time period of the data covered by the query. For example, it can query the past 10 minutes of data, or the past 6 hours of data. Defaults to `PT5H`.

-> **Note:** `query_period` must larger than or equal to `query_frequency`, which ensures there is no gaps in the overall query coverage.

* `suppression_duration` - (Optional) If `suppression_enabled` is `true`, this is ISO 8601 timespan duration, which specifies the amount of time the query should stop running after alert is generated. Defaults to `PT5H`.

-> **Note:** `suppression_duration` must larger than or equal to `query_frequency`, otherwise the suppression has no actual effect since no query will happen during the suppression duration.

* `suppression_enabled` - (Optional) Should the Sentinel Scheduled Alert Rulea stop running query after alert is generated? Defaults to `false`.

* `sentinel_entity_mapping` - (Optional) A list of `sentinel_entity_mapping` blocks as defined below.

-> **Note:** `entity_mapping` and `sentinel_entity_mapping` together can't exceed 10.

* `tactics` - (Optional) A list of categories of attacks by which to classify the rule. Possible values are `Collection`, `CommandAndControl`, `CredentialAccess`, `DefenseEvasion`, `Discovery`, `Execution`, `Exfiltration`, `ImpairProcessControl`, `InhibitResponseFunction`, `Impact`, `InitialAccess`, `LateralMovement`, `Persistence`, `PrivilegeEscalation`, `PreAttack`, `Reconnaissance` and `ResourceDevelopment`.

* `techniques` - (Optional) A list of techniques of attacks by which to classify the rule.

* `trigger_operator` - (Optional) The alert trigger operator, combined with `trigger_threshold`, setting alert threshold of this Sentinel Scheduled Alert Rule. Possible values are `Equal`, `GreaterThan`, `LessThan`, `NotEqual`. Defaults to `GreaterThan`.

* `trigger_threshold` - (Optional) The baseline number of query results generated, combined with `trigger_operator`, setting alert threshold of this Sentinel Scheduled Alert Rule. Defaults to `0`.

---

An `alert_details_override` block supports the following:

* `description_format` - (Optional) The format containing columns name(s) to override the description of this Sentinel Alert Rule.

* `display_name_format` - (Optional) The format containing columns name(s) to override the name of this Sentinel Alert Rule.

* `severity_column_name` - (Optional) The column name to take the alert severity from.

* `tactics_column_name` - (Optional) The column name to take the alert tactics from.

* `dynamic_property` - (Optional) A list of `dynamic_property` blocks as defined below.

---

A `dynamic_property` block supports the following:

* `name` - (Required) The name of the dynamic property. Possible Values are `AlertLink`, `ConfidenceLevel`, `ConfidenceScore`, `ExtendedLinks`, `ProductComponentName`, `ProductName`, `ProviderName`, `RemediationSteps` and `Techniques`.

* `value` - (Required) The value of the dynamic property. Pssible Values are `Caller`, `dcount_ResourceId` and `EventSubmissionTimestamp`.

---

An `entity_mapping` block supports the following:

* `entity_type` - (Required) The type of the entity. Possible values are `Account`, `AzureResource`, `CloudApplication`, `DNS`, `File`, `FileHash`, `Host`, `IP`, `Mailbox`, `MailCluster`, `MailMessage`, `Malware`, `Process`, `RegistryKey`, `RegistryValue`, `SecurityGroup`, `SubmissionMail`, `URL`.

* `field_mapping` - (Required) A list of `field_mapping` blocks as defined below.

---

A `sentinel_entity_mapping` block supports the following:

* `column_name` - (Required) The column name to be mapped to the identifier.

---

A `event_grouping` block supports the following:

* `aggregation_method` - (Required) The aggregation type of grouping the events. Possible values are `AlertPerResult` and `SingleAlert`.

---

A `field_mapping` block supports the following:

* `identifier` - (Required) The identifier of the entity.

* `column_name` - (Required) The column name to be mapped to the identifier.

---

A `incident` block supports the following:

* `create_incident_enabled` - (Required) Whether to create an incident from alerts triggered by this Sentinel Scheduled Alert Rule?

* `grouping` - (Required) A `grouping` block as defined below.

---

A `grouping` block supports the following:

* `enabled` - (Optional) Enable grouping incidents created from alerts triggered by this Sentinel Scheduled Alert Rule. Defaults to `true`.

* `lookback_duration` - (Optional) Limit the group to alerts created within the lookback duration (in ISO 8601 duration format). Defaults to `PT5M`.

* `reopen_closed_incidents` - (Optional) Whether to re-open closed matching incidents? Defaults to `false`.

* `entity_matching_method` - (Optional) The method used to group incidents. Possible values are `AnyAlert`, `Selected` and `AllEntities`. Defaults to `AnyAlert`.

* `by_entities` - (Optional) A list of entity types to group by, only when the `entity_matching_method` is `Selected`. Possible values are `Account`, `AzureResource`, `CloudApplication`, `DNS`, `File`, `FileHash`, `Host`, `IP`, `Mailbox`, `MailCluster`, `MailMessage`, `Malware`, `Process`, `RegistryKey`, `RegistryValue`, `SecurityGroup`, `SubmissionMail`, `URL`.

* `by_alert_details` - (Optional) A list of alert details to group by, only when the `entity_matching_method` is `Selected`. Possible values are `DisplayName` and `Severity`.

* `by_custom_details` - (Optional) A list of custom details keys to group by, only when the `entity_matching_method` is `Selected`. Only keys defined in the `custom_details` may be used.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Sentinel Scheduled Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Scheduled Alert Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Scheduled Alert Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel Scheduled Alert Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Scheduled Alert Rule.

## Import

Sentinel Scheduled Alert Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_alert_rule_scheduled.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRules/rule1
```
