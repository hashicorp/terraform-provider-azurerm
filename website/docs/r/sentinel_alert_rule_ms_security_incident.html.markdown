---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_alert_rule_ms_security_incident"
description: |-
  Manages a Sentinel MS Security Incident Alert Rule.
---

# azurerm_sentinel_alert_rule_ms_security_incident

Manages a Sentinel MS Security Incident Alert Rule.

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

resource "azurerm_sentinel_alert_rule_ms_security_incident" "example" {
  name                       = "example-ms-security-incident-alert-rule"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.example.workspace_id
  product_filter             = "Microsoft Cloud App Security"
  display_name               = "example rule"
  severity_filter            = ["High"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel MS Security Incident Alert Rule. Changing this forces a new Sentinel MS Security Incident Alert Rule to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace this Sentinel MS Security Incident Alert Rule belongs to. Changing this forces a new Sentinel MS Security Incident Alert Rule to be created.

* `display_name` - (Required) The friendly name of this Sentinel MS Security Incident Alert Rule.

* `product_filter` - (Required) The Microsoft Security Service from where the alert will be generated. Possible values are `Azure Active Directory Identity Protection`, `Azure Advanced Threat Protection`, `Azure Security Center`, `Azure Security Center for IoT`, `Microsoft Cloud App Security`, `Microsoft Defender Advanced Threat Protection` and `Office 365 Advanced Threat Protection`.

* `severity_filter` - (Required) Only create incidents from alerts when alert severity level is contained in this list. Possible values are `High`, `Medium`, `Low` and `Informational`.

~> **Note:** At least one of the severity filters need to be set.

---

* `alert_rule_template_guid` - (Optional) The GUID of the alert rule template which is used to create this Sentinel Scheduled Alert Rule. Changing this forces a new Sentinel MS Security Incident Alert Rule to be created.

* `description` - (Optional) The description of this Sentinel MS Security Incident Alert Rule.

* `enabled` - (Optional) Should this Sentinel MS Security Incident Alert Rule be enabled? Defaults to `true`.

* `display_name_filter` - (Optional) Only create incidents when the alert display name contain text from this list, leave empty to apply no filter.

* `display_name_exclude_filter` - (Optional) Only create incidents when the alert display name doesn't contain text from this list.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Sentinel MS Security Incident Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel MS Security Incident Alert Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel MS Security Incident Alert Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel MS Security Incident Alert Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel MS Security Incident Alert Rule.

## Import

Sentinel MS Security Incident Alert Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_alert_rule_ms_security_incident.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRules/rule1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.OperationalInsights`: 2023-12-01-preview
