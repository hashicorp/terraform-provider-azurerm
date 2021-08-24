---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_pipeline"
description: |-
  Manages a Pipeline inside a Azure Synapse Analytics workspace.
---

# azurerm_synapse_pipeline

Manages a Pipeline inside a Azure Synapse Analytics workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
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
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

resource "azurerm_synapse_pipeline" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id

  depends_on = [azurerm_synapse_firewall_rule.example]
}
```

## Example Usage with Activities

```hcl
resource "azurerm_synapse_pipeline" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  variables = {
    "bob" = "item1"
  }
  activities_json = <<JSON
[
	{
		"name": "Append variable1",
		"type": "AppendVariable",
		"dependsOn": [],
		"userProperties": [],
		"typeProperties": {
			"variableName": "bob",
			"value": "something"
		}
	}
]
  JSON

  depends_on = [azurerm_synapse_firewall_rule.example]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Synapse Pipeline. Changing this forces a new resource to be created. Must be globally unique.

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace that this Synapse Pipeline resides in. Changing this forces a new resource to be created.

* `description` - (Optional) The description for the Synapse Pipeline.

* `annotations` - (Optional)A list of annotation tags that can be used for describing the Synapse Pipeline.

* `parameters` - (Optional) A map of parameters to associate with the Synapse Pipeline. The key is the parameter name, and the value is the parameter's default value. Currently we only support type string for parameter values.

* `variables` - (Optional) A map of variables to associate with the Synapse Pipeline.

* `activities_json` - (Optional) A JSON object that contains the activities that will be associated with the Synapse Pipeline.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Synapse Pipeline.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Pipeline.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Pipeline.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Pipeline.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Pipeline.

## Import

Synapse Pipeline's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_pipeline.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/pipelines/pipeline1
```
