---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_linked_service"
description: |-
  Manages a Linked Service (connection) between a resource and Azure Synapse. This is a generic resource that supports all different Linked Service Types.
---

# azurerm_synapse_linked_service

Manages a Synapse Linked Service.

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

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_firewall_rule" "example" {
  name                 = "allowAll"
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


```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Linked Service. Changing this forces a new Synapse Linked Service to be created.

* `synapse_workspace_id` - (Required) The Synapse Workspace ID in which to associate the Linked Service with. Changing this forces a new Synapse Linked Service to be created.

* `type` - (Required) The type of data stores that will be connected to Synapse. For full list of supported data stores, please refer to [Azure Synapse connector](https://docs.microsoft.com/azure/data-factory/connector-overview). Changing this forces a new Synapse Linked Service to be created.

* `type_properties_json` - (Required) A JSON object that contains the properties of the Synapse Linked Service.

---

* `additional_properties` - (Optional) A map of additional properties to associate with the Synapse Linked Service.

* `annotations` - (Optional) List of tags that can be used for describing the Synapse Linked Service.

* `description` - (Optional) The description for the Synapse Linked Service.

* `integration_runtime` - (Optional) A `integration_runtime` block as defined below.

* `parameters` - (Optional) A map of parameters to associate with the Synapse Linked Service.

---

A `integration_runtime` block supports the following:

* `name` - (Required) The integration runtime reference to associate with the Synapse Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the integration runtime.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Linked Service.

## Import

Synapse Linked Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_linked_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/linkedservices/linkedservice1
```
