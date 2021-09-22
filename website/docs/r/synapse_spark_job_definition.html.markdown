---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_spark_job_definition"
description: |-
  Manages a Synapse Spark Job Definition.
---

# azurerm_synapse_spark_job_definition

Manages a Synapse Spark Job Definition.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true
}

resource "azurerm_synapse_firewall_rule" "example" {
  name                 = "allowAll"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

resource "azurerm_synapse_spark_pool" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
}

resource "azurerm_synapse_spark_job_definition" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  description          = "test"
  language             = "Java"
  spark_pool_name      = azurerm_synapse_spark_pool.example.name
  spark_version        = "2.4"
  job {
    file       = "abfss://test@test.dfs.core.windows.net/artefacts/sample.jar"
    class_name = "dev.test.tools.sample.Main"
    arguments = [
      "exampleArg"
    ]
    jars            = []
    files           = []
    archives        = []
    driver_memory   = "28g"
    driver_cores    = 4
    executor_memory = "28g"
    executor_cores  = 4
    num_executors   = 2
  }

  depends_on = [
    azurerm_synapse_firewall_rule.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Spark Job Definition. Changing this forces a new Synapse Spark Job Definition to be created.

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace. Changing this forces a new Synapse Spark Job Definition to be created.

---

* `description` - (Optional) The description for the Synapse Spark Job.

* `job` - (Optional) A `job` block as defined below.

* `language` - (Optional) The language of the Spark application.

* `spark_pool_name` - (Optional) Reference Synapse Spark Pool name.

* `spark_version` - (Optional) The required Spark version of the application.

---

A `job` block supports the following:

* `archives` - (Optional) Archives to be used in this job.

* `arguments` - (Optional) Command line arguments for the application.

* `class_name` - (Optional) Main class for Java/Scala application.

* `driver_cores` - (Optional) Number of cores to use for the driver.

* `driver_memory` - (Optional) Amount of memory to use for the driver process.

* `executor_cores` - (Optional) Number of cores to use for each executor.

* `executor_memory` - (Optional) Amount of memory to use per executor process.

* `file` - (Optional) File containing the application to execute.

* `files` - (Optional) Files to be used in this job.

* `jars` - (Optional) Jars to be used in this job.

* `num_executors` - (Optional) Number of executors to launch for this job.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse Spark Job Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Spark Job Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Spark Job Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Spark Job Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Spark Job Definition.

## Import

Synapse Spark Job Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_spark_job_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/sparkjobdefinitions/sparkjobdefinition1
```
