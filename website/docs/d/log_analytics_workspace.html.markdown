---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace"
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
  value = data.azurerm_log_analytics_workspace.example.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Log Analytics Workspace.
* `resource_group_name` - The name of the resource group in which the Log Analytics workspace is located in.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Log Analytics Workspace.

* `primary_shared_key` - The Primary shared key for the Log Analytics Workspace.

* `secondary_shared_key` - The Secondary shared key for the Log Analytics Workspace.

* `workspace_id` - The Workspace (or Customer) ID for the Log Analytics Workspace.

* `sku` - The Sku of the Log Analytics Workspace.

* `retention_in_days` - The workspace data retention in days.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Workspace.
