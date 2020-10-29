---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_role_assignment"
description: |-
  Manages a Synapse Role Assignment.
---

# azurerm_synapse_role_assignment

Allows you to Manages a Synapse Role Assignment.

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
}

resource "azurerm_synapse_firewall_rule" "example" {
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

data "azurerm_client_config" "current" {}

resource "azurerm_synapse_role_assignment" "example" {
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  role_name            = "Sql Admin"
  principal_id         = data.azurerm_client_config.current.object_id

  depends_on = [azurerm_synapse_firewall_rule.example]
}
```

## Argument Reference

The following arguments are supported:

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace on which to create the Role Assignment. Changing this forces a new resource to be created.

-> **NOTE:** A Synapse firewall rule including local IP is needed to allow access.

* `role_name` - (Required) The Role Name of the Synapse Built-In Role. Changing this forces a new resource to be created.

-> **NOTE:** Currently, the Synapse built-in roles are `Workspace Admin`, `Apache Spark Admin` and `Sql Admin`.

* `principal_id` - (Required) The ID of the Principal (User, Group or Service Principal) to assign the Synapse Role Definition to. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Synapse Role Assignment ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Role Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Role Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Role Assignment.

## Import

Synapse Role Assignment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_role_assignment.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1|000000000000"
```

-> **NOTE:** This ID is specific to Terraform - and is of the format `{synapseWorkspaceId}|{synapseRoleAssignmentId}`.
