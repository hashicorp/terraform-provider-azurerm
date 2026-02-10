---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_databricks_workspace_serverless"
description: |-
  Gets information about an existing Databricks Serverless Workspace.
---

# Data Source: azurerm_databricks_workspace_serverless

Use this data source to access information about an existing Databricks Serverless Workspace.

## Example Usage

```hcl
data "azurerm_databricks_workspace_serverless" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_databricks_workspace_serverless.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Databricks Serverless Workspace.

* `resource_group_name` - (Required) The name of the Resource Group where the Databricks Serverless Workspace exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Databricks Serverless Workspace.

* `location` - The Azure Region where the Databricks Serverless Workspace exists.

* `enhanced_security_compliance` - An `enhanced_security_compliance` block as defined below.

* `workspace_id` - Unique ID of this Databricks Serverless Workspace in Databricks management plane.

* `workspace_url` - URL this Databricks Serverless Workspace is accessible on.

* `tags` - A mapping of tags assigned to the Databricks Serverless Workspace.

---

An `enhanced_security_compliance` block exports the following:

* `automatic_cluster_update_enabled` - Whether automatic cluster updates for this Databricks Serverless Workspace is enabled.

* `compliance_security_profile_enabled` - Whether compliance security profile for this Databricks Serverless Workspace is enabled.

* `compliance_security_profile_standards` - A list of standards enforced on this Databricks Serverless Workspace.

* `enhanced_security_monitoring_enabled` - Whether enhanced security monitoring for this Databricks Serverless Workspace is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Serverless Workspace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Databricks` - 2026-01-01
