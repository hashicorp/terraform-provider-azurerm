---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace"
sidebar_current: "docs-azurerm-datasource-oms-log-analytics-workspace"
description: |-
  Gets information about an existing Log Analytics (formally Operational Insights) Workspace.
---

# Data Source: azurerm_log_analytics_workspace

Use this data source to access information about an existing Log Analytics (formally Operational Insights) Workspace.

## Example Usage

```hcl
data "azurerm_log_analytics_workspace" "example" {
  name                = "acctest-01"
  resource_group_name = "acctest"
}

output "log_analytics_workspace_id" {
  value = "${data.azurerm_log_analytics_workspace.example.workspace_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Log Analytics Workspace.
* `resource_group_name` - (Required) The name of the resource group in which the Log Analytics workspace is located in.

## Attributes Reference

The following attributes are exported:

* `id` - The Azure Resource ID of the Log Analytics Workspace.

* `primary_shared_key` - The Primary shared key for the Log Analytics Workspace.

* `secondary_shared_key` - The Secondary shared key for the Log Analytics Workspace.

* `workspace_id` - The Workspace (or Customer) ID for the Log Analytics Workspace.

* `portal_url` - The Portal URL for the Log Analytics Workspace.

* `sku` - The Sku of the Log Analytics Workspace.

* `retention_in_days` - The workspace data retention in days.

* `tags` - A mapping of tags assigned to the resource.
