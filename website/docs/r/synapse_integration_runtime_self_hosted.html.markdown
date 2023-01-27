---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_integration_runtime_self_hosted"
description: |-
  Manages a Synapse Self-hosted Integration Runtime.
---

# azurerm_synapse_integration_runtime_self_hosted

Manages a Synapse Self-hosted Integration Runtime.

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

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_firewall_rule" "example" {
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

resource "azurerm_synapse_integration_runtime_self_hosted" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Self-hosted Integration Runtime. Changing this forces a new Synapse Self-hosted Integration Runtime to be created.

* `synapse_workspace_id` - (Required) The Synapse Workspace ID in which to associate the Integration Runtime with. Changing this forces a new Synapse Self-hosted Integration Runtime to be created.

---

* `description` - (Optional) Integration runtime description.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Synapse Self-hosted Integration Runtime.
  
* `authorization_key_primary` - The primary integration runtime authentication key.

* `authorization_key_secondary` - The secondary integration runtime authentication key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Self-hosted Integration Runtime.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Self-hosted Integration Runtime.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Self-hosted Integration Runtime.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Self-hosted Integration Runtime.

## Import

Synapse Self-hosted Integration Runtimes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_integration_runtime_self_hosted.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/integrationRuntimes/IntegrationRuntime1
```
