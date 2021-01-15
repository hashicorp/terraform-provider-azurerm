---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_sentinel_alert_rule"
description: |-
  Gets information about an existing Sentinel Alert Rule.
---

# Data Source: azurerm_sentinel_alert_rule

Use this data source to access information about an existing Sentinel Alert Rule.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  resource_group_name = "example-resources"
}

data "azurerm_sentinel_alert_rule" "example" {
  name                       = "existing"
  log_analytics_workspace_id = data.azurerm_log_analytics_workspace.example.id
}

output "id" {
  value = data.azurerm_sentinel_alert_rule.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel Alert Rule.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace this Sentinel Alert Rule belongs to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Sentinel Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Alert Rule.
