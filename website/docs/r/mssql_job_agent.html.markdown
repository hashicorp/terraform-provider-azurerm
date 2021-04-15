---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_job_agent"
description: |-
  Manages an Elastic Job Agent.
---

# azurerm_mssql_job_agent

Manages an Elastic Job Agent.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "northeurope"
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Elastic Job Agent. Changing this forces a new Elastic Job Agent to be created.

* `location` - (Required) The Azure Region where the Elastic Job Agent should exist. Changing this forces a new Elastic Job Agent to be created.

* `database_id` - (Required) The ID of the database to store metadata for the Elastic Job Agent. Changing this forces a new Elastic Job Agent to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Database.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic Job Agent.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Database.
* `update` - (Defaults to 1 hour) Used when updating the Database.
* `delete` - (Defaults to 1 hour) Used when deleting the Database.

## Import

Elastic Job Agents can be imported using the `id`, e.g.

```shell
terraform import azurerm_mssql_job_agent.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Sql/servers/myserver1/jobAgents/myjobagent1
```
