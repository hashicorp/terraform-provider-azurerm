---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_workspace_extended_auditing_policy"
description: |-
  Manages a Synapse Workspace Extended Auditing Policy.
---

# azurerm_synapse_workspace_extended_auditing_policy

Manages a Synapse Workspace Extended Auditing Policy.

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
  account_kind             = "BlobStorage"
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

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "audit_logs" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_synapse_workspace_extended_auditing_policy" "example" {
  synapse_workspace_id                    = azurerm_synapse_workspace.example.id
  storage_endpoint                        = azurerm_storage_account.audit_logs.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.audit_logs.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
}
```

## Arguments Reference

The following arguments are supported:

* `synapse_workspace_id` - (Required) The ID of the Synapse workspace to set the extended auditing policy. Changing this forces a new resource to be created.

* `storage_endpoint` - (Optional) The blob storage endpoint (e.g. https://example.blob.core.windows.net). This blob storage will hold all extended auditing logs.

* `retention_in_days` - (Optional) The number of days to retain logs for in the storage account.

* `storage_account_access_key` - (Optional) The access key to use for the auditing storage account.

* `storage_account_access_key_is_secondary` - (Optional) Is `storage_account_access_key` value the storage's secondary key?

* `log_monitoring_enabled` - (Optional) Enable audit events to Azure Monitor? To enable server audit events to Azure Monitor, please enable its master database audit events to Azure Monitor.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Synapse Workspace Extended Auditing Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Workspace Extended Auditing Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Workspace Extended Auditing Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Workspace Extended Auditing Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Workspace Extended Auditing Policy.

## Import

Synapse Workspace Extended Auditing Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_workspace_extended_auditing_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/extendedAuditingSettings/default
```
