---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_sentinel_alert_rule_ms_security_incident"
description: |-
  Gets information about an existing Sentinel MS Security Incident Alert Rule.
---

# Data Source: azurerm_sentinel_alert_rule_ms_security_incident

Use this data source to access information about an existing Sentinel MS Security Incident Alert Rule.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_sentinel_alert_rule_ms_security_incident" "example" {
  name                       = "existing"
  log_analytics_workspace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1"
}

output "id" {
  value = data.azurerm_sentinel_alert_rule_ms_security_incident.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel MS Security Incident Alert Rule.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace this Sentinel MS Security Incident Alert Rule belongs to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Sentinel MS Security Incident Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel MS Security Incident Alert Rule.
