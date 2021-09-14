---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_data_flow"
description: |-
  Manages a Data Flow inside an Azure Synapse.
---

# azurerm_synapse_data_flow

Manages a Data Flow inside an Azure Synapse.

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

resource "azurerm_synapse_data_flow" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id

  source {
    name = "source1"

    linked_service {
      name = azurerm_synapse_linked_service.example.name
    }
  }

  sink {
    name = "sink1"

    linked_service {
      name = azurerm_synapse_linked_service.example.name
    }
  }

  script = <<EOT
source(
  allowSchemaDrift: true, 
  validateSchema: false, 
  limit: 100, 
  ignoreNoFilesFound: false, 
  documentForm: 'documentPerLine') ~> source1 
source1 sink(
  allowSchemaDrift: true, 
  validateSchema: false, 
  skipDuplicateMapInputs: true, 
  skipDuplicateMapOutputs: true) ~> sink1
EOT
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Data Flow. Changing this forces a new Synapse Data Flow to be created.

* `synapse_workspace_id` - (Required) The ID of the TODO. Changing this forces a new Synapse Data Flow to be created.
  
* `script` - (Required) The script for the Synapse Data Flow.

* `sink` - (Required) One or more `sink` blocks as defined below.

* `source` - (Required) One or more `source` blocks as defined below.

---

* `annotations` - (Optional) List of tags that can be used for describing the Synapse Data Flow.

* `description` - (Optional) The description for the Synapse Data Flow.

* `folder` - (Optional) The folder that this Data Flow is in. If not specified, the Data Flow will appear at the root level.

* `transformation` - (Optional) One or more `transformation` blocks as defined below.

---

A `dataset` block supports the following:

* `name` - (Required) The name for the Synapse Dataset.

* `parameters` - (Optional) A map of parameters to associate with the Synapse dataset.

---

A `linked_service` block supports the following:

* `name` - (Required) The name for the Synapse Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Synapse Linked Service.

---

A `schema_linked_service` block supports the following:

* `name` - (Required) The name for the Synapse Linked Service with schema.

* `parameters` - (Optional) A map of parameters to associate with the Synapse Linked Service.

---

A `sink` block supports the following:

* `name` - (Required)  The name for the Data Flow Source.

* `description` - (Optional) The description for the Data Flow Source.

* `dataset` - (Optional) A `dataset` block as defined above.

* `linked_service` - (Optional) A `linked_service` block as defined above.

* `schema_linked_service` - (Optional) A `schema_linked_service` block as defined above.

---

A `source` block supports the following:

* `name` - (Required) The name for the Data Flow Source.

* `description` - (Optional) The description for the Data Flow Source.
  
* `dataset` - (Optional) A `dataset` block as defined above.

* `linked_service` - (Optional) A `linked_service` block as defined above.

* `schema_linked_service` - (Optional) A `schema_linked_service` block as defined above.

---

A `transformation` block supports the following:

* `name` - (Required) The name for the Data Flow transformation.

* `description` - (Optional) The description for the Data Flow transformation.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse Data Flow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Data Flow.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Data Flow.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Data Flow.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Data Flow.

## Import

Synapse Data Flows can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_data_flow.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/dataflows/dataflow1
```
