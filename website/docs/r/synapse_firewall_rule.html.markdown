---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_firewall_rule"
description: |-
  Manages a Synapse Firewall Rule.
---

# azurerm_synapse_firewall_rule

Allows you to Manages a Synapse Firewall Rule.

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
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name of the firewall rule. Changing this forces a new resource to be created.

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace on which to create the Firewall Rule. Changing this forces a new resource to be created.

* `start_ip_address` - (Required) The starting IP address to allow through the firewall for this rule.

* `end_ip_address` - (Required) The ending IP address to allow through the firewall for this rule.

-> **NOTE:** The Azure feature `Allow access to Azure services` can be enabled by setting `start_ip_address` and `end_ip_address` to `0.0.0.0`.

## Attributes Reference

The following attributes are exported:

* `id` - The Synapse Firewall Rule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Firewall Rule.

## Import

Synapse Firewall Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_firewall_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.Synapse/workspaces/workspace1/firewallRules/rule1
```
