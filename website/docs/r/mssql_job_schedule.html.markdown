---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_job_schedule"
description: |-
  Manages an Elastic Job Schedule.
---

# azurerm_mssql_job_schedule

Manages an Elastic Job Schedule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "East US"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-server"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_database" "example" {
  name      = "example-db"
  server_id = azurerm_mssql_server.example.id
  collation = "SQL_Latin1_General_CP1_CI_AS"
  sku_name  = "S1"
}

resource "azurerm_mssql_job_agent" "example" {
  name        = "example-job-agent"
  location    = azurerm_resource_group.example.location
  database_id = azurerm_mssql_database.example.id
}

resource "azurerm_mssql_job_credential" "example" {
  name         = "example-job-credential"
  job_agent_id = azurerm_mssql_job_agent.example.id
  username     = "my-username"
  password     = "MyP4ssw0rd!!!"
}

resource "azurerm_mssql_job" "example" {
  name         = "example-job"
  job_agent_id = azurerm_mssql_job_agent.example.id
}

resource "azurerm_mssql_job_schedule" "example" {
  job_id = azurerm_mssql_job.example.id

  type       = "Recurring"
  enabled    = true
  end_time   = "2025-12-01T00:00:00Z"
  interval   = "PT5M"
  start_time = "2025-01-01T00:00:00Z"
}
```

## Arguments Reference

The following arguments are supported:

* `job_id` - (Required) The ID of the Elastic Job. Changing this forces a new Elastic Job Schedule to be created.

* `type` - (Required) The type of schedule. Possible values are `Once` and `Recurring`.

---

* `enabled` - (Optional) Should the Elastic Job Schedule be enabled? Defaults to `false`.

~> **Note:** When `type` is set to `Once` and `enabled` is set to `true`, it's recommended to add `enabled` to `ignore_changes`. This is because Azure will set `enabled` to `false` once the job has executed.

* `end_time` - (Optional) The end time of the schedule. Must be in RFC3339 format.

* `interval` - (Optional) The interval between job executions. Must be in ISO8601 duration format.

* `start_time` - (Optional) The start time of the schedule. Must be in RFC3339 format.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic Job Schedule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Elastic Job Schedule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic Job Schedule.
* `update` - (Defaults to 30 minutes) Used when updating the Elastic Job Schedule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Elastic Job Schedule.

## Import

Elastic Job Schedules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_job_schedule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Sql/servers/myserver1/jobAgents/myjobagent1/jobs/myjob1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Sql`: 2023-08-01-preview
