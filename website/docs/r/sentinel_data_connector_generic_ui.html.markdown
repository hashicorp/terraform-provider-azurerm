---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_generic_ui"
description: |-
  Manages a Generic UI Data Connector.
---

# azurerm_sentinel_data_connector_generic_ui

Manages a Generic UI Data Connector.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "east us"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  resource_group_name = azurerm_resource_group.example.name
  workspace_name      = azurerm_log_analytics_workspace.example.name
}

resource "azurerm_sentinel_data_connector_generic_ui" "example" {
  name                       = "example-generic-ui"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id

  title                    = "Slack"
  publisher                = "Slack"
  description_markdown     = "The [Slack](https://slack.com) data connector provides the capability to ingest [Slack Audit Records](https://api.slack.com/admins/audit-logs) events into Microsoft Sentinel through the REST API. Refer to [API documentation](https://api.slack.com/admins/audit-logs#the_audit_event) for more information. The connector provides ability to get events which helps to examine potential security risks, analyze your team's use of collaboration, diagnose configuration problems and more. This data connector uses Microsoft Sentinel native polling capability."
  graph_queries_table_name = "SlackAuditNativePoller_CL"
  graph_query {
    metric_name = "Total data received"
    legend      = "Slack audit events"
    base_query  = "{{graphQueriesTableName}}"
  }

  sample_query {
    description = "All Slack audit events"
    query       = "{{graphQueriesTableName}}\n| sort by TimeGenerated desc"
  }

  data_type {
    name                     = "{{graphQueriesTableName}}"
    last_data_received_query = "{{graphQueriesTableName}}\n            | summarize Time = max(TimeGenerated)\n            | where isnotempty(Time)"
  }

  connectivity_criteria {
    type = "IsConnectedQuery"
  }

  availability {
    enabled = true
    preview = true
  }

  permissions {
    resource_provider {
      name         = "Microsoft.OperationalInsights/workspaces"
      display_name = "read and write permissions are required."
      display_text = "Workspace"
      scope        = "Workspace"
      required_permissions {
        read   = true
        write  = true
        delete = true
      }
    }
    custom {
      name        = "Slack API credentials"
      description = "**SlackAPIBearerToken** is required for REST API. [See the documentation to learn more about API](https://api.slack.com/web#authentication). Check all [requirements and follow the instructions](https://api.slack.com/web#authentication) for obtaining credentials."
    }
  }

  instruction {
    title       = "Connect Slack to Microsoft Sentinel"
    description = "Enable Slack audit Logs"
    step {
      type = "InfoMessage"
      parameters = jsonencode({
        "enable" = "true"
      })
    }
  }

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this API Polling Data Connector. Changing this forces a new API Polling Data Connector to be created.

* `log_analytics_workspace_id` - (Required) The ID of the TODO. Changing this forces a new API Polling Data Connector to be created.

* `title` - (Required) Title displayed in the data connector page.

* `graph_queries_table_name` - (Required) The name of the Log Analytics table from which data for your queries is pulled. The table name can be any string, but must end in `_CL`. For example: `TableName_CL`.

* `publisher` - (Required) Publisher displayed in the data connector page.

* `availability` - (Optional) An `availability` block as defined below. Defines the availability of the Data Connector.

* `connectivity_criteria` - (Optional) A `connectivity_criteria` block as defined below. Defines the way the connector check connectivity.

* `custom_image` - (Optional) An optional custom image to be used when displaying the connector within Azure Sentinel's connector's gallery.

* `data_type` - (Optional) One or more `data_type` blocks as defined below. Defines Data types to check for last data received

* `description_markdown` - (Optional) The description in markdown format of the Data Connector.

* `graph_query` - (Optional) One or more `graph_query` blocks as defined below. Defines the graph query to show the current data status.

* `instruction` - (Optional) One or more `instruction` blocks as defined below. Define steps to enable the connector.

* `permission` - (Optional) One or more `permission` blocks as defined below. Define permission required by the Data Connector.

* `sample_query` - (Optional) One or more `sample_query` blocks as defined below. Define the sample queries for the Data Connector.

---

A `graph_query` block supports the following:

* `base_query` - (Required) The base query for the graph.

* `legend` - (Required) The legend for the graph.

* `metric_name` - (Required) The name of metric that the query is checking.

---

A `instruction` block supports the following:

* `title` - (Required) The title of the instruction.

* `description` - (Optional) The description of the instruction.

* `step` - (Optional) One or more `step` blocks as defined below.

---

A `step` block supports the following:

* `type` - (Required) The type of the setting. Possible values are `CopyableLabel`, `InfoMessage` and `InstructionStepsGroup`.

* `parameters` - (Optional) The parameters for the setting.

---

A `permissions` block supports the following:

* `resource_provider` - (Required) One or more `resource_provider` blocks as defined below.

* `custom` - (Optional) One or more `custom` blocks as defined above.

---

A `resource_provider` block supports the following:

* `name` - (Required) The name of the Resource Provider. Possible values are `Microsoft.Authorization/policyAssignments`, `Microsoft.OperationalInsights/solutions`, `Microsoft.OperationalInsights/workspaces`, Microsoft.OperationalInsights/workspaces/datasources`, `Microsoft.OperationalInsights/workspaces/sharedKeys` and `microsoft.aadiam/diagnosticSettings`.

* `permissions_display_name` - (Required) The display name of the provider.

* `permissions_display_text` - (Required) The description text of the provider.

* `required_permissions` - (Required) A `required_permissions` block as defined above.

* `scope` - (Required) The scope of the provider. Possible values are `PermissionProviderScopeResourceGroup`, `PermissionProviderScopeSubscription` and `PermissionProviderScopeWorkspace`.

---

A `required_permissions` block supports the following:

* `action` - (Optional) Whether require action permission.

* `delete` - (Optional) Whether require delete permission.

* `read` - (Optional) Whether require read permission.

* `write` - (Optional) Whether require write permission.

---

A `custom` block supports the following:

* `name` - (Required) The name which should be used for this custom permission.

* `description` - (Optional) The description of this custom permission.

---

A `sample_query` block supports the following:

* `description` - (Required) The description of the sample query.

* `query` - (Required) the sample query.

---

A `availability` block supports the following:

* `enabled` - (Required) Should the Data Connector be enabled?

* `preview` - (Required) Whether this Data Connector is in preview?

---

A `connectivity_criteria` block supports the following:

* `type` - (Required) type of connectivity. Possible value is `IsConnectedQuery`. 

* `value` - (Optional) Specifies a list of queries for checking connectivity.

---

A `data_type` block supports the following:

* `name` - (Required) The name which should be used for this Data Type.

* `last_data_received_query` - (Required) Query for indicate last data received.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Polling Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Polling Data Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Polling Data Connector.
* `update` - (Defaults to 30 minutes) Used when updating the API Polling Data Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Polling Data Connector.

## Import

API Polling Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_generic_ui.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
