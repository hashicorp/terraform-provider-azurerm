---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_job"
description: |-
  Manages an Elastic Job.
---

# azurerm_mssql_job

Manages an Elastic Job.

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
  description  = "example description"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Elastic Job. Changing this forces a new Elastic Job to be created.

* `job_agent_id` - (Required) The ID of the Elastic Job Agent. Changing this forces a new Elastic Job to be created.

---

* `description` - (Optional) The description of the Elastic Job.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic Job.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Elastic Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic Job.
* `update` - (Defaults to 30 minutes) Used when updating the Elastic Job.
* `delete` - (Defaults to 30 minutes) Used when deleting the Elastic Job.

## Import

Elastic Jobs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_job.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Sql/servers/myserver1/jobAgents/myjobagent1/jobs/myjob1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Sql`: 2023-08-01-preview
