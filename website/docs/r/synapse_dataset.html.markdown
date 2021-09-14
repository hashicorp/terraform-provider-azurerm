---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_dataset"
description: |-
  Manages a Dataset inside an Azure Synapse. This is a generic resource that supports all different Dataset Types.
---

# azurerm_synapse_dataset

Manages a Dataset inside an Azure Synapse. This is a generic resource that supports all different Dataset Types.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example"
  location                             = azurerm_resource_group.example.location
  resource_group_name                  = azurerm_resource_group.example.name
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true
}

resource "azurerm_synapse_firewall_rule" "example" {
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

resource "azurerm_synapse_linked_service" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.example.primary_connection_string}"
}
JSON

  depends_on = [
    azurerm_synapse_firewall_rule.example,
  ]
}

resource "azurerm_synapse_dataset" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  type                 = "Json"

  linked_service {
    name = azurerm_synapse_linked_service.example.name
  }

  type_properties_json = <<JSON
{
  "location": {
    "container": "${azurerm_storage_container.example.name}",
    "type": "AzureBlobStorageLocation"
  }
}
JSON
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Dataset. Changing this forces a new Synapse Dataset to be created.

* `synapse_workspace_id` - (Required) The Synapse Workspace ID in which to associate the Dataset with. Changing this forces a new Synapse Dataset to be created.

* `linked_service` - (Required) A `linked_service` block as defined below.

* `type` - (Required) The type of dataset that will be associated with Synapse. Changing this forces a new Synapse Dataset to be created.

* `type_properties_json` - (Required) A JSON object that contains the properties of the Synapse Dataset.

---

* `additional_properties` - (Optional) A map of additional properties to associate with the Synapse Dataset.

* `annotations` - (Optional) List of tags that can be used for describing the Synapse Dataset.

* `description` - (Optional) The description for the Synapse Dataset.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `parameters` - (Optional) A map of parameters to associate with the Synapse Dataset.

* `schema_json` - (Optional) A JSON object that contains the schema of the Synapse Dataset.

---

A `linked_service` block supports the following:

* `name` - (Required) The name of the Synapse Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Synapse Linked Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Dataset.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Dataset.

## Import

Synapse Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_dataset.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/datasets/dataset1
```
