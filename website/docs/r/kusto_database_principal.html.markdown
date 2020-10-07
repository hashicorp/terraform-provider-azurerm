---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_database_principal"
description: |-
  Manages Kusto / Data Explorer Database Principal
---

# azurerm_kusto_database_principal

Manages a Kusto (also known as Azure Data Explorer) Database Principal

~> **NOTE:** This resource is being **deprecated** due to API updates and should no longer be used.  Please use [azurerm_kusto_database_principal_assignment](./kusto_database_principal_assignment.html) instead.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "rg" {
  name     = "my-kusto-rg"
  location = "East US"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "kustocluster"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }
}

resource "azurerm_kusto_database" "database" {
  name                = "my-kusto-database"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name

  hot_cache_period   = "P7D"
  soft_delete_period = "P31D"
}

resource "azurerm_kusto_database_principal" "principal" {
  resource_group_name = azurerm_resource_group.rg.name
  cluster_name        = azurerm_kusto_cluster.cluster.name
  database_name       = azurerm_kusto_database.test.name

  role      = "Viewer"
  type      = "User"
  client_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.client_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Database Principal should exist. Changing this forces a new resource to be created.

* `cluster_name` - (Required) Specifies the name of the Kusto Cluster this database principal will be added to. Changing this forces a new resource to be created.

* `database_name` - (Required) Specified the name of the Kusto Database this principal will be added to. Changing this forces a new resource to be created.

* `role` - (Required) Specifies the permissions the Principal will have. Valid values include `Admin`, `Ingestor`, `Monitor`, `UnrestrictedViewers`, `User`, `Viewer`. Changing this forces a new resource to be created.

* `type` - (Required) Specifies the type of object the principal is. Valid values include `App`, `Group`, `User`. Changing this forces a new resource to be created.

* `object_id` - (Required) An Object ID of a User, Group, or App. Changing this forces a new resource to be created.

* `client_id` - (Required) The Client ID that owns the specified `object_id`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Kusto Database Principal ID.

* `app_id` - The app id, if not empty, of the principal.

* `email` - The email, if not empty, of the principal.

* `fully_qualified_name` - The fully qualified name of the principal.

* `name` - The name of the Kusto Database Principal.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Kusto Database Principal.
* `update` - (Defaults to 60 minutes) Used when updating the Kusto Database Principal.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Database Principal.
* `delete` - (Defaults to 60 minutes) Used when deleting the Kusto Database Principal.

## Import

Kusto Database Principals can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_database_principal.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1/Databases/database1/Role/role1/FQN/some-guid
```
