---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_database_principal_assignment"
description: |-
  Manages a Kusto / Data Explorer Database Principal Assignment
---

# azurerm_kusto_database_principal_assignment

Manages a Kusto (also known as Azure Data Explorer) Database Principal Assignment.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "KustoRG"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "example" {
  name                = "kustocluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }
}

resource "azurerm_kusto_database" "example" {
  name                = "KustoDatabase"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name

  hot_cache_period   = "P7D"
  soft_delete_period = "P31D"
}

resource "azurerm_kusto_database_principal_assignment" "example" {
  name                = "KustoPrincipalAssignment"
  resource_group_name = azurerm_resource_group.example.name
  cluster_name        = azurerm_kusto_cluster.example.name
  database_name       = azurerm_kusto_database.example.name

  tenant_id      = data.azurerm_client_config.current.tenant_id
  principal_id   = data.azurerm_client_config.current.client_id
  principal_type = "App"
  role           = "Viewer"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the kusto principal assignment. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the resource. Changing this forces a new resource to be created.

* `cluster_name` - (Required) The name of the cluster in which to create the resource. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the database in which to create the resource. Changing this forces a new resource to be created.

* `principal_id` - (Required) The object id of the principal. Changing this forces a new resource to be created.

* `principal_type` - (Required) The type of the principal. Valid values include `App`, `Group`, `User`. Changing this forces a new resource to be created.

* `role` - (Required) The database role assigned to the principal. Valid values include `Admin`, `Ingestor`, `Monitor`, `UnrestrictedViewer`, `User` and `Viewer`. Changing this forces a new resource to be created.

* `tenant_id` - (Required) The tenant id in which the principal resides. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kusto Database Principal Assignment.

* `principal_name` - The name of the principal.

* `tenant_name` - The name of the tenant.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto Database Principal Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Database Principal Assignment.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto Database Principal Assignment.

## Import

Kusto Database Principal Assignment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_database_principal_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/principalAssignments/assignment1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Kusto`: 2024-04-13
