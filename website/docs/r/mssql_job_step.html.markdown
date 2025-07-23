---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_job_step"
description: |-
  Manages an Elastic Job Step.
---

# azurerm_mssql_job_step

Manages an Elastic Job Step.

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
    database_name     = azurerm_mssql_database.example.name
    job_credential_id = azurerm_mssql_job_credential.example.id
  }
}

resource "azurerm_mssql_job" "example" {
  name         = "example-job"
  job_agent_id = azurerm_mssql_job_agent.example.id
  description  = "example description"
}

resource "azurerm_mssql_job_step" "test" {
  name                = "example-job-step"
  job_id              = azurerm_mssql_job.example.id
  job_credential_id   = azurerm_mssql_job_credential.example.id
  job_target_group_id = azurerm_mssql_job_target_group.example.id

  job_step_index = 1
  sql_script     = <<EOT
IF NOT EXISTS (SELECT * FROM sys.objects WHERE [name] = N'Pets')
  CREATE TABLE Pets (
    Animal NVARCHAR(50),
    Name NVARCHAR(50),
  );
EOT
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Elastic Job Step. Changing this forces a new Elastic Job Step to be created.

* `job_id` - (Required) The ID of the Elastic Job. Changing this forces a new Elastic Job Step to be created.

* `job_credential_id` - (Required) The ID of the Elastic Job Credential to use when executing this Elastic Job Step.

* `job_step_index` - (Required) The index at which to insert this Elastic Job Step into the Elastic Job.

~> **Note:** This value must be greater than or equal to 1 and less than or equal to the number of job steps in the Elastic Job.

* `job_target_group_id` - (Required) The ID of the Elastic Job Target Group.

* `sql_script` - (Required) The T-SQL script to be executed by this Elastic Job Step.

-> **Note:** While Azure places no restrictions on the script provided here, it is recommended to ensure the script is idempotent.

---

* `initial_retry_interval_seconds` - (Optional) The initial retry interval in seconds. Defaults to `1`.

* `maximum_retry_interval_seconds` - (Optional) The maximum retry interval in seconds. Defaults to `120`.

~> **Note:** `maximum_retry_interval_seconds` must be greater than `initial_retry_interval_seconds`.

* `output_target` - (Optional) An `output_target` block as defined below.

* `retry_attempts` - (Optional) The number of retry attempts. Defaults to `10`.

* `retry_interval_backoff_multiplier` - (Optional) The multiplier for time between retries. Defaults to `2`.

* `timeout_seconds` - (Optional) The execution timeout in seconds for this Elastic Job Step. Defaults to `43200`.

---

A `output_target` block supports the following:

* `job_credential_id` - (Required) The ID of the Elastic Job Credential to use when connecting to the output destination.

* `mssql_database_id` - (Required) The ID of the output database.

* `table_name` - (Required) The name of the output table.

* `schema_name` - (Optional) The name of the output schema. Defaults to `dbo`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic Job Step.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Elastic Job Step.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic Job Step.
* `update` - (Defaults to 30 minutes) Used when updating the Elastic Job Step.
* `delete` - (Defaults to 30 minutes) Used when deleting the Elastic Job Step.

## Import

Elastic Job Steps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_job_step.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Sql/servers/myserver1/jobAgents/myjobagent1/jobs/myjob1/steps/myjobstep1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Sql`: 2023-08-01-preview
