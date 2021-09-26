---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_notebook"
description: |-
  Manages a Synapse Notebook.
---

# azurerm_synapse_notebook

Manages a Synapse Notebook.

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

resource "azurerm_synapse_notebook" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  cells                = <<BODY
[
            {
              "cell_type": "code",
              "metadata": {},
              "source": [
                "def my_function():\n",
                " print(\"Hello from a function\")\n",
                "\n",
                "my_function()"
              ],
              "attachments": {},
              "outputs": [
                {
                  "execution_count": 3,
                  "output_type": "execute_result",
                  "data": {
                    "text/plain": "Hello from a function"
                  },
                  "metadata": {}
                }
              ]
            }
          ]
BODY

  //description     = "test"
  display_name    = "notebook example"
  language        = "python"
  major_version   = 4
  minor_version   = 2
  spark_pool_name = azurerm_synapse_spark_pool.example.name
  session_config {
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

* `name` - (Required) The name which should be used for this Synapse Notebook. Changing this forces a new Synapse Notebook to be created.

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace. Changing this forces a new Synapse Notebook to be created.

* `cells` - (Required)  Array of cells of the current notebook.
---

* `codemirror_mode` - (Optional) The codemirror mode to use for code in this language.

* `description` - (Optional) The description for the Synapse Notebook.

* `display_name` - (Optional) The display name for the Synapse Notebook.

* `language` - (Optional) The programming language which this kernel runs.

* `major_version` - (Optional) Notebook format (major number). Incremented between backwards incompatible changes to the notebook format.

* `minor_version` - (Optional) Notebook format (minor number). Incremented for backward compatible changes to the notebook format.

* `session_config` - (Optional) A `session_config` block as defined below.

* `spark_pool_name` - (Optional) Reference Synapse Spark Pool name.

---

A `session_config` block supports the following:

* `driver_cores` - (Optional) Number of cores to use for the driver.

* `driver_memory` - (Optional) Amount of memory to use for the driver process.

* `executor_cores` - (Optional) Number of cores to use for each executor.

* `executor_memory` - (Optional) Amount of memory to use per executor process.

* `num_executors` - (Optional) Number of executors to launch for this session.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse Notebook.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Notebook.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Notebook.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Notebook.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Notebook.

## Import

Synapse Notebooks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_notebook.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/Notebooks/notebook1
```
