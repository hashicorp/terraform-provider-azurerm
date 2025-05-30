---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_job_target_group"
description: |-
  Manages a Job Target Group.
---

# azurerm_mssql_job_target_group

Manages a Job Target Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "westeurope"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-server"
  location                     = azurerm_resource_group.example.location
  resource_group_name          = azurerm_resource_group.example.name
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
  username     = "testusername"
  password     = "testpassword"
}

resource "azurerm_mssql_job_target_group" "example" {
  name         = "example-target-group"
  job_agent_id = azurerm_mssql_job_agent.example.id

  job_target {
    server_name       = azurerm_mssql_server.example.name
    job_credential_id = azurerm_mssql_job_credential.example.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Job Target Group. Changing this forces a new Job Target Group to be created.

* `job_agent_id` - (Required) The ID of the Elastic Job Agent. Changing this forces a new Job Target Group to be created.

---

* `job_target` - (Optional) One or more `job_target` blocks as defined below.

---

A `job_target` block supports the following:

* `server_name` - (Required) The name of the MS SQL Server.

* `database_name` - (Optional) The name of the MS SQL Database.

~> **Note:** This cannot be set in combination with `elastic_pool_name`.

* `elastic_pool_name` - (Optional) The name of the MS SQL Elastic Pool.

~> **Note:** This cannot be set in combination with `database_name`.

* `job_credential_id` - (Optional) The ID of the job credential to use during execution of jobs.

~> **Note:** This is required when `membership_type` is `Include`, unless `database_name` is set.

* `membership_type` - (Optional) The membership type for this job target. Possible values are `Include` and `Exclude`. Defaults to `Include`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic Job Target Groups.

* `job_target` - One or more `job_target` blocks as defined below.

---

A `job_target` block exports the following:

* `type` - The job target type. This value is computed based on `server_name`, `database_name`, and `elastic_pool_name`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Elastic Job Target Groups.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic Job Target Groups.
* `update` - (Defaults to 30 minutes) Used when updating the Elastic Job Target Groups.
* `delete` - (Defaults to 30 minutes) Used when deleting the Elastic Job Target Groups.

## Import

Job Target Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_job_target_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Sql/servers/myserver1/jobAgents/myjobagent1/targetGroups/mytargetgroup1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Sql`: 2023-08-01-preview
