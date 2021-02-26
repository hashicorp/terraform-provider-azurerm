---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_managed_private_endpoint"
description: |-
  Manages a Synapse Managed Private Endpoint.
---

# azurerm_synapse_managed_private_endpoint

Allows you to Manages a Synapse Managed Private Endpoint.

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

resource "azurerm_storage_account" "example_connect" {
  name                     = "examplestorage2"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "BlobStorage"
}

resource "azurerm_synapse_managed_private_endpoint" "example" {
  name                 = "example-endpoint"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  target_resource_id   = azurerm_storage_account.example_connect.id
  subresource_name     = "blob"

  depends_on = [azurerm_synapse_firewall_rule.example]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Managed Private Endpoint. Changing this forces a new resource to be created.

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace on which to create the Managed Private Endpoint. Changing this forces a new resource to be created.

-> **NOTE:** A Synapse firewall rule including local IP is needed for managing current resource.

* `target_resource_id` - (Required) The ID of the Private Link Enabled Remote Resource which this Synapse Private Endpoint should be connected to. Changing this forces a new resource to be created.

* `subresource_name` - (Required) Specifies the sub resource name which the Synapse Private Endpoint is able to connect to. Changing this forces a new resource to be created.

-> **NOTE:** Possible values are listed in [documentation](https://docs.microsoft.com/en-us/azure/private-link/private-endpoint-overview#dns-configuration).

## Attributes Reference

The following attributes are exported:

* `id` - The Synapse Managed Private Endpoint ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Managed Private Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Managed Private Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Managed Private Endpoint.

## Import

Synapse Managed Private Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_managed_private_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/managedVirtualNetworks/default/managedPrivateEndpoints/endpoint1
```
