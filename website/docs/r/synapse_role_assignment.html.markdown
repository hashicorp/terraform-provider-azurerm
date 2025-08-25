---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_role_assignment"
description: |-
  Manages a Synapse Role Assignment.
---

# azurerm_synapse_role_assignment

Manages a Synapse Role Assignment.

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

data "azurerm_client_config" "current" {}

resource "azurerm_synapse_role_assignment" "example" {
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  role_name            = "Synapse SQL Administrator"
  principal_id         = data.azurerm_client_config.current.object_id

  depends_on = [azurerm_synapse_firewall_rule.example]
}
```

## Argument Reference

The following arguments are supported:

* `synapse_workspace_id` - (Optional) The Synapse Workspace which the Synapse Role Assignment applies to. Changing this forces a new resource to be created.

* `synapse_spark_pool_id` - (Optional) The Synapse Spark Pool which the Synapse Role Assignment applies to. Changing this forces a new resource to be created.

-> **Note:** A Synapse firewall rule including local IP is needed to allow access. Only one of `synapse_workspace_id`, `synapse_spark_pool_id` must be set.

* `role_name` - (Required) The Role Name of the Synapse Built-In Role. Possible values are `Apache Spark Administrator`, `Synapse Administrator`, `Synapse Artifact Publisher`, `Synapse Artifact User`, `Synapse Compute Operator`, `Synapse Contributor`, `Synapse Credential User`, `Synapse Linked Data Manager`, `Synapse Monitoring Operator`, `Synapse SQL Administrator` and `Synapse User`. Changing this forces a new resource to be created.

-> **Note:** Currently, the Synapse built-in roles are `Apache Spark Administrator`, `Synapse Administrator`, `Synapse Artifact Publisher`, `Synapse Artifact User`, `Synapse Compute Operator`, `Synapse Contributor`, `Synapse Credential User`, `Synapse Linked Data Manager`, `Synapse Monitoring Operator`, `Synapse SQL Administrator` and `Synapse User`.

-> **Note:** Old roles are still supported: `Workspace Admin`, `Apache Spark Admin`, `Sql Admin`. These values will be removed in the next Major Version 3.0.

* `principal_id` - (Required) The ID of the Principal (User, Group or Service Principal) to assign the Synapse Role Definition to. Changing this forces a new resource to be created.

* `principal_type` - (Optional) The Type of the Principal. One of `User`, `Group` or `ServicePrincipal`. Changing this forces a new resource to be created.

-> **Note:** While `principal_type` is optional, it's still recommended to set this value, as some Synapse use-cases may not work correctly if this is not specified. Service Principals for example can't run SQL statements using `Entra ID` authentication if `principal_type` is not set to `ServicePrincipal`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Synapse Role Assignment ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Role Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Role Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Role Assignment.

## Import

Synapse Role Assignment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_role_assignment.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1|000000000000"
```

-> **Note:** This ID is specific to Terraform - and is of the format `{synapseScope}|{synapseRoleAssignmentId}`.
