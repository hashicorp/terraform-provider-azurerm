---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_sentinel_alert_rule_anomaly"
description: |-
  Gets information about an existing Anomaly Alert Rule.
---

# Data Source: azurerm_sentinel_alert_rule_anomaly

Use this data source to access information about an existing Anomaly Alert Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-law"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  workspace_id                 = azurerm_log_analytics_workspace.example.id
  customer_managed_key_enabled = false
}

data "azurerm_sentinel_alert_rule_anomaly" "example" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.example.workspace_id
  display_name               = "Potential data staging"
}

output "id" {
  value = data.azurerm_sentinel_alert_rule_anomaly.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace.

* `name` - (Optional) The guid of this Sentinel Alert Rule Template. Either `display_name` or `name` have to be specified.

* `display_name` - (Optional) The display name of this Sentinel Alert Rule Template. Either `display_name` or `name` have to be specified.

~> **Note:** One of `name` or `display_name` must be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Anomaly Alert Rule.

* `anomaly_settings_version` - The version of the Anomaly Security ML Analytics Settings.

* `anomaly_version` - The anomaly version of the Anomaly Alert Rule.

* `description` - The description of the Anomaly Alert Rule.

* `enabled` - Is the Anomaly Alert Rule enabled?

* `frequency` - The frequency the Anomaly Alert Rule will be run.

* `required_data_connector` - A `required_data_connector` block as defined below.

* `settings_definition_id` - The ID of the anomaly settings definition Id.

* `Mode` - The Mode of the Anomaly Alert Rule.

* `tactics` - A list of categories of attacks by which to classify the rule.

* `techniques` - A list of techniques of attacks by which to classify the rule.

* `multi_select_observation` - A list of `multi_select_observation` blocks as defined below.

* `single_select_observation` - A list of `single_select_observation` blocks as defined below.

* `prioritized_exclude_observation` - A list of `prioritized_exclude_observation` blocks as defined below.

* `threshold_observation` - A list of `threshold_observation` blocks as defined below.

---

A `required_data_connector` block exports the following:

* `connector_id` - The ID of the required Data Connector.

* `data_types` - A list of data types of the required Data Connector.

---

A `multi_select_observation` block exports the following:

* `name` - The name of the multi select observation.

* `description` - The description of the multi select observation.

* `supported_values` - A list of supported values of the multi select observation.

* `values` - A list of values of the single select observation.

---

A `single_select_observation` block exports the following:

* `name` - The name of the single select observation.

* `description` - The description of the single select observation.

* `supported_values` - A list of supported values of the single select observation.

* `value` - The value of the multi select observation.

---

A `prioritized_exclude_observation` block exports the following:

* `name` - The name of the prioritized exclude observation.

* `description` - The description of the prioritized exclude observation.

* `prioritize` - The prioritized value per `description`.

* `exclude` - The excluded value per `description`.

---

A `threshold_observation` block exports the following:

* `name` - The name of the threshold observation.

* `description` - The description of the threshold observation.

* `max` - The max value of the threshold observation.

* `min` - The min value of the threshold observation.

* `value` - The value of the threshold observation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Anomaly Alert Rule.
