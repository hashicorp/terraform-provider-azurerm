---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_workspace"
description: |-
  Manages a synapse Workspace.
---

# azurerm_synapse_workspace

Manages a synapse Workspace.

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

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this synapse Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the synapse Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the synapse Workspace should exist. Changing this forces a new resource to be created.

* `storage_data_lake_gen2_filesystem_id` - (Required) Specifies the ID of storage data lake gen2 filesystem resource. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

* `sql_administrator_login` - (Required) Specifies The Login Name of the SQL administrator. Changing this forces a new resource to be created.

* `sql_administrator_login_password` - (Optional) The Password associated with the `sql_administrator_login` for the SQL administrator.

* `managed_virtual_network_enabled` - (Optional) Is Virtual Network enabled for all computes in this workspace. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Synapse Workspace.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Synapse Workspace. At this time the allowed value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the synapse Workspace.

* `connectivity_endpoints` - A list of Connectivity endpoints for this Synapse Workspace.

* `managed_resource_group_name` - Workspace managed resource group.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Synapse Workspace.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Workspace.

## Import

Synapse Workspace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1
```
