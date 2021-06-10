---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_kusto"
description: |-
  Manages a Linked Service (connection) between a Kusto Cluster and Azure Data Factory.
---

# azurerm_data_factory_linked_service_kusto

Manages a Linked Service (connection) between a Kusto Cluster and Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
  }
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
  name                = "my-kusto-database"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name
}

resource "azurerm_data_factory_linked_service_kusto" "example" {
  name                 = "example"
  data_factory_id      = azurerm_data_factory.example.id
  kusto_endpoint       = azurerm_kusto_cluster.example.uri
  kusto_database_name  = azurerm_kusto_database.example.name
  use_managed_identity = true
}

resource "azurerm_kusto_database_principal_assignment" "example" {
  name                = "KustoPrincipalAssignment"
  resource_group_name = azurerm_resource_group.example.name
  cluster_name        = azurerm_kusto_cluster.example.name
  database_name       = azurerm_kusto_database.example.name

  tenant_id      = azurerm_data_factory.example.identity.0.tenant_id
  principal_id   = azurerm_data_factory.example.identity.0.principal_id
  principal_type = "App"
  role           = "Viewer"
}
```

## Argument Reference

The following supported arguments are common across all Azure Data Factory Linked Services:

* `name` - (Required) Specifies the name of the Data Factory Linked Service. Changing this forces a new resource to be created. Must be unique within a data
  factory. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `description` - (Optional) The description for the Data Factory Linked Service.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service.

The following supported arguments are specific to Azure Kusto Linked Service:

* `kusto_endpoint` - (Required) The URI of the Kusto Cluster endpoint.

* `kusto_database_name` - (Required) The Kusto Database Name.

* `use_managed_identity` - (Optional) Whether to use the Data Factory's managed identity to authenticate against the Kusto Database.

* `service_principal_id` - (Optional) The service principal id in which to authenticate against the Kusto Database.

* `service_principal_key` - (Optional) The service principal key in which to authenticate against the Kusto Database.

* `tenant` - (Required) The service principal tenant id or name in which to authenticate against the Kusto Database.

~> **NOTE** If `service_principal_id` is used, `service_principal_key` and `tenant` is also required.

~> **NOTE** One of Managed Identity authentication and Service Principal authentication must be set.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Linked Service.

## Import

Data Factory Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_kusto.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
