---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_job_credential"
description: |-
  Manages an Elastic Job Credential.
---

# azurerm_mssql_job_credential

Manages an Elastic Job Credential.

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

resource "azurerm_mssql_job_credential" "example" {
  name         = "example-credential"
  job_agent_id = azurerm_mssql_job_agent.example.id
  username     = "my-username"
  password     = "MyP4ssw0rd!!!"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Elastic Job Credential. Changing this forces a new Elastic Job Credential to be created.

* `job_agent_id` - (Required) The ID of the Elastic Job Agent. Changing this forces a new Elastic Job Credential to be created.

* `username` - (Required) The username part of the credential.

* `password` - (Required) The password part of the credential.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic Job Credential.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Elastic Job Credential.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic Job Credential.
* `update` - (Defaults to 1 hour) Used when updating the Elastic Job Credential.
* `delete` - (Defaults to 1 hour) Used when deleting the Elastic Job Credential.

## Import

Elastic Job Credentials can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_job_credential.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Sql/servers/myserver1/jobAgents/myjobagent1/credentials/credential1
```
